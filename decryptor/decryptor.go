package decryptor

import (
	"crypto/aes"
	"crypto/cipher"

	"github.com/pkg/errors"
)

type Decryptor struct {
	key    []byte
	vector []byte
}

func NewSiiDecryptor(vector []byte) *Decryptor {
	dec := Decryptor{
		key: []byte{0x2A, 0x5F, 0xCB, 0x17, 0x91, 0xD2, 0x2F, 0xB6, 0x02, 0x45, 0xB3, 0xD8, 0x36,
			0x9E, 0xD0, 0xB2, 0xC2, 0x73, 0x71, 0x56, 0x3F, 0xBF, 0x1F, 0x3C, 0x9E, 0xDF, 0x6B, 0x11, 0x82, 0x5A, 0x5D, 0x0A},
		vector: vector,
	}

	return &dec
}

func (d *Decryptor) Decrypt(data []byte) ([]byte, error) {
	trimming := func (encrypt []byte) []byte {
		padding := encrypt[len(encrypt)-1]
		return encrypt[:len(encrypt)-int(padding)]
	}

	block, err := aes.NewCipher(d.key)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errors.New("invalid crypt length")
	}
	ecb := cipher.NewCBCDecrypter(block, []byte(d.vector))
	decrypted := make([]byte, len(data))
	ecb.CryptBlocks(decrypted, data)

	return trimming(decrypted), nil
}