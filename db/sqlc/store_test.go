package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB) // create a new store with the test database connection

	account1 := createRandomAccount(t) // create a random account for the transfer
	account2 := createRandomAccount(t) // create another random account for the transfer

	// print out balances before transactions
	fmt.Println(">> before[from:to]", account1.Balance, account2.Balance)

	// run n concurrent transfer transactions to make sure the transfer transactions work well
	n := 5              // to make it easier to debug, we should not run too many concurrent transactions
	amount := int64(10) // amount to transfer

	errs := make(chan error) // channel to receive errors from the go-routines
	// channel to receive transfer results from the go-routines
	results := make(chan TransferTxResult) // channel to receive transfer results from the go-routines

	// run n concurrent transfer transactions
	for i := 0; i < n; i++ {
		// In order to figure out why deadlocks occurred, we need to print some logs to see which transaction is calling which query, and which order
		// We have to assign a name for each transaction and pass it into the TransferTx function, via the context argument.
		txName := fmt.Sprintf("tx %d,", i+1)

		// different go-routine ,,, send back to main go-routine, using channels since they're used to connect go routines
		go func() {
			// calling the Transaction function

			// Instead of passing the background context, we will pass in the new context with the transaction name
			ctx := context.WithValue(context.Background(), txKey, txName) // after this the context in TransferTx will hold the tx name

			// it will create a new transaction and execute the transfer transaction
			result, err := store.TransferTx(ctx, TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})

			errs <- err       // send the error to the `errs` channel
			results <- result // send the result to the `results` channel
		}() // - parantheses run the go routine immediately
	}

	// Check results and errors
	existed := make(map[int]bool) // map to check if the transfer ID already exists
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

		// Check Accounts
		fromAccount := result.FromAccount             // where money is going out
		require.NotEmpty(t, fromAccount)              // check if the from account is not empty
		require.Equal(t, account1.ID, fromAccount.ID) // check if the from account ID is correct

		// Money In
		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)              // check if the to account is not empty
		require.Equal(t, account2.ID, toAccount.ID) // check if the to account ID is correct

		// Check account balance
		// print balance after the transaction
		fmt.Println(">> tx[from:to]:", fromAccount.Balance, toAccount.Balance)
		difference1 := account1.Balance - fromAccount.Balance
		difference2 := toAccount.Balance - account2.Balance
		require.Equal(t, difference1, difference2)
		require.True(t, difference1 > 0)         // check if the difference is greater than 0`
		require.True(t, difference1%amount == 0) // difference should be a multiple of the amount // amount, 2 * amount, 3 * amount, ..., n * amount

		// k
		k := int(difference1 / amount) // compute k
		// k must be int between and n, n is no. of executed txs, k must be unique for each tx
		require.True(t, k >= 1 && k <= n)  // check if k is between 1 and n
		require.NotContains(t, existed, k) // check if k is not in the existed map
		existed[k] = true                  // add k to the existed map
	}

	// Check the final updated balance
	updatedAccount1, err := testQueries.GetAccount(context.Background(), account1.ID) // get the updated account
	require.NoError(t, err)                                                           // check if there is no error

	// Updated account 2
	updatedAccount2, err := testQueries.GetAccount(context.Background(), account2.ID) // get the updated account
	require.NoError(t, err)                                                           // check if there is no error

	// print out balances after transactions
	fmt.Println(">> after:", updatedAccount1.Balance, updatedAccount2.Balance)

	// test
	require.Equal(t, account1.Balance-int64(n)*amount, updatedAccount1.Balance) // check if the updated account balance is correct
	require.Equal(t, account2.Balance+int64(n)*amount, updatedAccount2.Balance) // check if the updated account balance is correct
}
