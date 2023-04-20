package storage

import (
	"fmt"
	"github.com/Spear5030/yagophkeeper/pkg/logger"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

//var appFs = afero.NewMemMapFs()

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

type mockCrypt struct {
	storage
}

func (m *mockCrypt) encrypt(b []byte) []byte {
	return b
}

func (m *mockCrypt) decrypt(b []byte) []byte {
	return b
}

func Test_write(t *testing.T) {
	lg, _ := logger.New(true)
	fst, _ := New("test", "N1PCdw3M2B1TfJhoaY2mL736p2vCUc47", lg)
	st := mockCrypt{*fst}
	fmt.Println(st.masterPass)
}
