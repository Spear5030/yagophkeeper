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
	filename   string
	masterPass string
	logger     *zap.Logger
	lps        []domain.LoginPassword // TODO maps for delete
	tds        []domain.TextData
	bds        []domain.BinaryData
	cards      []domain.CardData
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
func New(filename string, masterPass string, logger *zap.Logger) (*storage, error) {
	var s storage
	fstat, err := os.Stat(filename)
	s.filename = filename
	s.logger = logger
	s.masterPass = masterPass

	if errors.Is(err, os.ErrNotExist) || fstat.Size() == 0 {
		s.UpdatedAt = time.Time{} //zero time

	} else {
		err = s.readFile()
		if err != nil {
			return nil, err
		}
	}

	return &s, nil
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

// AddLoginPassword добавляет структуру логин-пароль и записывает файл
func (s *storage) AddLoginPassword(lp domain.LoginPassword) error {
	lp.Key = len(s.lps) + 1
	s.lps = append(s.lps, lp)
	return s.writeFile()
}

func (s *storage) AddTextData(td domain.TextData) error {
	td.Key = len(s.tds) + 1
	s.tds = append(s.tds, td)
	return s.writeFile()
}

func (s *storage) AddBinaryData(bd domain.BinaryData) error {
	bd.Key = len(s.bds) + 1
	s.bds = append(s.bds, bd)
	return s.writeFile()
}

func (s *storage) AddCardData(card domain.CardData) error {
	card.Key = len(s.cards) + 1
	s.cards = append(s.cards, card)
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
	for _, bd := range s.bds {
		_, err = buf.Write([]byte{TypeBinary})
		if err != nil {
			s.logger.Debug(err.Error())
			return nil, err
		}
		if err = gob.NewEncoder(&buf).Encode(bd); err != nil {
			s.logger.Debug(err.Error())
			return nil, err
		}
		_, err = buf.Write([]byte{30}) // record separator
		if err != nil {
			return nil, err
		}
	}
	for _, td := range s.tds {
		_, err = buf.Write([]byte{TypeText})
		if err != nil {
			s.logger.Debug(err.Error())
			return nil, err
		}
		if err = gob.NewEncoder(&buf).Encode(td); err != nil {
			s.logger.Debug(err.Error())
			return nil, err
		}
		_, err = buf.Write([]byte{30}) // record separator
		if err != nil {
			return nil, err
		}
	}
	for _, card := range s.cards {
		_, err = buf.Write([]byte{TypeCard})
		if err != nil {
			s.logger.Debug(err.Error())
			return nil, err
		}
		if err = gob.NewEncoder(&buf).Encode(card); err != nil {
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
	encrypted, err := s.encrypt(full, s.masterPass)
	if err != nil {
		s.logger.Error("encrypt file error", zap.Error(err))
		return err
	}
	err = os.WriteFile(s.filename, encrypted, 0644)
	if err != nil {
		s.logger.Debug(err.Error())
		return err
	}
	return nil
}

func (s *storage) readFile() error {
	encrypted, err := os.ReadFile(s.filename)
	if err != nil {
		s.logger.Error("read file error", zap.Error(err))
		return err
	}
	b, err := s.decrypt(encrypted, s.masterPass)
	if err != nil {
		s.logger.Error("decrypt file error", zap.Error(err))
		return err
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
			bufferGob.Write(b)
			err := gob.NewDecoder(bufferGob).Decode(lp)
			if err != nil {
				s.logger.Debug(err.Error())
			}
			bufferGob.Reset()
			s.lps = append(s.lps, *lp)
		case TypeText:
			td := &domain.TextData{}
			bufferGob.Write(b)
			err := gob.NewDecoder(bufferGob).Decode(td)
			if err != nil {
				s.logger.Debug(err.Error())
			}
			bufferGob.Reset()
			s.tds = append(s.tds, *td)
		case TypeBinary:
			bd := &domain.BinaryData{}
			bufferGob.Write(b)
			err := gob.NewDecoder(bufferGob).Decode(bd)
			if err != nil {
				s.logger.Debug(err.Error())
			}
			bufferGob.Reset()
			s.bds = append(s.bds, *bd)
		case TypeCard:
			card := &domain.CardData{}
			bufferGob.Write(b)
			err := gob.NewDecoder(bufferGob).Decode(card)
			if err != nil {
				s.logger.Debug(err.Error())
			}
			bufferGob.Reset()
			s.cards = append(s.cards, *card)
		default:
			continue
		}
	}
	return nil
}

func (s *storage) GetLogins() []domain.LoginPassword {
	return s.lps
}

func (s *storage) GetTextData() []domain.TextData {
	return s.tds
}

func (s *storage) GetBinaryData() []domain.BinaryData {
	return s.bds
}

func (s *storage) GetCardsData() []domain.CardData {
	return s.cards
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

func (s *storage) encrypt(b []byte, keyString string) (encryptedBytes []byte, err error) {

	block, err := aes.NewCipher([]byte(keyString))
	if err != nil {
		s.logger.Debug(err.Error())
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		s.logger.Debug(err.Error())
		return nil, err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		s.logger.Debug(err.Error())
		return nil, err
	}
	encryptedBytes = aesGCM.Seal(nonce, nonce, b, nil)
	return encryptedBytes, nil
}

func (s *storage) decrypt(b []byte, keyString string) (decryptedBytes []byte, err error) {

	block, err := aes.NewCipher([]byte(keyString))
	if err != nil {
		s.logger.Debug(err.Error())
		return nil, err
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		s.logger.Debug(err.Error())
		return nil, err
	}

	nonceSize := aesGCM.NonceSize()

	nonce, ciphertext := b[:nonceSize], b[nonceSize:]

	decryptedBytes, err = aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		s.logger.Debug(err.Error())
		return nil, err
	}

	return decryptedBytes, nil
}
