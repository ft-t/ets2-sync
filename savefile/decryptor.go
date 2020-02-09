package savefile

//https://gist.github.com/hothero/7d085573f5cb7cdb5801d7adcf66dcf3

import (
	"crypto/aes"
	"crypto/cipher"
	"errors"
)

func decryptSii(crypt []byte, key []byte, vector []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if len(crypt) == 0 {
		return nil, errors.New("invalid crypt length")
	}
	ecb := cipher.NewCBCDecrypter(block, []byte(vector))
	decrypted := make([]byte, len(crypt))
	ecb.CryptBlocks(decrypted, crypt)

	return pKCS5Trimming(decrypted), nil
}

func pKCS5Trimming(encrypt []byte) []byte {
	padding := encrypt[len(encrypt)-1]
	return encrypt[:len(encrypt)-int(padding)]
}
