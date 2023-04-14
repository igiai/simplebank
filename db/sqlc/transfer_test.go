package db

import (
	"context"
	"testing"

	"github.com/igiai/simplebank/db/util"
	"github.com/stretchr/testify/require"
)

func createRandomTransferForAccounts(t *testing.T, fromAccount Account, toAccount Account) Transfer {
	arg := CreateTransferParams{
		FromAccountID: fromAccount.ID,
		ToAccountID:   toAccount.ID,
		Amount:        util.RandomMoney(),
	}

	transfer, err := testQueries.CreateTransfer(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, arg.FromAccountID, transfer.FromAccountID)
	require.Equal(t, arg.ToAccountID, transfer.ToAccountID)
	require.Equal(t, arg.Amount, transfer.Amount)

	require.NotEmpty(t, transfer.ID)
	require.NotEmpty(t, transfer.CreatedAt)

	return transfer
}

func TestCreateTransfer(t *testing.T) {
	fromAccount := createRandomAccount(t)
	toAccount := createRandomAccount(t)

	createRandomTransferForAccounts(t, fromAccount, toAccount)
}

func TestGetTransfer(t *testing.T) {
	fromAccount := createRandomAccount(t)
	toAccount := createRandomAccount(t)
	expectedTransfer := createRandomTransferForAccounts(t, fromAccount, toAccount)

	actualTransfer, err := testQueries.GetTransfer(context.Background(), expectedTransfer.ID)
	require.NoError(t, err)
	require.NotEmpty(t, actualTransfer)

	require.Equal(t, expectedTransfer.ID, actualTransfer.ID)
	require.Equal(t, expectedTransfer.FromAccountID, actualTransfer.FromAccountID)
	require.Equal(t, expectedTransfer.ToAccountID, actualTransfer.ToAccountID)
	require.Equal(t, expectedTransfer.Amount, actualTransfer.Amount)
	require.Equal(t, expectedTransfer.CreatedAt, actualTransfer.CreatedAt)
}

func TestListTransfers(t *testing.T) {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	for i := 0; i < 5; i++ {
		createRandomTransferForAccounts(t, account1, account2)
		createRandomTransferForAccounts(t, account2, account1)
	}

	arg := ListTransfersParams{
		FromAccountID: account1.ID,
		ToAccountID:   account1.ID,
		Limit:         5,
		Offset:        5,
	}

	transfers, err := testQueries.ListTransfers(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, transfers, 5)

	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)
		require.True(t, transfer.FromAccountID == account1.ID || transfer.ToAccountID == account1.ID)
	}
}
