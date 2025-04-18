package db

import (
	"context"
	"database/sql"
	"fmt"
)

// store provides all functions to execute db queries individually, as well as their combinations within a transaction.
type Store struct {
	*Queries // composition instead of inheritance // by embedding Queries within Store, we can access all the methods of Queries directly on Store
	db       *sql.DB
}

// NewStore creates a new Store.
func NewStore(db *sql.DB) *Store {
	return &Store{
		Queries: New(db), // initialize Queries with the provided db connection
		db:      db,      // store the db connection
	}
}

// execTx executes a function within a database transaction.
// takes a context and a callback function as input
// it starts a new database transaction
// it creates a new queries object with that transaction
// it calls the callback function with the queries object
// finally, it commits the transaction if no error occurred, or rolls it back if an error occurred
func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil) // read-commited isolation level is the default

	if err != nil {
		return err
	}

	q := New(tx) // create a new queries object with the transaction
	err = fn(q)  // call the callback function with the queries object

	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr) // if the callback function returns an error, rollback the transaction
		}

		// rollback is successful, so we return tx error
		return err
	}

	// if all operations are successful, commit the transaction
	return tx.Commit() // return it's error to the caller
}

// TransferTxParams contains the input parameters of the transfer transaction
type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"` // ID of the account to transfer money from
	ToAccountID   int64 `json:"to_account_id"`   // ID of the account to transfer money to
	Amount        int64 `json:"amount"`
}

// TransferTxResult is the result of the transfer transaction
type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`     // the transfer record
	FromAccount Account  `json:"from_account"` // the account from which money is transferred, after balance is updated
	ToAccount   Account  `json:"to_account"`   // the account to which money is transferred, after balance is updated
	FromEntry   Entry    `json:"from_entry"`   // the entry record for the account from which money is transferred
	ToEntry     Entry    `json:"to_entry"`     // the entry record for the account to which money is transferred
}

// TransferTx performs a money transfer from 1 account to the otehr
// it creates a transfer record, add account entries, and update accounts' balance within a single tx
func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult // initialize the result variable

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})

		if err != nil {
			return err
		}

		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount, // money is moving out
		})
		if err != nil {
			return err
		}

		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount, // money is moving in
		})
		if err != nil {
			return err
		}

		// TODO: Update accounts' balance - needs locking mechanisms
		// Get account from database -> Update its balance ... proper locking mechanism is required
		account1, err := q.GetAccountForUpdate(ctx, arg.FromAccountID)
		if err != nil {
			return err
		}

		result.FromAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
			ID:      account1.ID,
			Balance: account1.Balance - arg.Amount, // money is moving out
		})
		if err != nil {
			return err
		}

		// Move money into account 2
		account2, err := q.GetAccountForUpdate(ctx, arg.ToAccountID)
		if err != nil {
			return err
		}

		result.ToAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
			ID:      account2.ID,
			Balance: account2.Balance + arg.Amount, // money is moving in
		})
		if err != nil {
			return err
		}

		return nil
	})

	return result, err
}
