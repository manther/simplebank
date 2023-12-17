package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/manther/simplebank/db/util"
	"github.com/stretchr/testify/require"
)

func createRandomAccount(t *testing.T) Account {
	arg := CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	accountNew := createRandomAccount(t)
	accountRet, err := testQueries.GetAccount(context.Background(), accountNew.ID)
	require.NoError(t, err)
	require.NotEmpty(t, accountRet)
	require.Equal(t, accountNew.ID, accountRet.ID)
	require.Equal(t, accountNew.Currency, accountRet.Currency)
	require.Equal(t, accountNew.Owner, accountRet.Owner)
	require.Equal(t, accountNew.Balance, accountRet.Balance)
	require.WithinRange(t, accountNew.CreatedAt, accountRet.CreatedAt, time.Now().Add(time.Second))
}

func TestUpdateAccount(t *testing.T) {
	account := createRandomAccount(t)
	accntUpParms := UpdateAccountParams{
		ID:      account.ID,
		Balance: util.RandomMoney(),
	}
	accoutUp, err := testQueries.UpdateAccount(context.Background(), accntUpParms)
	require.NoError(t, err)
	require.NotEmpty(t, accoutUp)

	require.Equal(t, accoutUp.ID, account.ID)
	require.Equal(t, accoutUp.Currency, account.Currency)
	require.Equal(t, accoutUp.Owner, account.Owner)
	require.Equal(t, accoutUp.Balance, accntUpParms.Balance)
	require.WithinRange(t, accoutUp.CreatedAt, account.CreatedAt, time.Now().Add(time.Second))
}

func TestDeletAccount(t *testing.T) {
	account := createRandomAccount(t)

	accountRet, err := testQueries.GetAccount(context.Background(), account.ID)
	require.NoError(t, err)
	require.NotEmpty(t, accountRet)
	require.Equal(t, account.ID, accountRet.ID)

	err = testQueries.DeleteAccount(context.Background(), account.ID)
	require.NoError(t, err)

	accountRet, err = testQueries.GetAccount(context.Background(), account.ID)
	require.Empty(t, accountRet)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
}

func TestListAccounts(t *testing.T){
	for i := 0; i < 10; i++ {
		createRandomAccount(t)
	}
	arg := ListAccountsParams{
		Limit: 5,
		Offset: 5,
	}

	accounts, err := testQueries.ListAccounts(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, accounts, 5)

	for _, accnt := range accounts {
		require.NotEmpty(t, accnt)
	}
}
