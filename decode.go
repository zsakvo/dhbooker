package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"
)

var (
	//IV 偏移量
	IV   = []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	test = "c1SR02T7X+xmq37zfs0U8NAj73eedAs3tnXMQKDNUPlI2vcaNRXpKA3JktMoffp3EYPCsvCjzeCJUynjDISbNP4D5HjaCp6tMrOsBBfQzVI="
)

//SHA256 sha256编码
func SHA256(data []byte) []byte {
	ret := sha256.Sum256(data)
	return ret[:]
}

//Base64Decode Base64解码
func Base64Decode(encoded string) ([]byte, error) {
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return decoded, err
	}
	return decoded, nil
}

//LoadKey 读取解密密钥
func LoadKey(EncryptKey string) []byte {
	Key := SHA256([]byte(EncryptKey))
	return Key[:32]
}

//AESDecrypt AES解密
func AESDecrypt(EncryptKey string, ciphertext []byte) []byte {
	key := LoadKey(EncryptKey)
	block, err := aes.NewCipher(key)
	check(err)
	// Generally use first 16 bytes cipher text as IV
	// in this case they use 16 bytes 0x00
	blockModel := cipher.NewCBCDecrypter(block, IV)
	plainText := make([]byte, len(ciphertext))
	blockModel.CryptBlocks(plainText, ciphertext)
	plainText = PKCS7UnPadding(plainText)
	return plainText
}

//PKCS7UnPadding 对齐
func PKCS7UnPadding(plainText []byte) []byte {
	length := len(plainText)
	unpadding := int(plainText[length-1])
	return plainText[:(length - unpadding)]
}

//GetContent 入口函数
func decode(str string, EncryptKey string) string {
	decoded, err := Base64Decode(str)
	check(err)
	raw := AESDecrypt(EncryptKey, decoded)
	return string(raw)
}
