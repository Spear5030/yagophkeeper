package usecase

import (
	"github.com/Spear5030/yagophkeeper/internal/domain"
	"go.uber.org/zap"

	"time"
)

type network interface {
	RegisterUser(user domain.User) (string, error)
	LoginUser(user domain.User) (string, error)
	CheckSync(email string) (time.Time, error)
	GetData(email string) ([]byte, error)
	SendData(email string, data []byte) error
}

type storage interface {
	ListSecrets() []domain.LoginPassword
	AddLoginPassword(domain.LoginPassword) error
	SaveUserData(user domain.User, token string) error
	UpdateTime() error
	GetData() ([]byte, error)
	SetData(data []byte) error
	GetLocalSyncTime() time.Time
}

type usecase struct {
	logger         *zap.Logger
	storage        storage
	network        network
	email          string
	serverSyncTime time.Time
	localSyncTime  time.Time
}

func New(storage storage, network network, logger *zap.Logger) *usecase {
	var uc = &usecase{
		logger:  logger,
		storage: storage,
		network: network,
	}

	uc.localSyncTime = storage.GetLocalSyncTime()
	return uc
}

func (u *usecase) ListSecrets() []domain.LoginPassword {
	return u.storage.ListSecrets()
}

func (u *usecase) AddLoginPassword(lp domain.LoginPassword) error {
	err := u.storage.AddLoginPassword(lp)
	if err != nil {
		return err
	}
	u.localSyncTime = time.Now()
	return u.storage.UpdateTime()
}

func (u *usecase) RegisterUser(user domain.User) error {
	token, err := u.network.RegisterUser(user)
	if err != nil {
		return err
	}
	err = u.storage.SaveUserData(user, token)
	u.email = user.Email
	return err
}

func (u *usecase) LoginUser(user domain.User) error {
	token, err := u.network.LoginUser(user)
	if err != nil {
		u.logger.Debug(err.Error())
		return err
	}
	u.serverSyncTime, err = u.network.CheckSync(user.Email)
	if err != nil {
		u.logger.Debug(err.Error())
		return err
	}
	err = u.storage.SaveUserData(user, token)
	u.email = user.Email
	return err
}

func (u *usecase) CheckSync() (time.Time, error) {
	t, err := u.network.CheckSync(u.email)
	u.serverSyncTime = t
	if err != nil {
		return time.Time{}, err
	}
	return t, err
}

func (u *usecase) GetLocalSyncTime() time.Time {
	return u.storage.GetLocalSyncTime()
}

func (u *usecase) SyncData() error {
	var data []byte
	var err error
	if u.serverSyncTime.IsZero() {
		_, err = u.CheckSync()
		if err != nil {
			return err
		}
	}
	if u.localSyncTime.IsZero() {
		u.localSyncTime = u.storage.GetLocalSyncTime()
	}

	if u.serverSyncTime.After(u.localSyncTime) {
		data, err = u.network.GetData(u.email)
		if err != nil {
			return err
		}
		u.storage.SetData(data)
	} else {
		data, err = u.storage.GetData()
		if err != nil {
			return err
		}
		err = u.network.SendData(u.email, data)
		if err != nil {
			return err
		}
	}
	return nil
}
