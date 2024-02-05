package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store interface {
    Querier
    TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
}

// SQLStore provides all functions to execute SQL queries and transactions
type SQLStore struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}

func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}
	return tx.Commit()
}


type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

type AddMoneyParams struct {
    ID int64
    accountID1 int64
    amount1 int64
    accountID2 int64
    amount2 int64
}


// TransferTx performs a money transfer from one account to the other.
// It creates a transfer record, add accoubnt entry, and update accounts' balance with new values within a single database transaction.
// If any of the operations fail, it returns an error.
func (store *SQLStore) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
    var result TransferTxResult

    err := store.execTx(ctx, func(q *Queries) error {
        var err error


        // Step 1: Create a new entry in the transfers table
        result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
            FromAccountID: arg.FromAccountID,
            ToAccountID:   arg.ToAccountID,
            Amount:        arg.Amount,
        })
        if err != nil {
            return err
        }

        // Step 2: Create entries in the account_entries table
        result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
            AccountID: arg.FromAccountID,
            Amount:    -arg.Amount,
        })
        if err != nil {
            return err
        }

        result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
            AccountID: arg.ToAccountID,
            Amount:    arg.Amount,
        })
        if err != nil {
            return err
        }

        if arg.FromAccountID < arg.ToAccountID {

            result.FromAccount, result.ToAccount, err = addMoney(ctx, q, AddMoneyParams{
                accountID1: arg.FromAccountID,
                amount1: -arg.Amount,
                accountID2: arg.ToAccountID,
                amount2: arg.Amount,
            })

            if err != nil {
                return err
            }
            
        } else {
            result.ToAccount, result.FromAccount, err = addMoney(ctx, q, AddMoneyParams{
                accountID1: arg.ToAccountID,
                amount1: arg.Amount,
                accountID2: arg.FromAccountID,
                amount2: -arg.Amount,
            })

            if err != nil {
                return err
            }
        }

        return nil
    })

    return result, err
}


func addMoney(ctx context.Context, q *Queries, arg AddMoneyParams) (account1 Account, account2 Account, err error) {
    account1, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
        ID:      arg.accountID1,
        Amount: arg.amount1,
    })

    if err != nil {
        return
    }

    account2, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
        ID:      arg.accountID2,
        Amount: arg.amount2,
    })
    if err != nil {
        return
    }

    return

}