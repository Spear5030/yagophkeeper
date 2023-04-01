package usecase

import (
	"github.com/Spear5030/yagophkeeper/internal/domain"
	"go.uber.org/zap"
)

type secreter interface {
	ListSecrets() []domain.LoginPassword
	AddLoginPassword(domain.LoginPassword) error
}

type auth interface {
	RegisterUser(user domain.User) (string, error)
	LoginUser(user domain.User) (string, error)
}

type storage interface {
	secreter
	auth
}

type usecase struct {
	logger   *zap.Logger
	secreter secreter
	auth     auth
}

func New(logger *zap.Logger, storage storage) *usecase {
	return &usecase{
		logger:   logger,
		secreter: storage,
		auth:     storage,
	}
}

func (u *usecase) ListSecrets() []domain.LoginPassword {
	return u.secreter.ListSecrets()
}

func (u *usecase) AddLoginPassword(lp domain.LoginPassword) error {
	return u.secreter.AddLoginPassword(lp)
}

func (u *usecase) RegisterUser(user domain.User) (token string, err error) {
	return "", nil
}

func (u *usecase) LoginUser(user domain.User) (token string, err error) {
	return "", nil
}
