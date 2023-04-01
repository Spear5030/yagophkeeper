package storage

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/Spear5030/yagophkeeper/internal/domain"
	"go.uber.org/zap"
	"io"
	"os"
)

const (
	TypeLoginPassword byte = 0x1
	TypeText          byte = 0x2
	TypeBinary        byte = 0x3
	TypeCard          byte = 0x4
)

type Storage struct {
	filename string
	logger   *zap.Logger
	lps      []domain.LoginPassword
}

func New(filename string, logger *zap.Logger) (*Storage, error) {
	var storage Storage

	file, err := os.OpenFile(filename, os.O_RDONLY|os.O_APPEND|os.O_CREATE, 0644)
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
	storage.filename = filename
	storage.logger = logger
	return &storage, nil
}

func (s Storage) AddLoginPassword(lp domain.LoginPassword) error {
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
	_, err = file.Write(append(buffer.Bytes(), 30)) // record separator
	if err != nil {
		return err
	}
	return file.Sync()
}

func (s Storage) ListSecrets() []domain.LoginPassword {
	return s.lps
}

func (s Storage) RegisterUser(user domain.User) (string, error) {
	return "", nil
}

func (s Storage) LoginUser(user domain.User) (string, error) {
	return "", nil
}
