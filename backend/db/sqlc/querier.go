// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.22.0

package db

import (
	"context"
	"database/sql"
)

type Querier interface {
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	GetUser(ctx context.Context, username sql.NullString) (User, error)
	SearchUsers(ctx context.Context, arg SearchUsersParams) ([]SearchUsersRow, error)
	UpdateUser(ctx context.Context, arg UpdateUserParams) (int32, error)
}

var _ Querier = (*Queries)(nil)
