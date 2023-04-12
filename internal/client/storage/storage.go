package storage

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
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

func New(filename string, logger *zap.Logger) (*storage, error) {
	var storage storage
	fstat, err := os.Stat(filename)
	storage.filename = filename
	if (errors.Is(err, os.ErrNotExist)) || (fstat.Size() == 0) {
		storage.UpdatedAt = time.Time{} //zero time
		err = storage.writeHeaders()
		if err != nil {
			logger.Error("new file error", zap.Error(err))
		}
	} else {
		err = storage.readHeaders()
		if err != nil {
			logger.Error("new file error", zap.Error(err))
		}
		file, err := os.OpenFile(filename, os.O_RDONLY, 0644)
		if err != nil {
			return nil, err
		}
		defer file.Close()
		rd := bufio.NewReader(file)
		var buffer bytes.Buffer
		for {
			b, err := rd.ReadBytes(30) // record separator
			if err != nil {
				if err == io.EOF {
					break
				} else {
					fmt.Println(err)
					break
				}
			}
			var secretType byte
			if len(b) > 1 {
				secretType = b[0]
				b = b[1:]
			}
			switch secretType {
			case TypeLoginPassword:
				lp := &domain.LoginPassword{}
				buffer.Write(b)
				err := gob.NewDecoder(&buffer).Decode(lp)
				if err != nil {
					fmt.Println(err)
				}
				buffer.Reset()
				storage.lps = append(storage.lps, *lp)
			case TypeText:
			case TypeBinary:
			case TypeCard:
			default:
				continue
			}
		}
	}
	storage.logger = logger
	return &storage, nil
}

func (s *storage) writeHeaders() error {

	file, err := os.OpenFile(s.filename, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		s.logger.Debug(err.Error())
		return err
	}
	defer file.Close()
	var buffer bytes.Buffer
	//var tmp fileHeaders
	if err = gob.NewEncoder(&buffer).Encode(s.fileHeaders); err != nil {
		s.logger.Debug(err.Error())
		return err
	}

	_, err = file.Write(append(buffer.Bytes(), 30)) // record separator
	if err != nil {
		return err
	}
	return file.Sync()
}

func (s *storage) readHeaders() error {
	file, err := os.OpenFile(s.filename, os.O_RDONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	var buffer bytes.Buffer
	rd := bufio.NewReader(file)
	b, err := rd.ReadBytes(30)
	if err != nil {
		s.logger.Debug("error readHeaders:" + err.Error())
	}
	buffer.Write(b)
	if err = gob.NewDecoder(&buffer).Decode(&s.fileHeaders); err != nil {
		s.logger.Debug(err.Error())
		return err
	}
	return nil
}

func (s *storage) SaveUserData(user domain.User, token string) error {
	var err error
	s.Email = user.Email
	s.Token = token
	s.HashedPass, err = bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	err = s.writeHeaders()
	if err != nil {
		return err
	}
	return nil
}

func (s *storage) UpdateTime() error {
	s.UpdatedAt = time.Now()
	err := s.writeHeaders()
	if err != nil {
		return err
	}
	return nil
}

func (s *storage) AddLoginPassword(lp domain.LoginPassword) error {
	s.lps = append(s.lps, lp)
	file, err := os.OpenFile(s.filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	var buffer bytes.Buffer
	if err = gob.NewEncoder(&buffer).Encode(lp); err != nil {
		fmt.Println(err)
		return err
	}

	_, err = file.Write([]byte{TypeLoginPassword})
	if err != nil {
		return err
	}
	_, err = file.Write(append(buffer.Bytes(), 30)) // record separator
	if err != nil {
		return err
	}
	return file.Sync()
}

func (s *storage) ListSecrets() []domain.LoginPassword {
	if len(s.lps) == 0 {
		return nil
	}
	return s.lps
}

func (s *storage) GetData() ([]byte, error) {
	b, err := os.ReadFile(s.filename)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (s *storage) SetData(data []byte) error {
	err := os.WriteFile(s.filename, data, 0644)
	if err != nil {
		return err
	}
	return nil
}

func (s *storage) GetToken() string {
	return s.Token
}

func (s *storage) GetLocalSyncTime() time.Time {
	return s.UpdatedAt
}
