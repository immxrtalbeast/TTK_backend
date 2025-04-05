package user

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/immxrtalbeast/TTK_backend/internal/domain"
	"github.com/immxrtalbeast/TTK_backend/internal/lib"
	"github.com/immxrtalbeast/TTK_backend/storage/prisma/db"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type UserInteractor struct {
	userRepo  domain.UserRepository
	tokenTTL  time.Duration
	appSecret string
}

func NewUserInteractor(userRepo domain.UserRepository, tokenTTL time.Duration, appSecret string) *UserInteractor {
	return &UserInteractor{
		userRepo:  userRepo,
		tokenTTL:  tokenTTL,
		appSecret: appSecret,
	}
}

func (ui *UserInteractor) Login(ctx context.Context, login string, passhash string) (string, error) {
	const op = "uc.user.login"
	user, err := ui.userRepo.UserByLogin(ctx, login)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}
	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(passhash)); err != nil {
		return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}
	return lib.NewToken(user, ui.tokenTTL, ui.appSecret)

}
func (ui *UserInteractor) CreateUser(ctx context.Context, login string, name string, pass string) error {
	const op = "uc.user.create"
	passHash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	user := domain.User{
		Login:    login,
		Name:     name,
		PassHash: passHash,
	}
	if err := ui.userRepo.CreateUser(ctx, &user); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (ui *UserInteractor) User(ctx context.Context, id string) (*domain.User, error) {
	const op = "uc.user.get"
	user, err := ui.userRepo.User(ctx, id)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			return nil, fmt.Errorf("%s: %w", op, db.ErrNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return user, nil
}
