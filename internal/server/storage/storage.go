package storage

import (
	"encoding/binary"
	"errors"
	"go.etcd.io/bbolt"
	"go.uber.org/zap"
	"log"
	"time"
)

type storage struct {
	db     *bbolt.DB
	logger *zap.Logger
}

func New(path string, lg *zap.Logger) *storage {
	db, err := bbolt.Open(path, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	err = db.Update(func(tx *bbolt.Tx) error {
		_, errCreate := tx.CreateBucketIfNotExists([]byte("users"))
		if errCreate != nil {
			return errCreate
		}
		_, errCreate = tx.CreateBucketIfNotExists([]byte("sync"))
		if errCreate != nil {
			return errCreate
		}
		_, errCreate = tx.CreateBucketIfNotExists([]byte("data"))
		if errCreate != nil {
			return errCreate
		}
		return nil
	})
	if err != nil {
		lg.Fatal("errCreate buckets", zap.Error(err))
	}
	return &storage{
		db:     db,
		logger: lg,
	}
}

func (pp *storage) RegisterUser(email string, hashedPassword []byte) (err error) {
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

func (pp *storage) GetUserHashedPassword(email string) (hashedPassword []byte, err error) {
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

func (pp *storage) GetLastSyncTime(email string) (lastSync time.Time, err error) {
	lastSync = time.Time{}
	err = pp.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("sync"))
		bytesLastSync := b.Get([]byte(email))
		if len(bytesLastSync) == 0 {
			return errors.New("no last sync time")
		}
		lastSync = time.Unix(int64(binary.BigEndian.Uint64(bytesLastSync)), 0)
		return nil
	},
	)
	return
}

func (pp *storage) SetLastSyncTime(email string, lastSync time.Time) (err error) {
	err = pp.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("sync"))
		bytesLastSync := make([]byte, 8)
		binary.BigEndian.PutUint64(bytesLastSync, uint64(lastSync.Unix()))
		err = b.Put([]byte(email), bytesLastSync)
		pp.logger.Debug("err", zap.Error(err))
		return err
	},
	)
	return err
}

func (pp *storage) SetData(email string, data []byte) (err error) {
	err = pp.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("data"))
		err = b.Put([]byte(email), data)
		if err != nil {
			pp.logger.Debug("err", zap.Error(err))
		}
		return err
	},
	)
	return
}

func (pp *storage) GetData(email string) (data []byte, err error) {
	err = pp.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("data"))
		data = b.Get([]byte(email))
		if len(data) == 0 {
			return errors.New("no data")
		}
		return err
	},
	)
	return
}
