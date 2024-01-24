package box

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"github.com/sagernet/sing-box/constant"
)

func DecryptAES(ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher([]byte(constant.ENCRYPT_KEY))
	mode := cipher.NewCBCDecrypter(block, []byte(constant.ENCRYPT_KEY_IV))
	plaintext := make([]byte, len(ciphertext))
	ciphertext, err = base64.StdEncoding.DecodeString(string(ciphertext))
	mode.CryptBlocks(plaintext, ciphertext)
	padSize := int(plaintext[len(plaintext)-1])
	return plaintext[:len(plaintext)-padSize], err
}
