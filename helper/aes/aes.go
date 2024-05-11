package aes

import (
	"encoding/base64"
	"encoding/hex"

	"git.woa.com/kf_cdms/go-public/helper/aes/internal"
)

// Base64Encrypt AES 16位 秘钥加密
func Base64Encrypt(str string, key string) (string, error) {
	return internal.EnCrypt(str, key, base64.StdEncoding.EncodeToString, internal.AesECBEncrypt)
}

// Base64Decrypt AES 16位 秘钥解密
func Base64Decrypt(str string, key string) ([]byte, error) {
	return internal.DeCrypt(str, key, base64.StdEncoding.DecodeString, internal.AesECBDecrypt)
}

// HexEncrypt aes搭配hex加密
func HexEncrypt(str string, key string) (string, error) {
	return internal.EnCrypt(str, key, hex.EncodeToString, internal.AesECBEncrypt)
}

// HexDecrypt aes搭配hex解密
func HexDecrypt(str string, key string) ([]byte, error) {
	return internal.DeCrypt(str, key, hex.DecodeString, internal.AesECBDecrypt)
}
