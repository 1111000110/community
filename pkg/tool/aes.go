package tool

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
)

// 密钥和IV长度必须符合AES的要求（16、24、32字节）
var key = []byte("zhangxuan     liuxueqian") // 16字节密钥
var iv = []byte("1111100000011111")          // 16字节IV

// 加密字符串
func Encrypt(data string) (string, error) {
	if data == "" {
		return "", nil
	}
	// 创建AES加密器
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// 将字符串填充到16字节的倍数
	dataBytes := []byte(data)
	dataBytes = pad(dataBytes, aes.BlockSize)

	// 使用CBC模式加密
	ciphertext := make([]byte, len(dataBytes))
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, dataBytes)

	// 返回base64编码的密文
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// 解密字符串
func Decrypt(encryptedData string) (string, error) {
	if encryptedData == "" {
		return "", nil
	}
	// 创建AES解密器
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// 解码base64密文
	ciphertext, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return "", err
	}

	// 使用CBC模式解密
	plaintext := make([]byte, len(ciphertext))
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(plaintext, ciphertext)

	// 去除填充
	plaintext = unpad(plaintext, aes.BlockSize)

	// 返回解密后的字符串
	return string(plaintext), nil
}

// 填充数据至指定块大小
func pad(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	paddingText := make([]byte, padding)
	for i := range paddingText {
		paddingText[i] = byte(padding)
	}
	return append(data, paddingText...)
}

// 去除填充数据
func unpad(data []byte, blockSize int) []byte {
	padding := data[len(data)-1]
	return data[:len(data)-int(padding)]
}
