package usecase

import (
	"fmt"
	"github.com/Spear5030/yagophkeeper/internal/domain"
	"go.uber.org/zap"
	"time"
)

type network interface {
	RegisterUser(user domain.User) (string, error)
	LoginUser(user domain.User) (string, error)
	CheckSync(email string) (time.Time, error)
	GetData() ([]byte, error)
	SendData(data []byte) error
}

//go:generate mockery --name "storage"
type storage interface {
	GetLogins() []domain.LoginPassword
	GetTextData() []domain.TextData
	GetBinaryData() []domain.BinaryData
	GetCardsData() []domain.CardData
	AddLoginPassword(domain.LoginPassword) error
	AddTextData(domain.TextData) error
	AddBinaryData(domain.BinaryData) error
	AddCardData(domain.CardData) error
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
	version        string
	buildTime      string
}

func New(storage storage, network network, version string, buildTime string, logger *zap.Logger) *usecase {
	var uc = &usecase{
		logger:    logger,
		storage:   storage,
		network:   network,
		version:   version,
		buildTime: buildTime,
	}

	uc.localSyncTime = storage.GetLocalSyncTime()
	return uc
}

func (u usecase) GetLoginsPasswords() []domain.LoginPassword {
	return u.storage.GetLogins()
}

func (u *usecase) ListSecrets() []string {
	var result []string
	result = append(result, "Logins:")
	for _, l := range u.storage.GetLogins() {
		result = append(result, fmt.Sprintf("Key[%d],%s:%s", l.Key, l.Login, l.Password))
	}
	result = append(result, "Texts:")
	for _, txt := range u.storage.GetTextData() {
		result = append(result, fmt.Sprintf("Key[%d],%s", txt.Key, txt.Text))
	}
	result = append(result, "Cards:")
	for _, card := range u.storage.GetCardsData() {
		result = append(result, fmt.Sprintf("Key[%d],%s,%s,%s", card.Key, card.Number, card.CVC, card.CardHolder))
	}
	result = append(result, "Binary:")
	for _, b := range u.storage.GetBinaryData() {
		result = append(result, fmt.Sprintf("Key[%d],%s,%b", b.Key, b.Meta, b.BinaryData[:20]))
	}
	return result
}

func (u *usecase) AddLoginPassword(lp domain.LoginPassword) error {
	err := u.storage.AddLoginPassword(lp)
	if err != nil {
		return err
	}
	u.localSyncTime = time.Now()
	return u.storage.UpdateTime()
}

func (u *usecase) AddTextData(td domain.TextData) error {
	err := u.storage.AddTextData(td)
	if err != nil {
		return err
	}
	u.localSyncTime = time.Now()
	return u.storage.UpdateTime()
}

func (u *usecase) AddBinaryData(bd domain.BinaryData) error {
	err := u.storage.AddBinaryData(bd)
	if err != nil {
		return err
	}
	u.localSyncTime = time.Now()
	return u.storage.UpdateTime()
}

func (u *usecase) AddCardData(card domain.CardData) error {
	err := u.storage.AddCardData(card)
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
	tSync, err := u.network.CheckSync(user.Email)
	if err != nil {
		u.logger.Debug(err.Error())
		return err
	}
	u.serverSyncTime = tSync
	err = u.storage.SaveUserData(user, token)
	u.email = user.Email
	return err
}

func (u *usecase) CheckSync() (time.Time, error) {
	t, err := u.network.CheckSync(u.email)
	if err != nil {
		return time.Time{}, err
	}
	u.serverSyncTime = t
	return t, err
}

func (u *usecase) GetLocalSyncTime() time.Time {
	return u.storage.GetLocalSyncTime()
}

// SyncData сравнивает время обновления на сервере и локально. Синхронизирует файлы секретов на сервере и локально
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
		data, err = u.network.GetData()
		if err != nil {
			return err
		}
		u.storage.SetData(data)
	} else {
		data, err = u.storage.GetData()
		if err != nil {
			return err
		}
		err = u.network.SendData(data)
		if err != nil {
			return err
		}
	}
	return nil
}

func (u *usecase) GetVersion() string {
	return u.version
}

func (u *usecase) GetBuildTime() string {
	return u.buildTime
}
