package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"csbackend/util"

	"github.com/stretchr/testify/require"
)

// create a random user and test if successful.
func createRandomUser(t *testing.T) User {
	hashedPassword, err := util.HashPassword(util.RandomString(6))
	require.NoError(t, err)
	arg := CreateUserParams{
		Username: sql.NullString{String: util.RandomOwner(), Valid: true},
		Password: sql.NullString{String: hashedPassword, Valid: true},
		Email:    sql.NullString{String: util.RandomEmail(), Valid: true},
	}

	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.Password, user.Password)
	require.Equal(t, arg.Email, user.Email)
	require.NotEmpty(t, user.CreatedAt)

	return user
}

// unit-test for CreateUser
func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

// unit-test for GetUser
func TestGetUser(t *testing.T) {
	user1 := createRandomUser(t)
	user2, err := testQueries.GetUser(context.Background(), user1.Username)
	require.NoError(t, err)
	require.NotEmpty(t, user2)
	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.Password, user2.Password)
	require.Equal(t, user1.Email, user2.Email)
	require.WithinDuration(t, user1.CreatedAt.Time, user2.CreatedAt.Time, time.Second)
}
