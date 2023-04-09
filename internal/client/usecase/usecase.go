package usecase

import (
	"github.com/Spear5030/yagophkeeper/internal/domain"
	"go.uber.org/zap"
)

type network interface {
	RegisterUser(user domain.User) error
	LoginUser(user domain.User) error
	CheckSync(email string) error
}

type storage interface {
	ListSecrets() []domain.LoginPassword
	AddLoginPassword(domain.LoginPassword) error
	SaveUserData(user domain.User) error
	UpdateTime() error
}

type usecase struct {
	logger  *zap.Logger
	storage storage
	network network
}

func New(storage storage, network network, logger *zap.Logger) *usecase {
	return &usecase{
		logger:  logger,
		storage: storage,
		network: network,
	}
}

func (u *usecase) ListSecrets() []domain.LoginPassword {
	return u.storage.ListSecrets()
}

func (u *usecase) AddLoginPassword(lp domain.LoginPassword) error {
	err := u.storage.AddLoginPassword(lp)
	if err != nil {
		return err
	}
	return u.storage.UpdateTime()
}

func (u *usecase) RegisterUser(user domain.User) (err error) {
	err = u.network.RegisterUser(user) //todo return hash?
	if err != nil {
		return err
	}
	err = u.network.CheckSync(user.Email)
	err = u.storage.SaveUserData(user)
	return err
}

func (u *usecase) LoginUser(user domain.User) (err error) {
	err = u.network.LoginUser(user)
	if err != nil {
		return err
	}
	err = u.network.CheckSync(user.Email)
	err = u.storage.SaveUserData(user)
	return err
}
