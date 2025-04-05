package domain

import (
	"context"
)

type Role string

const (
	UserRole  Role = "USER"
	AdminRole Role = "ADMIN"
)

type User struct {
	ID       string
	Name     string
	Login    string
	PassHash []byte
	IsAdmin  Role
}

type UserInteractor interface {
	CreateUser(ctx context.Context, login string, name string, pass string) error
	User(ctx context.Context, id string) (*User, error)
	Login(ctx context.Context, login string, passhash string) (string, error)
}

type UserRepository interface {
	CreateUser(ctx context.Context, user *User) error
	User(ctx context.Context, id string) (*User, error)
	UserByLogin(ctx context.Context, login string) (*User, error)
	// UpdateUser(ctx context.Context, user *User) error
	// UpdateUserPassword(ctx context.Context, passHash []byte) error
}
