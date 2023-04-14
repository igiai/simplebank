package db

import (
	"context"
	"testing"

	"github.com/igiai/simplebank/db/util"
	"github.com/stretchr/testify/require"
)

func createRandomEntryForAccount(t *testing.T, account Account) Entry {
	arg := CreateEntryParams{
		AccountID: account.ID,
		Amount:    util.RandomMoney(),
	}

	entry, err := testQueries.CreateEntry(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.Equal(t, arg.AccountID, entry.AccountID)
	require.Equal(t, arg.Amount, entry.Amount)

	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)

	return entry
}

func TestCreateEntry(t *testing.T) {
	account := createRandomAccount(t)
	createRandomEntryForAccount(t, account)
}

func TestGetEntry(t *testing.T) {
	account := createRandomAccount(t)
	expectedEntry := createRandomEntryForAccount(t, account)

	actualEntry, err := testQueries.GetEntry(context.Background(), expectedEntry.ID)
	require.NoError(t, err)
	require.NotEmpty(t, actualEntry)

	require.Equal(t, expectedEntry.ID, actualEntry.ID)
	require.Equal(t, expectedEntry.AccountID, actualEntry.AccountID)
	require.Equal(t, expectedEntry.Amount, actualEntry.Amount)
	require.Equal(t, expectedEntry.CreatedAt, actualEntry.CreatedAt)
}

func TestListEntries(t *testing.T) {
	account := createRandomAccount(t)
	for i := 0; i < 10; i++ {
		createRandomEntryForAccount(t, account)
	}

	arg := ListEntriesParams{
		AccountID: account.ID,
		Limit:     5,
		Offset:    5,
	}

	entries, err := testQueries.ListEntries(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, entries, 5)

	for _, entry := range entries {
		require.NotEmpty(t, entry)
		require.Equal(t, arg.AccountID, entry.AccountID)
	}
}
