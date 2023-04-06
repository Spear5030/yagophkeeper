package storage

import (
	"errors"
	"go.etcd.io/bbolt"
	"go.uber.org/zap"
	"log"
)

type Storage struct {
	db     *bbolt.DB
	logger *zap.Logger
}

func New(path string, lg *zap.Logger) *Storage {
	db, err := bbolt.Open(path, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	err = db.Update(func(tx *bbolt.Tx) error {
		_, errCreate := tx.CreateBucketIfNotExists([]byte("users"))
		if errCreate != nil {
			return errCreate
		}
		_, errCreate = tx.CreateBucketIfNotExists([]byte("data"))
		if errCreate != nil {
			return errCreate
		}
		return nil
	})
	return &Storage{
		db:     db,
		logger: lg,
	}
}

func (pp *Storage) RegisterUser(email string, hashedPassword []byte) (err error) {
	err = pp.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("users"))
		user := b.Get([]byte(email))
		if len(user) > 0 {
			return errors.New("user exists")
		}
		return err
	},
	)
	if err != nil {
		pp.logger.Debug("err", zap.Error(err))
		return err
	}
	err = pp.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("users"))
		err = b.Put([]byte(email), hashedPassword)
		pp.logger.Debug("err", zap.Error(err))
		return err
	})
	return err
}

func (pp *Storage) GetUserHashedPassword(email string) (hashedPassword []byte, err error) {
	err = pp.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("users"))
		hashedPassword = b.Get([]byte(email))
		if len(hashedPassword) == 0 {
			return errors.New("user no exists")
		}
		return nil
	},
	)
	return hashedPassword, err
}
