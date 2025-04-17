package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/reinhardbuyabo/simplebank/util"
	"github.com/stretchr/testify/require"
)

// a simple change to one test shouldn't affect others

// to avoid code duplication, let's write a separate function to get a random account
func createRandomAccount(t *testing.T) Account {
	arg := CreateAccountParams{
		Owner:    util.RandomOwner(), // randomly generate?
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)                          // check if there is no error // fail test if there is an error
	require.NotEmpty(t, account)                     // check if the account is not empty
	require.Equal(t, arg.Owner, account.Owner)       // check if the owner is the same
	require.Equal(t, arg.Balance, account.Balance)   // check if the balance is the same
	require.Equal(t, arg.Currency, account.Currency) // check if the currency is the same

	require.NotZero(t, account.ID)        // check that the id is automatically generated
	require.NotZero(t, account.CreatedAt) // check that the created at is automatically generated

	return account
}

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t) // call the function to create a random account
}

func TestGetAccount(t *testing.T) {
	account1 := createRandomAccount(t) // create a random account

	// call testQueries
	account2, err := testQueries.GetAccount(context.Background(), account1.ID) // get the account by id

	require.NoError(t, err)                                                        // check if there is no error
	require.NotEmpty(t, account2)                                                  // check if the account is not empty
	require.Equal(t, account1.ID, account2.ID)                                     // check if the id is the same
	require.Equal(t, account1.Owner, account2.Owner)                               // check if the owner is the same
	require.Equal(t, account1.Balance, account2.Balance)                           // check if the balance is the same
	require.Equal(t, account1.Currency, account2.Currency)                         // check if the currency is the same
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second) // check that 2 timestamps are different by at most some delta duration, in this case 1 second
}

func TestUpdateAccount(t *testing.T) {
	account1 := createRandomAccount(t) // create a random account

	// declare arguments
	arg := UpdateAccountParams{
		ID:      account1.ID,
		Balance: util.RandomMoney(), // random amount of money
	}

	account2, err := testQueries.UpdateAccount(context.Background(), arg) // update the account
	require.NoError(t, err)                                               // check if there is no error
	require.NotEmpty(t, account2)                                         // check if the account is not empty

	// compare each individual field for account 1 and 2
	require.Equal(t, account1.ID, account2.ID)                                     // check if the id is the same
	require.Equal(t, account1.Owner, account2.Owner)                               // check if the owner is the same
	require.Equal(t, account1.Currency, account2.Currency)                         // check if the currency is the same
	require.Equal(t, arg.Balance, account2.Balance)                                // check if the balance is the same
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second) // check that 2 timestamps are different by at most some delta duration, in this case 1 second
}

func TestDeleteAccount(t *testing.T) {
	account1 := createRandomAccount(t) // create a random account

	err := testQueries.DeleteAccount(context.Background(), account1.ID) // delete the account
	require.NoError(t, err)                                             // check if there is no error

	account2, err := testQueries.GetAccount(context.Background(), account1.ID) // get the account by id
	require.Error(t, err)                                                      // check if there is an error
	require.EqualError(t, err, sql.ErrNoRows.Error())                          // "sql: no rows in result set"
	require.Empty(t, account2)                                                 // check if the account is empty
}

func TestListAccounts(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomAccount(t) // create a random account
	}

	// skip the first 5 accounts and return the next 5 accounts
	arg := ListAccountsParams{
		Limit:  5, //
		Offset: 5,
	}

	accounts, err := testQueries.ListAccounts(context.Background(), arg) // list the accounts
	require.NoError(t, err)                                              // check if there is no error
	require.Len(t, accounts, 5)                                          // check if the length of the accounts is 5

	for _, account := range accounts { // loop through the accounts
		require.NotEmpty(t, account) // check if the account is not empty
	}
}
