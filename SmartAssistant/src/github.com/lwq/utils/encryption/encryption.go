package encryption

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
)

func EncryptAES(plaintext, key string) (string, error) {
	// 将密钥转换为字节数组
	keyBytes := []byte(key)

	// 创建一个新的 AES 加密块
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", err
	}

	// 填充明文数据
	plaintextBytes := []byte(plaintext)
	blockSize := block.BlockSize()
	padding := blockSize - len(plaintextBytes)%blockSize
	paddingBytes := bytes.Repeat([]byte{byte(padding)}, padding)
	paddedPlaintextBytes := append(plaintextBytes, paddingBytes...)

	// 创建一个密码分组模式（CBC）
	iv := make([]byte, blockSize)
	mode := cipher.NewCBCEncrypter(block, iv)

	// 加密数据
	ciphertext := make([]byte, len(paddedPlaintextBytes))
	mode.CryptBlocks(ciphertext, paddedPlaintextBytes)

	// 返回加密后的字符串
	return base64.URLEncoding.EncodeToString(ciphertext), nil
}
