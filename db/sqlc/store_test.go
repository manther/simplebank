package db

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestTransferTX(t *testing.T) {
	store := NewStore(testDb)

	accnt1 := createRandomAccount(t)
	accnt2 := createRandomAccount(t)
	fmt.Println(">>Before:", accnt1.Balance, accnt2.Balance)

	// run multiple concurrent transfer transactions
	n := 5
	amount := int64(10)

	errs := make(chan error)
	results := make(chan TransferTxResult)

	existed := make(map[int]bool)
	for i := 0; i < n; i++ {
		go func() {
			result, err := store.TransferTX(context.Background(), TransferTxParms{
				FromAccountID: accnt1.ID,
				ToAccountID:   accnt2.ID,
				Amount:        amount,
			})

			errs <- err
			results <- result

		}()
	}

	for i := 0; i < n; i++ {

		err := <-errs
		require.NoError(t, err)

		res := <-results

		require.NotEmpty(t, res)

		// Check Transfer
		require.Equal(t, accnt1.ID, res.Transfer.FromAccount)
		require.Equal(t, accnt2.ID, res.Transfer.ToAccount)
		require.Equal(t, amount, res.Transfer.Amount)
		require.NotZero(t, res.Transfer.ID)
		require.WithinDuration(t, time.Now(), res.Transfer.CreatedAt, 1*time.Second)

		_, err = store.GetTransfer(context.Background(), res.Transfer.ID)
		require.NoError(t, err)
		// fmt.Println("Transfer ok. Transfer ID:", res.Transfer.ID)

		// Check entries - From
		require.NotEmpty(t, res.FromEntry)
		require.Equal(t, accnt1.ID, res.FromEntry.AccountID)
		require.Equal(t, -amount, res.FromEntry.Amount)
		require.NotZero(t, res.FromEntry.ID)
		require.WithinDuration(t, time.Now(), res.FromEntry.CreatedAt, 1*time.Second)

		_, err = store.GetEntry(context.Background(), res.FromEntry.ID)
		require.NoError(t, err)
		// fmt.Println("From Entry ok. Transfer ID:", res.Transfer.ID)

		// Check entries - To
		require.NotEmpty(t, res.ToEntry)
		require.Equal(t, accnt2.ID, res.ToEntry.AccountID)
		require.Equal(t, amount, res.ToEntry.Amount)
		require.NotZero(t, res.ToEntry.ID)
		require.WithinDuration(t, time.Now(), res.ToEntry.CreatedAt, 1*time.Second)

		_, err = store.GetEntry(context.Background(), res.ToEntry.ID)
		require.NoError(t, err)
		// fmt.Println("To ok. Transfer ID:", res.Transfer.ID)

		// Check account
		require.NotEmpty(t, res.FromAccount)
		require.Equal(t, accnt1.ID, res.FromAccount.ID)

		require.NotEmpty(t, res.ToAccount)
		require.Equal(t, accnt2.ID, res.ToAccount.ID)

		// Check Accnt balance
		fmt.Println(">>tx:", res.FromAccount.Balance, res.ToAccount.Balance)
		diff1 := accnt1.Balance - res.FromAccount.Balance
		diff2 := res.ToAccount.Balance - accnt2.Balance
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0)
		require.True(t, diff2%amount == 0)

		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true
	}
	updateAccount1, err := testQueries.GetAccount(context.Background(), accnt1.ID)
	require.NoError(t, err)

	updateAccount2, err := testQueries.GetAccount(context.Background(), accnt2.ID)
	require.NoError(t, err)

	fmt.Println(">>Balance after update:", updateAccount1.Balance, updateAccount2.Balance)
	require.Equal(t, accnt1.Balance-int64(n)*amount, updateAccount1.Balance)
	require.Equal(t, accnt2.Balance+int64(n)*amount, updateAccount2.Balance)
}

func TestTransferTXDeadLock(t *testing.T) {
	store := NewStore(testDb)

	accnt1 := createRandomAccount(t)
	accnt2 := createRandomAccount(t)
	fmt.Println(">>Before:", accnt1.Balance, accnt2.Balance)

	// run multiple concurrent transfer transactions
	n := 10
	amount := int64(10)
	errs := make(chan error)

	for i := 0; i < n; i++ {
		fromAccountID := accnt1.ID
		toAccountID := accnt2.ID

		if i%2 == 0 {
			fromAccountID = accnt2.ID
			toAccountID = accnt1.ID
		}

		go func() {
			res, err := store.TransferTX(context.Background(), TransferTxParms{
				FromAccountID: fromAccountID,
				ToAccountID:   toAccountID,
				Amount:        amount,
			})

			fmt.Println(">>tx:", res.FromAccount.Balance, res.ToAccount.Balance)

			errs <- err
		}()
	}

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
	}

	updateAccount1, err := testQueries.GetAccount(context.Background(), accnt1.ID)
	require.NoError(t, err)

	updateAccount2, err := testQueries.GetAccount(context.Background(), accnt2.ID)
	require.NoError(t, err)

	fmt.Println(">>Balance after update:", updateAccount1.Balance, updateAccount2.Balance)
	require.Equal(t, accnt1.Balance, updateAccount1.Balance)
	require.Equal(t, accnt2.Balance, updateAccount2.Balance)
}
