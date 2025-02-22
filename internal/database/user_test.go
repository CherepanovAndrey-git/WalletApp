package database

//
//import (
//	"context"
//	"testing"
//	"time"
//
//	"github.com/stretchr/testify/require"
//	"wallet-app/internal/utils"
//)
//
//func createRandomUser(t *testing.T) User {
//	arg := CreateUserParams{
//		ID:           utils.RandomUUID(),
//		Username:     utils.RandomOwner(),
//		Email:        utils.RandomEmail(),
//		PasswordHash: "hashedpassword",
//		CreatedAt:    time.Now(),
//		UpdatedAt:    time.Now(),
//	}
//
//	user, err := testQueries.CreateUser(context.Background(), arg)
//	require.NoError(t, err)
//	require.NotEmpty(t, user)
//
//	require.Equal(t, arg.Username, user.Username)
//	require.Equal(t, arg.Email, user.Email)
//	require.Equal(t, arg.PasswordHash, user.PasswordHash)
//
//	require.NotZero(t, user.ID)
//	require.NotZero(t, user.CreatedAt)
//
//	return user
//}
//
//func TestCreateUser(t *testing.T) {
//	createRandomUser(t)
//}
//
//func TestGetUserByEmail(t *testing.T) {
//	user1 := createRandomUser(t)
//	user2, err := testQueries.GetUserByEmail(context.Background(), user1.Email)
//	require.NoError(t, err)
//	require.NotEmpty(t, user2)
//
//	require.Equal(t, user1.ID, user2.ID)
//	require.Equal(t, user1.Username, user2.Username)
//	require.Equal(t, user1.Email, user2.Email)
//	require.Equal(t, user1.PasswordHash, user2.PasswordHash)
//	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
//}
