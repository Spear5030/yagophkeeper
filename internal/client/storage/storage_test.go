package storage

import (
	"github.com/Spear5030/yagophkeeper/internal/domain"
	"github.com/Spear5030/yagophkeeper/pkg/logger"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	lg, _ := logger.New(true)
	fst, err := New("test", "12345", lg)
	require.NoError(t, err)
	require.Equal(t, time.Time{}, fst.UpdatedAt)
}

func TestEncryptDecrypt(t *testing.T) {
	lg, _ := logger.New(true)
	fst, _ := New("test", "N1PCdw3M2B1TfJhoaY2mL736p2vCUc47", lg)
	testData := []byte("very secret data")
	enc, err := fst.encrypt(testData)
	require.NoError(t, err)
	dec, err := fst.decrypt(enc)
	require.NoError(t, err)
	require.Equal(t, testData, dec)
}

func TestWriteRead(t *testing.T) {
	lg, _ := logger.New(true)
	appFs = afero.NewMemMapFs()
	fst, _ := New("test", "N1PCdw3M2B1TfJhoaY2mL736p2vCUc47", lg)
	fst.Email = "test@test.ts"
	err := fst.AddLoginPassword(domain.LoginPassword{
		Key:      1,
		Login:    "atata",
		Password: "dsada",
	})
	require.NoError(t, err)
	err = fst.AddTextData(domain.TextData{
		Key:  1,
		Text: "test\ntext",
	})
	require.NoError(t, err)
	fst2, _ := New("test", "N1PCdw3M2B1TfJhoaY2mL736p2vCUc47", lg)
	require.Equal(t, fst, fst2)
}
