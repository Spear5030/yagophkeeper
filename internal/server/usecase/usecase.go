package usecase

import (
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type storage interface {
	RegisterUser(email string, hashedPassword []byte) (err error)
	GetUserHashedPassword(email string) (hashedPassword []byte, err error)
	GetLastSyncTime(email string) (lastSync time.Time, err error)
	SetLastSyncTime(email string, lastSync time.Time) (err error)
	SetData(email string, data []byte) (err error)
	GetData(email string) (data []byte, err error)
}

type usecase struct {
	storage   storage
	logger    *zap.Logger
	secretKey string
}

func New(s storage, lg *zap.Logger, secret string) *usecase {
	return &usecase{
		storage:   s,
		logger:    lg,
		secretKey: secret,
	}
}

func (uc *usecase) RegisterUser(email string, password string) (token string, err error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		uc.logger.Debug("genToken error", zap.Error(err))
		return "", err
	}
	err = uc.storage.RegisterUser(email, hashedPassword)
	if err != nil {
		uc.logger.Debug("Register error", zap.Error(err))
		return "", err
	}
	err = uc.storage.SetLastSyncTime(email, time.Now())
	if err != nil {
		uc.logger.Debug("Set Last Sync error", zap.Error(err))
		return "", err
	}
	if err != nil {
		uc.logger.Debug("Register error", zap.Error(err))
		return "", err
	}
	token, err = genJWT(uc.secretKey, email)
	if err != nil {
		uc.logger.Debug("genToken error", zap.Error(err))
		return "", err
	}
	return token, err
}

func (uc *usecase) LoginUser(email string, password string) (token string, err error) {
	hash, err := uc.storage.GetUserHashedPassword(email)
	if err != nil {
		return "", err
	}
	err = bcrypt.CompareHashAndPassword(hash, []byte(password))
	if err != nil {
		return "", err
	}
	token, err = genJWT(uc.secretKey, email)
	if err != nil {
		return "", err
	}
	return token, err
}

func (uc *usecase) GetLastSyncTime(email string) (lastSync time.Time, err error) {
	return uc.storage.GetLastSyncTime(email)
}

func (uc *usecase) SetData(email string, data []byte) (err error) {
	err = uc.storage.SetData(email, data)
	if err != nil {
		return err
	}
	err = uc.storage.SetLastSyncTime(email, time.Now())
	return
}

func (uc *usecase) GetData(email string) (data []byte, err error) {
	return uc.storage.GetData(email)
}

func genJWT(secretKey string, email string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(time.Hour).Unix(),
	})
	tokenString, err := token.SignedString([]byte(secretKey))
	return tokenString, err
}
