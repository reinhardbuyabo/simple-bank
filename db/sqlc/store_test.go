package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB) // create a new store with the test database connection

	account1 := createRandomAccount(t) // create a random account for the transfer
	account2 := createRandomAccount(t) // create another random account for the transfer

	// run n concurrent transfer transactions to make sure the transfer transactions work well
	n := 5
	amount := int64(10) // amount to transfer

	errs := make(chan error) // channel to receive errors from the go-routines
	// channel to receive transfer results from the go-routines
	results := make(chan TransferTxResult) // channel to receive transfer results from the go-routines

	// run n concurrent transfer transactions
	for i := 0; i < n; i++ {
		// different go-routine ,,, send back to main go-routine, using channels since they're used to connect go routines
		go func() {
			// calling the Transaction function
			// it will create a new transaction and execute the transfer transaction
			result, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})

			errs <- err       // send the error to the `errs` channel
			results <- result // send the result to the `results` channel
		}() // - parantheses run the go routine immediately
	}

	// Check results and errors
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err) // check if there is no error

		result := <-results
		require.NotEmpty(t, result) // check if the result is not empty

		// Check Transfer Fields
		transfer := result.Transfer
		require.NotEmpty(t, transfer)                         // check if the transfer is not empty
		require.Equal(t, account1.ID, transfer.FromAccountID) // check if the from account ID is correct
		require.Equal(t, account2.ID, transfer.ToAccountID)   // check if the to account ID is correct
		require.Equal(t, amount, transfer.Amount)             // check if the amount is correct
		require.NotZero(t, transfer.ID)                       // check if the transfer ID is not zero
		require.NotZero(t, transfer.CreatedAt)                // check if the transfer created at is not zero

		// get the transfer [READ]
		_, err = store.GetTransfer(context.Background(), transfer.ID) // Queries object is embedded within the store, so the queries object is in the store
		require.NoError(t, err)                                       // check if there is no error

		// Check Entries
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)                     // check if the from entry is not empty
		require.Equal(t, account1.ID, fromEntry.AccountID) // check if the from account ID is correct
		require.Equal(t, -amount, fromEntry.Amount)        // check if the from entry amount is correct
		require.NotZero(t, fromEntry.ID)                   // check if the from entry ID is not zero
		require.NotZero(t, fromEntry.CreatedAt)            // check if the from entry created at is not zero

		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err) // check if there is no error

		// To Entry
		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)                     // check if the to entry is not empty
		require.Equal(t, account2.ID, toEntry.AccountID) // check if the to account ID is correct
		require.Equal(t, amount, toEntry.Amount)         // check if the to entry amount is correct
		require.NotZero(t, toEntry.ID)                   // check if the to entry ID is not zero
		require.NotZero(t, toEntry.CreatedAt)            // check if the to entry created at is not zero

		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err) // check if there is no error
	}
}
