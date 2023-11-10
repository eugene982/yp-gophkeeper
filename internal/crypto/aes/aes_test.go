package aes

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {

	t.Run("error", func(t *testing.T) {
		_, err := New([]byte("err"))
		require.Error(t, err) // invalid key size 3
	})

	t.Run("ok", func(t *testing.T) {
		b, err := New([]byte("key-size-16bytes"))
		require.NoError(t, err)
		require.NotNil(t, b)
	})
}

func TestEncrypt(t *testing.T) {

	t.Run("error decrypt", func(t *testing.T) {
		crypt, err := New([]byte("key-size-16bytes"))
		require.NoError(t, err)

		_, err = crypt.Decrypt([]byte("some--string"))
		require.Error(t, err)
	})

	t.Run("ecrypt decrypt", func(t *testing.T) {
		crypt, err := New([]byte("key-size-16bytes"))
		require.NoError(t, err)

		str := []byte("some--string")
		b, err := crypt.Encrypt(str)
		require.NoError(t, err)

		res, err := crypt.Decrypt(b)
		require.NoError(t, err)
		require.Equal(t, str, res)
	})
}
