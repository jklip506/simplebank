package db

import (
	"context"
	"database/sql"
	"fmt"
	"simplebank/db/utils"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestQueries_CreateAccount(t *testing.T) {

	// Prepare test data
	arg := CreateAccountParams{
		Owner:    utils.RandomOwner(),
		Balance:  utils.RandomMoney(),
		Currency: utils.RandomCurrency(),
	}

	// Create an account
	account, err := testQueries.CreateAccount(context.Background(), arg)
	if err != nil {
		t.Fatalf("failed to create account: %v", err)
	}

	// Verify the account details
	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)
	
	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	// Clean up the test data


	// Close the test database connection
	// ...
}

func TestQueries_GetAccount(t *testing.T) {
	// Prepare test data
	arg := CreateAccountParams{
		Owner:    utils.RandomOwner(),
		Balance:  utils.RandomMoney(),
		Currency: utils.RandomCurrency(),
	}
	// Create an account
	account1, err := testQueries.CreateAccount(context.Background(), arg)
	if err != nil {
		t.Fatalf("failed to create account: %v", err)
	}
	// Get the account
	account2, err := testQueries.GetAccount(context.Background(), account1.ID)
	if err != nil {
		t.Fatalf("failed to get account: %v", err)
	}
	// Verify the account details
	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, account1.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)
}

func TestQueries_DeleteAccount(t *testing.T) {
    // Prepare test data
    arg := CreateAccountParams{
        Owner:    utils.RandomOwner(),
        Balance:  utils.RandomMoney(),
        Currency: utils.RandomCurrency(),
    }
    // Create an account
    account1, err := testQueries.CreateAccount(context.Background(), arg)
    if err != nil {
        t.Fatalf("failed to create account: %v", err)
    }
    // Delete the account
    deletedAccount, err := testQueries.DeleteAccount(context.Background(), account1.ID)
    if err != nil {
        t.Fatalf("failed to delete account: %v", err)
    }
    // Verify the deleted account is the same as the created one
    require.NotEmpty(t, deletedAccount)
    require.Equal(t, account1.ID, deletedAccount.ID)
    // Get the account
    account2, err := testQueries.GetAccount(context.Background(), account1.ID)
    if err != nil {
        if err == sql.ErrNoRows {
            // This is what we want! We successfully deleted the account.
        } else {
            t.Fatalf("failed to get account: %v", err)
        }
    } else {
        // If we didn't get an error, that means we fetched an account. This is not what we want.
        t.Fatalf("unexpectedly retrieved a deleted account: %v", account2)
    }
    // Verify the account is deleted
    require.Empty(t, account2)
}

func TestQueries_UpdateAccount(t *testing.T) {
    // Create a test account
    createArg := CreateAccountParams{
        Owner:    utils.RandomOwner(),
        Balance:  utils.RandomMoney(),
        Currency: utils.RandomCurrency(),
    }
    account, err := testQueries.CreateAccount(context.Background(), createArg)
    require.NoError(t, err)

    // Prepare test data
    arg := UpdateAccountParams{
        ID:      account.ID, // Use the ID of the created account
        Balance: 1000,
    }

    // Update the account
    account, err = testQueries.UpdateAccount(context.Background(), arg)
    require.NoError(t, err)

    // Verify the updated account details
    require.Equal(t, arg.ID, account.ID)
    require.Equal(t, arg.Balance, account.Balance)

    // Get the updated account
    updatedAccount, err := testQueries.GetAccount(context.Background(), arg.ID)
    require.NoError(t, err)

    // Verify the updated account details
    require.Equal(t, arg.ID, updatedAccount.ID)
    require.Equal(t, arg.Balance, updatedAccount.Balance)
    require.Equal(t, account.Owner, updatedAccount.Owner)
    require.Equal(t, account.Currency, updatedAccount.Currency)
    require.WithinDuration(t, account.CreatedAt, updatedAccount.CreatedAt, time.Second)
}

func TestQueries_ListAccounts(t *testing.T) {
	// Prepare test data
	arg := ListAccountsParams{
		Limit:  10,
		Offset: 0,
	}

	// Insert test accounts
	limit := int(arg.Limit)
	for i := 0; i < limit; i++ {
		createArg := CreateAccountParams{
			Owner:    utils.RandomOwner(),
			Balance:  utils.RandomMoney(),
			Currency: utils.RandomCurrency(),
		}
		_, err := testQueries.CreateAccount(context.Background(), createArg)
		require.NoError(t, err)
	}

	// Call the ListAccounts function
	limit = int(arg.Limit)
	accounts, err := testQueries.ListAccounts(context.Background(), arg)
	require.NoError(t, err)

	// Print out the accounts
	for _, account := range accounts {
		fmt.Printf("Account: %+v\n", account)
	}

	// Verify the number of accounts returned
	require.Len(t, accounts, limit)

	// Verify the account details
	for _, account := range accounts {
		require.NotZero(t, account.ID)
		require.NotEmpty(t, account.Owner)
		require.NotZero(t, account.Balance)
		require.NotEmpty(t, account.Currency)
		require.NotZero(t, account.CreatedAt)
	}

}

func TestQueries_ListAccounts_Pagination(t *testing.T) {
    // Assume accounts are already created in setup

    // Test with different limit and offset
    arg := ListAccountsParams{
        Limit:  5,
        Offset: 5,
    }

    // Call the ListAccounts function
    accounts, err := testQueries.ListAccounts(context.Background(), arg)
    require.NoError(t, err)

    // Verify the number of accounts returned
    require.Len(t, accounts, 5)
}

// Add more tests for boundary cases, error handling, and invalid parameters
func TestQueries_AddAccountBalance(t *testing.T) {
	// Prepare test data
	arg := CreateAccountParams{
		Owner:    utils.RandomOwner(),
		Balance:  utils.RandomMoney(),
		Currency: utils.RandomCurrency(),
	}
	// Create an account
	account, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)

	// Prepare test data for adding balance
	addArg := AddAccountBalanceParams{
		ID:     account.ID,
		Amount: 100,
	}

	// Add balance to the account
	updatedAccount, err := testQueries.AddAccountBalance(context.Background(), addArg)
	require.NoError(t, err)

	// Verify the updated account details
	require.Equal(t, account.ID, updatedAccount.ID)
	require.Equal(t, account.Owner, updatedAccount.Owner)
	require.Equal(t, account.Currency, updatedAccount.Currency)
	require.Equal(t, account.CreatedAt, updatedAccount.CreatedAt)
	require.Equal(t, account.Balance+addArg.Amount, updatedAccount.Balance)
}


// Test for updating a non-existing account
func TestQueries_UpdateAccount_NonExisting(t *testing.T) {
    arg := UpdateAccountParams{
        ID:      -1, // Non-existing ID
        Balance: 1000,
    }

    _, err := testQueries.UpdateAccount(context.Background(), arg)
    require.Error(t, err)
}

// Test for deleting a non-existing account
func TestQueries_DeleteAccount_NonExisting(t *testing.T) {
    account, err := testQueries.DeleteAccount(context.Background(), -1) // Non-existing ID
    require.ErrorIs(t, err, sql.ErrNoRows)
    require.Empty(t, account)
}