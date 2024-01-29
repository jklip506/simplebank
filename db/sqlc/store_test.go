package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)


func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	// Create test accounts
	fromAccount := createRandomAccount(t)
	toAccount := createRandomAccount(t)


	// Run n concurrent transfer transactions
	n := 5
	amount := int64(10)

	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		go func() {
			result, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: fromAccount.ID,
				ToAccountID:   toAccount.ID,
				Amount:        amount,
			})

			errs <- err
			results <- result
		}()
	}

	// Check results
	existed := make(map[int]bool)
	for i := 0; i < n; i++ { 
		err := <- errs
		require.NoError(t, err)

		result := <- results
		require.NotEmpty(t, result)

		// Check the transfer details
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, fromAccount.ID, transfer.FromAccountID)
		require.Equal(t, toAccount.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		// get transfer from db to check it was created
		transferFromDB, err := store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)
		require.NotEmpty(t, transferFromDB)

		// check the account entries
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, fromEntry.AccountID, fromAccount.ID)
		require.Equal(t, fromEntry.Amount, -amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		// get entry from db to check it was created
		fromEntryFromDB, err := store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)
		require.NotEmpty(t, fromEntryFromDB)

		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, toEntry.AccountID, toAccount.ID)
		require.Equal(t, toEntry.Amount, amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		// get entry from db to check it was created
		toEntryFromDB, err := store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)
		require.NotEmpty(t, toEntryFromDB)

		// check accounts
		account1 := result.FromAccount
		require.NotEmpty(t, account1)
		require.Equal(t, account1.ID, fromAccount.ID)

		account2 := result.ToAccount
		require.NotEmpty(t, account2)
		require.Equal(t, account2.ID, toAccount.ID)

		// check account balances
		diff1 := fromAccount.Balance - account1.Balance
		diff2 := account2.Balance - toAccount.Balance
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0) // amount, 2 * amount, 3 * amount, etc... n * amount

		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true

	}

	// check the final updated account balances
	updatedAccount1, err := testQueries.GetAccount(context.Background(), fromAccount.ID)
	require.NoError(t, err)
	require.NotEmpty(t, updatedAccount1)

	updatedAccount2, err := testQueries.GetAccount(context.Background(), toAccount.ID)
	require.NoError(t, err)
	require.NotEmpty(t, updatedAccount2)

	require.Equal(t, fromAccount.Balance-int64(n)*amount, updatedAccount1.Balance)
	require.Equal(t, toAccount.Balance+int64(n)*amount, updatedAccount2.Balance)

}