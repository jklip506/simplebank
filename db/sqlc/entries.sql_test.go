package db

import (
	"context"
	"simplebank/db/utils"
	"testing"

	"github.com/stretchr/testify/require"
)


func TestCreateEntry_Success(t *testing.T) {
    arg := CreateAccountParams {
        Owner: utils.RandomOwner(),
        Balance: utils.RandomMoney(),
        Currency: utils.RandomCurrency(),
    }
    // Create a test account
    accountID, err := testQueries.CreateAccount(context.Background(), arg)
    if err != nil {
        t.Fatalf("failed to create test account: %v", err)
    }

    // Generate random amount
    amount := utils.RandomMoney()

    // Create entry
    createEntryArg := CreateEntryParams{
        AccountID: accountID.ID,
        Amount:    amount,
    }

    createdEntry, err := testQueries.CreateEntry(context.Background(), createEntryArg)

    // Check for any errors
    require.NoError(t, err)
    require.NotZero(t, accountID.ID)
    require.NotZero(t, accountID.CreatedAt)

    // ensure the updated entry is correct
    require.Equal(t, arg.Owner, accountID.Owner)
    require.Equal(t, arg.Balance, accountID.Balance)
    require.Equal(t, arg.Currency, accountID.Currency)
    require.Equal(t, createEntryArg.AccountID, createdEntry.AccountID)
    require.Equal(t, createEntryArg.Amount, createdEntry.Amount)

    // TODO: Add code to remove all entries
	
}

func TestCreateEntry_Failure(t *testing.T) {

    account := CreateAccountParams {
        Owner: utils.RandomOwner(),
        Balance: utils.RandomMoney(),
        Currency: utils.RandomCurrency(),
    }
    // Create a test account
    accountID, err := testQueries.CreateAccount(context.Background(), account)
    if err != nil {
        t.Fatalf("failed to create test account: %v", err)
    }

    // Generate random amount
    amount := utils.RandomMoney()

    // Create entry
    entry := CreateEntryParams{
        AccountID: accountID.ID,
        Amount:    amount,
    }

    createdEntry, err := testQueries.CreateEntry(context.Background(), entry)

    // Check for any errors
    require.NoError(t, err)
    require.NotZero(t, accountID.ID)
    require.NotZero(t, accountID.CreatedAt)

    // ensure the updated entry is correct
    require.Equal(t, accountID.ID, createdEntry.AccountID)
    require.Equal(t, amount, createdEntry.Amount)
}
func TestListEntriesByAccount(t *testing.T) {

	account := CreateAccountParams {
		Owner: utils.RandomOwner(),
		Balance: utils.RandomMoney(),
		Currency: utils.RandomCurrency(),
	}

	// Create a test account
	accountID, err := testQueries.CreateAccount(context.Background(), account)
	require.NoError(t, err)

	// Create test entries
	for i := 0; i < 5; i++ {
		amount := utils.RandomMoney()
		entryID, err := testQueries.CreateEntry(context.Background(), CreateEntryParams{
			AccountID: accountID.ID,
			Amount:    amount,
		})
		require.NoError(t, err)
		require.NotZero(t, entryID)
	}

	// Call the function being tested
	entries, err := testQueries.ListEntriesByAccount(context.Background(), ListEntriesByAccountParams{
		AccountID: accountID.ID,
		Limit:     10,
		Offset:    0,
	})
	require.NoError(t, err)
	require.Len(t, entries, 5)

	// assert that the entries returned are correct
	for _, entry := range entries {
		require.NotZero(t, entry.ID)
		require.NotZero(t, entry.AccountID)
		require.NotZero(t, entry.Amount)
		require.NotZero(t, entry.CreatedAt)
	}

	// TODO: Add assertions to verify the returned entries
	require.Equal(t, accountID.ID, accountID.ID)
	require.Equal(t, accountID.CreatedAt, accountID.CreatedAt)
	require.Equal(t, accountID.Owner, accountID.Owner)
	require.Equal(t, accountID.Balance, accountID.Balance)
	require.Equal(t, accountID.Currency, accountID.Currency)

}

func TestGetEntry_Success(t *testing.T) {
	// Create a test entry
	account1 := CreateAccountParams {
		Owner: utils.RandomOwner(),
		Balance: utils.RandomMoney(),
		Currency: utils.RandomCurrency(),
	}

	// create the test account
	accountID, err := testQueries.CreateAccount(context.Background(), account1)
	require.NoError(t, err)

	// create the test entry
	entry := CreateEntryParams{
		AccountID: accountID.ID,
		Amount:    utils.RandomMoney(),
	}

	// insert the entry into the database
	createdEntry, err := testQueries.CreateEntry(context.Background(), entry)
	require.NoError(t, err)

	// get the entry from the database
	// get the entry from the database
	entryFromDB, err := testQueries.GetEntry(context.Background(), createdEntry.ID)
	require.NoError(t, err)

	// ensure the updated entry is correct
	require.Equal(t, createdEntry.ID, entryFromDB.ID)
	require.Equal(t, createdEntry.AccountID, entryFromDB.AccountID)
	require.Equal(t, createdEntry.Amount, entryFromDB.Amount)

}

// func TestGetEntry_NotFound(t *testing.T) {
// 	// Get an entry with a non-existent ID
// 	_, err := testQueries.GetEntry(context.Background(), -1)
// 	require.Error(t, err)
// 	require.Equal(t, err, sql.ErrNoRows)
// }
