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
}

type usecase struct {
	storage storage
	logger  *zap.Logger
}

func New(s storage, lg *zap.Logger) *usecase {
	return &usecase{
		storage: s,
		logger:  lg,
	}
}

func (uc *usecase) RegisterUser(email string, password string) (token string, err error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	err = uc.storage.RegisterUser(email, hashedPassword)
	token, err = genJWT("secret") //todo cfg
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
	token, err = genJWT("secret") //todo cfg
	if err != nil {
		return "", err
	}
	return token, err
}

func genJWT(secretKey string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": time.Now().Add(time.Hour).Unix(),
	})
	tokenString, err := token.SignedString([]byte(secretKey))
	return tokenString, err
}
