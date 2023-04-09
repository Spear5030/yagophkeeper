package storage

import (
	"bufio"
	"bytes"
	"encoding/binary"
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

const (
	VersionFile byte = 0x1
	CryptoAlg   byte = 0x1 //tmp
)

type Storage struct {
	filename   string
	logger     *zap.Logger
	lps        []domain.LoginPassword
	updatedAt  time.Time
	hashedPass []byte
}

func New(filename string, logger *zap.Logger) (*Storage, error) {
	var storage Storage
	fstat, err := os.Stat(filename)
	if (errors.Is(err, os.ErrNotExist)) || (fstat.Size() == 0) {
		err = NewFileHeaders(filename)
		if err != nil {
			logger.Error("new file error", zap.Error(err))
		}
	} else {
		file, err := os.OpenFile(filename, os.O_RDONLY, 0644)
		if err != nil {
			return nil, err
		}
		defer file.Close()
		rd := bufio.NewReader(file)
		var buffer bytes.Buffer

		headers := make([]byte, 1+1+8+254+60) //file "headers" version+cryptoalg+timestamp+email+hashedpass
		_, err = rd.Read(headers)
		if err != nil {
			logger.Error("read file error", zap.Error(err))
		}

		err = storage.readHeaders(headers)
		if err != nil {
			return nil, err
		}
		//_, err = rd.Discard(1 + 1 + 8 + 254 + 60)

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
	storage.filename = filename
	storage.logger = logger
	storage.updatedAt = time.Now()
	return &storage, nil
}

func NewFileHeaders(path string) error {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	var headers []byte
	sliceheaders := make([]byte, 8+254+60) //timestamp + email +hashedpassword
	headers = append([]byte{VersionFile, CryptoAlg}, sliceheaders...)
	_, err = file.Write(headers)
	if err != nil {
		return err
	}
	return file.Sync()
}

func (s Storage) readHeaders(headers []byte) error {
	if len(headers) != 1+1+8+254+60 {
		return errors.New("invalid file header")
	}
	//version := headers[0]
	//cryptoalg := headers[1]
	s.updatedAt = time.Unix(int64(binary.BigEndian.Uint64(headers[2:10])), 0)
	//email := strings.TrimSpace(string(headers[10:264]) //wrong trim
	s.hashedPass = headers[264:]
	//fmt.Println(bcrypt.CompareHashAndPassword(hashedPass, []byte("123456")))
	return nil
}

func (s Storage) SaveUserData(user domain.User) error {
	file, err := os.OpenFile(s.filename, os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.Seek(10, 0)
	if err != nil {
		return err
	}
	emailBytes := make([]byte, 254) // max length email
	copy(emailBytes, user.Email)
	file.Write(emailBytes) //maybe store hashed email?
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	file.Write(hashedPassword) // 60 length
	if err != nil {
		return err
	}
	return file.Sync()
}

func (s Storage) UpdateTime() error {
	file, err := os.OpenFile(s.filename, os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.Seek(2, 0)
	if err != nil {
		return err
	}
	btimestamp := make([]byte, 8)
	binary.BigEndian.PutUint64(btimestamp, uint64(time.Now().Unix()))
	_, err = file.Write(btimestamp)
	if err != nil {
		return err
	}
	s.updatedAt = time.Now()
	return file.Sync()
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
	if len(s.lps) == 0 {
		return nil
	}
	return s.lps
}
