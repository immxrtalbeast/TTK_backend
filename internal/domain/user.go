package domain

import (
	"context"
	"time"
)

type Role string

const (
	UserRole  Role = "USER"
	AdminRole Role = "ADMIN"
)

type User struct {
	ID        string
	Name      string
	Login     string
	PassHash  []byte
	CreatedAt time.Time
	IsAdmin   Role
}

type UserInteractor interface {
	CreateUser(ctx context.Context, login string, name string, pass string) error
	User(ctx context.Context, id string) (*User, error)
	Login(ctx context.Context, login string, passhash string) (string, error)
	Users(ctx context.Context, page int, limit int) ([]*User, error)
	UpdateUser(ctx context.Context, id string, name string, login string, passhash string, role Role) error
	DeleteUser(ctx context.Context, id string) error
}

type UserRepository interface {
	CreateUser(ctx context.Context, user *User) error
	User(ctx context.Context, id string) (*User, error)
	UserByLogin(ctx context.Context, login string) (*User, error)
	Users(ctx context.Context, page int, limit int) ([]*User, error)
	UpdateUser(ctx context.Context, user *User) error
	DeleteUser(ctx context.Context, id string) error
	// UpdateUserPassword(ctx context.Context, passHash []byte) error
}
