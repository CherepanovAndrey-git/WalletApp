package database

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func createRandomWallet(t *testing.T) Wallet {
	user := createRandomUser(t)

	arg := CreateWalletParams{
		UserID:  user.ID,
		Balance: "0.00",
	}

	wallet, err := testQueries.CreateWallet(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, wallet)

	require.Equal(t, arg.UserID, wallet.UserID)
	require.Equal(t, arg.Balance, wallet.Balance)

	require.NotZero(t, wallet.ID)
	require.NotZero(t, wallet.Uuid)

	return wallet
}

func TestCreateWallet(t *testing.T) {
	createRandomWallet(t)
}

func TestGetWalletByUserID(t *testing.T) {
	wallet1 := createRandomWallet(t)
	wallet2, err := testQueries.GetWalletByUserID(context.Background(), wallet1.UserID)
	require.NoError(t, err)
	require.NotEmpty(t, wallet2)

	require.Equal(t, wallet1.ID, wallet2.ID)
	require.Equal(t, wallet1.UserID, wallet2.UserID)
	require.Equal(t, wallet1.Balance, wallet2.Balance)
	require.Equal(t, wallet1.Uuid, wallet2.Uuid)
}

func TestUpdateWalletBalance(t *testing.T) {
	wallet := createRandomWallet(t)

	// Test deposit
	depositArg := UpdateWalletBalanceParams{
		Amount: "100.00",
		Userid: wallet.UserID,
	}

	err := testQueries.UpdateWalletBalance(context.Background(), depositArg)
	require.NoError(t, err)

	updatedWallet, err := testQueries.GetWalletByUserID(context.Background(), wallet.UserID)
	require.NoError(t, err)
	require.Equal(t, "100.00", updatedWallet.Balance)

	// Test withdrawal
	withdrawArg := UpdateWalletBalanceParams{
		Amount: "-50.00",
		Userid: wallet.UserID,
	}

	err = testQueries.UpdateWalletBalance(context.Background(), withdrawArg)
	require.NoError(t, err)

	finalWallet, err := testQueries.GetWalletByUserID(context.Background(), wallet.UserID)
	require.NoError(t, err)
	require.Equal(t, "50.00", finalWallet.Balance)
}
