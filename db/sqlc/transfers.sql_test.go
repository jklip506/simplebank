package db

import (
	"context"
	"simplebank/db/utils"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateTransfer(t *testing.T) {
    CreateAccountParams := CreateAccountParams{
        Owner:    utils.RandomOwner(),
        Balance:  utils.RandomMoney(),
        Currency: utils.RandomCurrency(),
    }

    // Create a test account
    account1, err := testQueries.CreateAccount(context.Background(), CreateAccountParams)
    require.NoError(t, err)

    // Create a test account
    account2, err := testQueries.CreateAccount(context.Background(), CreateAccountParams)
    require.NoError(t, err)

    // Define the input parameters for the CreateTransfer method
    params := CreateTransferParams{
        FromAccountID: account1.ID,
        ToAccountID:   account2.ID,
        Amount:        utils.RandomMoney(),
    }

    // Call the CreateTransfer method
    transfer, err := testQueries.CreateTransfer(context.Background(), params)
    require.NoError(t, err)

    // Retrieve the transfer from the database
    transferFromDB, err := testQueries.GetTransfer(context.Background(), transfer.ID)
    require.NoError(t, err)

    // Define the expected result
    expected := Transfer{
        ID:           transfer.ID,
        FromAccountID: params.FromAccountID,
        ToAccountID:   params.ToAccountID,
        Amount:       params.Amount,
        CreatedAt:    transfer.CreatedAt,
    }

    // Compare the result with the expected value
    require.Equal(t, expected, transferFromDB)
}

func TestGetTransfer(t *testing.T) {
	// Create a test account
	CreateAccountParams := CreateAccountParams{
		Owner:    utils.RandomOwner(),
		Balance:  utils.RandomMoney(),
		Currency: utils.RandomCurrency(),
	}
	account, err := testQueries.CreateAccount(context.Background(), CreateAccountParams)
	require.NoError(t, err)

	// Create a test transfer
	CreateTransferParams := CreateTransferParams{
		FromAccountID: account.ID,
		ToAccountID:   account.ID,
		Amount:        utils.RandomMoney(),
	}
	transfer, err := testQueries.CreateTransfer(context.Background(), CreateTransferParams)
	require.NoError(t, err)

	// Call the GetTransfer method
	transferFromDB, err := testQueries.GetTransfer(context.Background(), transfer.ID)
	require.NoError(t, err)

	// Compare the result with the expected value
	require.Equal(t, transfer, transferFromDB)
}

func TestGetTransfer_failure(t *testing.T) {
	// Call the GetTransfer method
	transferFromDB, err := testQueries.GetTransfer(context.Background(), 0)

	// Check for any errors
	require.Error(t, err)
	require.Empty(t, transferFromDB)
}

func TestListTransfers(t *testing.T) {
	// Create a test account
	CreateAccountParams := CreateAccountParams{
		Owner:    utils.RandomOwner(),
		Balance:  utils.RandomMoney(),
		Currency: utils.RandomCurrency(),
	}
	account, err := testQueries.CreateAccount(context.Background(), CreateAccountParams)
	require.NoError(t, err)

	// Create test transfers
	CreateTransferParams := CreateTransferParams{
		FromAccountID: account.ID,
		ToAccountID:   account.ID,
		Amount:        utils.RandomMoney(),
	}
	_, err = testQueries.CreateTransfer(context.Background(), CreateTransferParams)
	require.NoError(t, err)

	_, err = testQueries.CreateTransfer(context.Background(), CreateTransferParams)
	require.NoError(t, err)

	// Define the input parameters for the ListTransfers method
	params := ListTransfersParams{
		FromAccountID: account.ID,
		ToAccountID:   account.ID,
		Limit:         10,
		Offset:        0,
	}

	// Call the ListTransfers method
	transfers, err := testQueries.ListTransfers(context.Background(), params)
	require.NoError(t, err)

	// Check the number of transfers returned
	require.Len(t, transfers, 2)

	// Check the correctness of the transfers
	for _, transfer := range transfers {
		require.Equal(t, account.ID, transfer.FromAccountID)
		require.Equal(t, account.ID, transfer.ToAccountID)
	}
}

func TestListTransfersFromAccount(t *testing.T) {
	// Create a test account
	CreateAccountParams := CreateAccountParams{
		Owner:    utils.RandomOwner(),
		Balance:  utils.RandomMoney(),
		Currency: utils.RandomCurrency(),
	}
	account, err := testQueries.CreateAccount(context.Background(), CreateAccountParams)
	require.NoError(t, err)

	// Create test transfers
	CreateTransferParams := CreateTransferParams{
		FromAccountID: account.ID,
		ToAccountID:   account.ID,
		Amount:        utils.RandomMoney(),
	}
	_, err = testQueries.CreateTransfer(context.Background(), CreateTransferParams)
	require.NoError(t, err)

	_, err = testQueries.CreateTransfer(context.Background(), CreateTransferParams)
	require.NoError(t, err)

	// Define the input parameters for the ListTransfersFromAccount method
	params := ListTransfersFromAccountParams{
		FromAccountID: account.ID,
		Limit:         10,
		Offset:        0,
	}

	// Call the ListTransfersFromAccount method
	transfers, err := testQueries.ListTransfersFromAccount(context.Background(), params)
	require.NoError(t, err)

	// Check the number of transfers returned
	require.Len(t, transfers, 2)

	// Check the correctness of the transfers
	for _, transfer := range transfers {
		require.Equal(t, account.ID, transfer.FromAccountID)
		require.Equal(t, account.ID, transfer.ToAccountID)
	}
}

func TestListTransfersToAccount(t *testing.T) {
	// Create a test account
	CreateAccountParams := CreateAccountParams{
		Owner:    utils.RandomOwner(),
		Balance:  utils.RandomMoney(),
		Currency: utils.RandomCurrency(),
	}
	account, err := testQueries.CreateAccount(context.Background(), CreateAccountParams)
	require.NoError(t, err)

	// Create test transfers
	CreateTransferParams := CreateTransferParams{
		FromAccountID: account.ID,
		ToAccountID:   account.ID,
		Amount:        utils.RandomMoney(),
	}
	_, err = testQueries.CreateTransfer(context.Background(), CreateTransferParams)
	require.NoError(t, err)

	_, err = testQueries.CreateTransfer(context.Background(), CreateTransferParams)
	require.NoError(t, err)

	// Define the input parameters for the ListTransfersToAccount method
	params := ListTransfersToAccountParams{
		ToAccountID: account.ID,
		Limit:       10,
		Offset:      0,
	}

	// Call the ListTransfersToAccount method
	transfers, err := testQueries.ListTransfersToAccount(context.Background(), params)
	require.NoError(t, err)

	// Check the number of transfers returned
	require.Len(t, transfers, 2)

	// Check the correctness of the transfers
	for _, transfer := range transfers {
		require.Equal(t, account.ID, transfer.FromAccountID)
		require.Equal(t, account.ID, transfer.ToAccountID)
	}
}