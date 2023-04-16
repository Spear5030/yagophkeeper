package storage

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/gob"
	"errors"
	"github.com/Spear5030/yagophkeeper/internal/domain"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"io"
	"os"
	"time"
)

const (
	TypeLoginPassword byte = 0x1
	TypeText          byte = 0x2
	TypeBinary        byte = 0x3
	TypeCard          byte = 0x4
)

type storage struct {
	filename string
	logger   *zap.Logger
	lps      []domain.LoginPassword
	fileHeaders
}

type fileHeaders struct {
	Version    int32
	CryptoAlg  int32
	UpdatedAt  time.Time
	Email      string
	HashedPass []byte
	Token      string
}

// New возвращает файловое хранилище.
// Если существует - считывает служебные данные
func New(filename string, logger *zap.Logger) (*storage, error) {
	var storage storage
	fstat, err := os.Stat(filename)
	storage.filename = filename
	storage.logger = logger

	if (errors.Is(err, os.ErrNotExist)) || (fstat.Size() == 0) {
		storage.UpdatedAt = time.Time{} //zero time

	} else {
		err = storage.readFile()
		if err != nil {
			return nil, err
		}
	}

	return &storage, nil
}

// readHeaders считывает служебные поля в структуру fileHeaders
func (s *storage) readHeaders(b []byte) error {
	buf := bytes.NewBuffer(b)
	if err := gob.NewDecoder(buf).Decode(&s.fileHeaders); err != nil {
		s.logger.Debug(err.Error())
		return err
	}
	return nil
}

// SaveUserData сохраняет данные о пользователе в структуру fileHeaders и файл
func (s *storage) SaveUserData(user domain.User, token string) error {
	var err error
	s.Email = user.Email
	s.Token = token
	s.HashedPass, err = bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	err = s.writeFile()
	if err != nil {
		return err
	}
	return nil
}

// UpdateTime сохраняет время обновления в структуру fileHeaders и файл
func (s *storage) UpdateTime() error {
	s.UpdatedAt = time.Now()
	err := s.writeFile()
	if err != nil {
		return err
	}
	return nil
}

// AddLoginPassword добавляет(через append) структуру секрета логин-пароль в файл и память
func (s *storage) AddLoginPassword(lp domain.LoginPassword) error {
	s.lps = append(s.lps, lp)
	return s.writeFile()
}

func (s *storage) makeHeaders() ([]byte, error) {
	var buf bytes.Buffer
	var err error
	if err = gob.NewEncoder(&buf).Encode(s.fileHeaders); err != nil {
		s.logger.Debug(err.Error())
		return nil, err
	}
	_, err = buf.Write([]byte{30})
	if err != nil {
		s.logger.Debug(err.Error())
		return nil, err
	}
	return buf.Bytes(), nil
}

func (s *storage) makeBody() ([]byte, error) {
	var err error
	var buf bytes.Buffer
	for _, lp := range s.lps {
		_, err = buf.Write([]byte{TypeLoginPassword})
		if err != nil {
			s.logger.Debug(err.Error())
			return nil, err
		}
		if err = gob.NewEncoder(&buf).Encode(lp); err != nil {
			s.logger.Debug(err.Error())
			return nil, err
		}
		_, err = buf.Write([]byte{30}) // record separator
		if err != nil {
			return nil, err
		}
	}
	return buf.Bytes(), nil
}

func (s *storage) writeFile() error {
	headers, err := s.makeHeaders()
	if err != nil {
		s.logger.Debug(err.Error())
		return err
	}
	body, err := s.makeBody()
	if err != nil {
		s.logger.Debug(err.Error())
		return err
	}
	full := make([]byte, 0, len(headers)+len(body))
	full = append(full, headers...)
	full = append(full, body...)
	encrypted := s.encrypt(full, "N1PCdw3M2B1TfJhoaY2mL736p2vCUc47")
	err = os.WriteFile(s.filename, encrypted, 0644)
	if err != nil {
		s.logger.Debug(err.Error())
		return err
	}
	return nil
}

func (s *storage) readFile() error {
	encrypted, err := os.ReadFile(s.filename)
	b := s.decrypt(encrypted, "N1PCdw3M2B1TfJhoaY2mL736p2vCUc47") // todo
	if err != nil {
		s.logger.Error("read file error", zap.Error(err))
	}
	buf := bytes.NewBuffer(b)
	headers, err := buf.ReadBytes(30)
	if err != nil {
		return err
	}
	err = s.readHeaders(headers)
	if err != nil {
		return err
	}
	for {
		b, err := buf.ReadBytes(30) // record separator
		if err != nil {
			if err == io.EOF {
				break
			} else {
				s.logger.Debug(err.Error())
				break
			}
		}
		var secretType byte
		if len(b) > 2 {
			secretType = b[0]
			b = b[1:]
		}
		bufferGob := bytes.NewBuffer(b)

		switch secretType {
		case TypeLoginPassword:
			lp := &domain.LoginPassword{}
			//buffer.Write(decrypted)
			bufferGob.Write(b)
			err := gob.NewDecoder(bufferGob).Decode(lp)
			if err != nil {
				s.logger.Debug(err.Error())
			}
			bufferGob.Reset()
			s.lps = append(s.lps, *lp)
		case TypeText:
		case TypeBinary:
		case TypeCard:
		default:
			continue
		}
	}
	return nil
}

// ListSecrets вывод секретов
func (s *storage) ListSecrets() []domain.LoginPassword {
	if len(s.lps) == 0 {
		return nil
	}
	return s.lps
}

// GetData Чтение всего файла секретов
func (s *storage) GetData() ([]byte, error) {
	b, err := os.ReadFile(s.filename)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// SetData Запись всего файла секретов
func (s *storage) SetData(data []byte) error {
	err := os.WriteFile(s.filename, data, 0644)
	if err != nil {
		return err
	}
	return nil
}

// GetToken возвращает токен пользователя
func (s *storage) GetToken() string {
	return s.Token
}

// GetLocalSyncTime возвращает время обновления
func (s *storage) GetLocalSyncTime() time.Time {
	return s.UpdatedAt
}

func (s *storage) encrypt(b []byte, keyString string) (encryptedBytes []byte) {

	block, err := aes.NewCipher([]byte(keyString))
	if err != nil {
		s.logger.Debug(err.Error())
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		s.logger.Debug(err.Error())
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		s.logger.Debug(err.Error())
	}
	encryptedBytes = aesGCM.Seal(nonce, nonce, b, nil)
	return encryptedBytes
}

func (s *storage) decrypt(b []byte, keyString string) (decryptedBytes []byte) {

	//key, _ := hex.DecodeString([]byte(keyString))

	block, err := aes.NewCipher([]byte(keyString))
	if err != nil {
		s.logger.Debug(err.Error())
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		s.logger.Debug(err.Error())
	}

	nonceSize := aesGCM.NonceSize()

	nonce, ciphertext := b[:nonceSize], b[nonceSize:]

	decryptedBytes, err = aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		s.logger.Debug(err.Error())
	}

	return decryptedBytes
}
