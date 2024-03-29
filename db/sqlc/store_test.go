package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	// It is good to check if our transaction runs well while being one of the many transactions
	// running concurrently, to do that we will
	// run n concurrent transfer transactions
	n := 5
	amount := int64(10)

	// Channels are used to send data between concurrently running goroutines
	// This TestTransferTX function will also we run in a goroutine which will be the main one
	// To send the errors and results from invoked goroutines we declare channels below
	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		// To utilize concurrency we will use goroutines
		go func() {
			ctx := context.Background()
			result, err := store.TransferTX(ctx, TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})
			// In a normal test here we would have the assertions using require package
			// But we cannot do it because each of the functions is run concurrently in a different goroutine and
			// we cannot be sure if the execution would be stopped if it was requested when the error occurred
			// That is why we use channels to send errors and results from each of concurrently running tests to to main goroutine
			// where they will be asserted
			errs <- err
			results <- result
		}()
		// Bracket above after definition of a anonymous function is necessary to run it
	}

	// Here we check the results from all tests
	existed := make(map[int]bool)

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		// Check transfer
		// Check values from the result
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, account1.ID, transfer.FromAccountID)
		require.Equal(t, account2.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		// Also check if it was added to db
		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		// Check entires
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, account1.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, account2.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		// Check accounts
		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, account1.ID, fromAccount.ID)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, account2.ID, toAccount.ID)

		// Check accounts' balance
		diff1 := account1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - account2.Balance
		require.Equal(t, diff1, diff2)
		// The amount that was subtracted from the account1 cannot be negative
		require.True(t, diff1 > 0)
		// This amount must also be divisible by the amount(var amount) of money transferred in each transaction
		require.True(t, diff1%amount == 0) // amount, 2*amount, 3*amount, ...

		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= n)
		// We also want to check if k is unique for each transaction
		require.NotContains(t, existed, k)
		existed[k] = true
	}

	// Check the final updated balance, after all transactions
	updatedAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	require.Equal(t, account1.Balance-int64(n)*amount, updatedAccount1.Balance)
	require.Equal(t, account2.Balance+int64(n)*amount, updatedAccount2.Balance)
}

func TestTransferTxDeadlock(t *testing.T) {
	store := NewStore(testDB)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	// In this example we will run 5 transactions transferring money from account1 to account2
	// and 5 transactions doing the opposite
	// if not handled properly this would cause deadlock, because transactions would like to access same accounts at the same time
	n := 10
	amount := int64(10)
	errs := make(chan error)

	for i := 0; i < n; i++ {
		// for 5 of 10 tests fromAccount is account1 and toAccount is account2
		fromAccountID := account1.ID
		toAccountID := account2.ID

		// for 5 of 10 tests fromAccount is account2 and toAccount is account1
		if i%2 == 1 {
			fromAccountID = account2.ID
			toAccountID = account1.ID
		}

		go func() {
			ctx := context.Background()
			_, err := store.TransferTX(ctx, TransferTxParams{
				FromAccountID: fromAccountID,
				ToAccountID:   toAccountID,
				Amount:        amount,
			})
			errs <- err
		}()
	}

	// Here we check the results from all tests
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
	}

	// Check the final updated balance, after all transactions
	updatedAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	// This time we expect the balances to be exactly the same before and after transfers
	// because 5 transfers will subtract money from them and 5 will add money to them
	require.Equal(t, account1.Balance, updatedAccount1.Balance)
	require.Equal(t, account2.Balance, updatedAccount2.Balance)
}
