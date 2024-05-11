package internal

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/des"
	"errors"
)

type encryptFunc func(text, k []byte) ([]byte, error) // 加（解）密方法
type textEncodeFunc func(src []byte) string           // 文本加工方法
type textDecodeFunc func(s string) ([]byte, error)    // 文本逆加工方法

// EnCrypt 字符串加密结果转16进制字符（注：后一个参数为加密函数）
func EnCrypt(str, key string, textEnCB textEncodeFunc, cryptCallback encryptFunc) (string, error) {
	b, err := cryptCallback([]byte(str), []byte(key))
	if err != nil {
		return "", err
	}
	return textEnCB(b), nil
}

// DeCrypt 16进制字符转换并解密出字符串（注：后一个参数为解密函数）
func DeCrypt(str, key string, textDcCB textDecodeFunc, cryptCallback encryptFunc) ([]byte, error) {
	origData, err := textDcCB(str)
	if err != nil {
		return []byte{}, err
	}
	b, err := cryptCallback(origData, []byte(key))
	if err != nil {
		return []byte{}, err
	}
	return b, nil
}

// DesCBCDecrypt DES-EDE3-CBC 解密 8位key
func DesCBCDecrypt(origData, key []byte) ([]byte, error) {
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockMode := cipher.NewCBCDecrypter(block, key)
	cipherText := make([]byte, len(origData))
	blockMode.CryptBlocks(cipherText, origData)
	cipherText = zeroUnPadding(pkcs5UnPadding(cipherText))
	return cipherText, nil
}

// DesCBCEncrypt DES-EDE3-CBC加密 8位key
func DesCBCEncrypt(origData, key []byte) ([]byte, error) {
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}
	origData = pkcs5Padding(origData, block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, key)
	cipherText := make([]byte, len(origData))
	blockMode.CryptBlocks(cipherText, origData)
	return cipherText, nil
}

// AesECBEncrypt AES-128-ECB加密  16位key
// AES-192-ECB加密  24位key
// AES-256-ECB加密  32位key
func AesECBEncrypt(origData, key []byte) ([]byte, error) {
	if len(key) != 16 && len(key) != 24 && len(key) != 32 {
		return nil, errors.New("key length error")
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	origData = pKCS7Padding(origData, block.BlockSize())
	decrypted := make([]byte, len(origData))
	size := block.BlockSize()
	for bs, be := 0, size; bs < len(origData); bs, be = bs+size, be+size {
		block.Encrypt(decrypted[bs:be], origData[bs:be])
	}
	return decrypted, nil
}

// AesECBDecrypt AES-128-ECB解密  16位key
// AES-192-ECB解密  24位key
// AES-256-ECB解密  32位key
func AesECBDecrypt(origData, key []byte) ([]byte, error) {
	if len(key) != 16 && len(key) != 24 && len(key) != 32 {
		return nil, errors.New("key length error")
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	decrypted := make([]byte, len(origData))
	size := block.BlockSize()
	for bs, be := 0, size; bs < len(origData); bs, be = bs+size, be+size {
		block.Decrypt(decrypted[bs:be], origData[bs:be])
	}
	return pKCS7UnPadding(decrypted), nil
}

// pkcs5Padding Pkcs5Padding
func pkcs5Padding(cipherText []byte, blockSize int) []byte {
	padding := blockSize - len(cipherText)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(cipherText, padText...)
}

// pkcs5UnPadding 去掉字符串后面的填充字符
func pkcs5UnPadding(origData []byte) []byte {
	length := len(origData)
	unPadding := int(origData[length-1])
	return origData[:(length - unPadding)]
}

// pKCS7Padding PKCS7Padding
func pKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padText...)
}

// pKCS7UnPadding PKCS7UnPadding
func pKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unPadding := int(origData[length-1])
	return origData[:(length - unPadding)]
}

// zeroPadding ZeroPadding
func zeroPadding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	plaintext := bytes.Repeat([]byte{0}, padding) // 用0去填充
	return append(ciphertext, plaintext...)
}

// zeroUnPadding ZeroUnPadding
func zeroUnPadding(origData []byte) []byte {
	return bytes.TrimFunc(origData,
		func(r rune) bool {
			return r == rune(0)
		})
}
