package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/manther/simplebank/db/util"
	"github.com/stretchr/testify/require"
)

func createRandomEntry(t *testing.T) (Account, Entry) {
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

	argE := CreateEntryParams{
		Amount:    util.RandomMoney(),
		AccountID: account.ID,
	}

	entry, err := testQueries.CreateEntry(context.Background(), argE)
	require.NoError(t, err)
	require.NotEmpty(t, entry)

	return account, entry
}

func TestCreateEntry(t *testing.T) {
	createRandomEntry(t)
}

func TestGetEntry(t *testing.T) {
	account, entry := createRandomEntry(t)
	entryRet, err := testQueries.GetEntry(context.Background(), entry.ID)
	require.NoError(t, err)
	require.NotEmpty(t, entryRet)
	require.Equal(t, entry.ID, entryRet.ID)
	require.Equal(t, entry.AccountID, account.ID)
	require.Equal(t, entry.Amount, entryRet.Amount)
	require.WithinRange(t, entry.CreatedAt, entryRet.CreatedAt, time.Now().Add(time.Second))
}

func TestUpdateEntry(t *testing.T) {
	account, entry := createRandomEntry(t)
	entUpParms := UpdatEntryParams{
		ID:        entry.ID,
		AccountID: account.ID,
		Amount:    entry.Amount,
	}
	entryUpd, err := testQueries.UpdatEntry(context.Background(), entUpParms)
	require.NoError(t, err)
	require.NotEmpty(t, entryUpd)
	require.Equal(t, entryUpd.ID, entUpParms.ID)
	require.Equal(t, entryUpd.AccountID, entUpParms.AccountID)
	require.Equal(t, entryUpd.Amount, entUpParms.Amount)
	require.WithinRange(t, entryUpd.CreatedAt, account.CreatedAt, time.Now().Add(time.Second))
}

func TestDeletEntry(t *testing.T) {
	_, entry := createRandomEntry(t)

	entrRet, err := testQueries.GetEntry(context.Background(), entry.ID)
	require.NoError(t, err)
	require.NotEmpty(t, entrRet)
	require.Equal(t, entry.ID, entrRet.ID)

	err = testQueries.DeleteEntriy(context.Background(), entrRet.ID)
	require.NoError(t, err)

	entrRet, err = testQueries.GetEntry(context.Background(), entry.ID)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, entrRet)
	require.NotEqual(t, entry.ID, entrRet.ID)
}

func TestListEntriess(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomEntry(t)
	}
	arg := ListEntriesParams{
		Limit:  5,
		Offset: 5,
	}

	accounts, err := testQueries.ListEntries(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, accounts, 5)

	for _, accnt := range accounts {
		require.NotEmpty(t, accnt)
	}
}
