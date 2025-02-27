package database

import (
	"context"
	_ "github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"testing"
	_ "time"
)

func createRandomWallet(t *testing.T) Wallet {
	user := createRandomUser(t)

	wallet, err := testQueries.CreateWallet(context.Background(), user.ID)
	require.NoError(t, err)
	require.NotEmpty(t, wallet)

	require.Equal(t, user.ID, wallet.UserID)
	require.Equal(t, "0.00", wallet.BalanceUsd)
	require.Equal(t, "0.00", wallet.BalanceRub)
	require.Equal(t, "0.00", wallet.BalanceEur)

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
	require.Equal(t, wallet1.BalanceUsd, wallet2.BalanceUsd)
	require.Equal(t, wallet1.BalanceRub, wallet2.BalanceRub)
	require.Equal(t, wallet1.BalanceEur, wallet2.BalanceEur)
	require.Equal(t, wallet1.Uuid, wallet2.Uuid)
}

func TestUpdateUSDBalance(t *testing.T) {
	wallet := createRandomWallet(t)

	depositArg := UpdateUSDBalanceParams{
		Amount: "100.00",
		UserID: wallet.UserID,
	}

	err := testQueries.UpdateUSDBalance(context.Background(), depositArg)
	require.NoError(t, err)

	updatedWallet, err := testQueries.GetWalletByUserID(context.Background(), wallet.UserID)
	require.NoError(t, err)
	require.Equal(t, "100.00", updatedWallet.BalanceUsd)

	withdrawArg := UpdateUSDBalanceParams{
		Amount: "-50.00",
		UserID: wallet.UserID,
	}

	err = testQueries.UpdateUSDBalance(context.Background(), withdrawArg)
	require.NoError(t, err)

	finalWallet, err := testQueries.GetWalletByUserID(context.Background(), wallet.UserID)
	require.NoError(t, err)
	require.Equal(t, "50.00", finalWallet.BalanceUsd)
}

func TestUpdateRUBBalance(t *testing.T) {
	wallet := createRandomWallet(t)

	depositArg := UpdateRUBBalanceParams{
		Amount: "5000.00",
		UserID: wallet.UserID,
	}

	err := testQueries.UpdateRUBBalance(context.Background(), depositArg)
	require.NoError(t, err)

	updatedWallet, err := testQueries.GetWalletByUserID(context.Background(), wallet.UserID)
	require.NoError(t, err)
	require.Equal(t, "5000.00", updatedWallet.BalanceRub)

	withdrawArg := UpdateRUBBalanceParams{
		Amount: "-2000.00",
		UserID: wallet.UserID,
	}

	err = testQueries.UpdateRUBBalance(context.Background(), withdrawArg)
	require.NoError(t, err)

	finalWallet, err := testQueries.GetWalletByUserID(context.Background(), wallet.UserID)
	require.NoError(t, err)
	require.Equal(t, "3000.00", finalWallet.BalanceRub)
}

func TestUpdateEURBalance(t *testing.T) {
	wallet := createRandomWallet(t)

	depositArg := UpdateEURBalanceParams{
		Amount: "200.00",
		UserID: wallet.UserID,
	}

	err := testQueries.UpdateEURBalance(context.Background(), depositArg)
	require.NoError(t, err)

	updatedWallet, err := testQueries.GetWalletByUserID(context.Background(), wallet.UserID)
	require.NoError(t, err)
	require.Equal(t, "200.00", updatedWallet.BalanceEur)

	withdrawArg := UpdateEURBalanceParams{
		Amount: "-75.00",
		UserID: wallet.UserID,
	}

	err = testQueries.UpdateEURBalance(context.Background(), withdrawArg)
	require.NoError(t, err)

	finalWallet, err := testQueries.GetWalletByUserID(context.Background(), wallet.UserID)
	require.NoError(t, err)
	require.Equal(t, "125.00", finalWallet.BalanceEur)
}
