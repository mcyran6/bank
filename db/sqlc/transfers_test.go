package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/mcyran6/bank/util"
	"github.com/stretchr/testify/require"
)

func createRandomTransfer(t *testing.T, toID sql.NullInt64, fromID sql.NullInt64) Transfer {
	arg := CreateTransferParams{
		FromAccountID: fromID,
		ToAccountID:   toID,
		Amount:        util.RandomMoney(),
	}

	transfer, err := testQueries.CreateTransfer(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)
	require.Equal(t, arg.FromAccountID, transfer.FromAccountID)
	require.Equal(t, arg.ToAccountID, transfer.ToAccountID)
	require.Equal(t, arg.Amount, transfer.Amount)
	require.NotZero(t, transfer.ID)
	require.NotZero(t, transfer.CreatedAt)
	return transfer
}

func TestCreateTransfer(t *testing.T) {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	require.NotEqual(t, account1.ID, account2.ID)
	transfer := createRandomTransfer(t, sql.NullInt64{Int64: account1.ID, Valid: true}, sql.NullInt64{Int64: account2.ID, Valid: true})
	require.Equal(t, transfer.ToAccountID.Int64, account1.ID)
	require.Equal(t, transfer.FromAccountID.Int64, account2.ID)
}

func TestGetTransfer(t *testing.T) {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	transfer1 := createRandomTransfer(t, sql.NullInt64{Int64: account1.ID, Valid: true}, sql.NullInt64{Int64: account2.ID, Valid: true})
	transfer2, err := testQueries.GetTransfer(context.Background(), transfer1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, transfer2)

	require.Equal(t, transfer1.ID, transfer2.ID)
	require.Equal(t, transfer1.FromAccountID, transfer2.FromAccountID)
	require.Equal(t, transfer1.ToAccountID, transfer2.ToAccountID)
	require.Equal(t, transfer1.Amount, transfer2.Amount)
	require.WithinDuration(t, transfer1.CreatedAt, transfer2.CreatedAt, 0)
}

func TestUpdateTransfer(t *testing.T) {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	transfer1 := createRandomTransfer(t, sql.NullInt64{Int64: account1.ID, Valid: true}, sql.NullInt64{Int64: account2.ID, Valid: true})
	arg := UpdateTransferParams{
		ID:     transfer1.ID,
		Amount: util.RandomMoney(),
	}
	transfer2, err := testQueries.UpdateTransfer(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, transfer2)

	require.Equal(t, transfer1.ID, transfer2.ID)
	require.NotEqual(t, transfer1.Amount, transfer2.Amount)
}

func TestDeleteTransfer(t *testing.T) {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	transfer1 := createRandomTransfer(t, sql.NullInt64{Int64: account1.ID, Valid: true}, sql.NullInt64{Int64: account2.ID, Valid: true})
	err := testQueries.DeleteTransfer(context.Background(), transfer1.ID)
	require.NoError(t, err)

	transfer2, err := testQueries.GetTransfer(context.Background(), transfer1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, transfer2)
}

func TestListTransfers(t *testing.T) {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	for i := 0; i < 10; i++ {
		createRandomTransfer(t, sql.NullInt64{Int64: account1.ID, Valid: true}, sql.NullInt64{Int64: account2.ID, Valid: true})
	}
	arg := ListTransfersParams{
		ToAccountID:   sql.NullInt64{Int64: account1.ID, Valid: true},
		FromAccountID: sql.NullInt64{Int64: account2.ID, Valid: true},
		Limit:         5,
		Offset:        5,
	}
	transfers, err := testQueries.ListTransfers(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, transfers, 5)

	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)
	}
}
