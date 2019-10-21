package helper

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"golang.org/x/crypto/bcrypt"
)

// 加密函数 实现php的password_hash()
func PasswordHash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

//验证密码 实现php的password_verify()
func PasswordVerify(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// base64 加密
func Base64Encode(s string) string {
	return base64.StdEncoding.EncodeToString([]byte(s))
}

/*因为在加密用到Base64Encode以后，替换了=为空字符串，所以这里需要处理下，PHP里面的base64_decode函数是可以直接处理，
这里自己卡了几个小时去研究了下BASE64得原理，base64是说3个字符转成4个字符的方法。因为3个字符二进制刚好24位，
分成4个，那么每个只有6位二进制，算一下刚好64，所以base64的所有字符就是64个，从A..Za..z0..9+/刚好64个。
如果当转换中，如果少了，后面补=号，所以会出现最后一个=号和两个=号的情况。这里根据这个原理反补了等号回去，
因为GO本身的解码函数没有处理这个。*/
func Base64Decode(str string) []byte {
	var b []byte
	var err error
	x := len(str) * 3 % 4
	switch {
	case x == 2:
		str += "=="
	case x == 1:
		str += "="
	}
	if b, err = base64.StdEncoding.DecodeString(str); err != nil {
		return b
	}

	return b
}

//生成32位md5字串
func GetMd5String(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

//------------------------------------------Aes

func AesEncrypt(orig string, key string) string {
	// 转成字节数组
	origData := []byte(orig)
	k := []byte(key)

	// 分组秘钥
	block, _ := aes.NewCipher(k)
	// 获取秘钥块的长度
	blockSize := block.BlockSize()
	// 补全码
	origData = PKCS7Padding(origData, blockSize)
	// 加密模式
	blockMode := cipher.NewCBCEncrypter(block, k[:blockSize])
	// 创建数组
	cryted := make([]byte, len(origData))
	// 加密
	blockMode.CryptBlocks(cryted, origData)

	return base64.StdEncoding.EncodeToString(cryted)

}

func AesDecrypt(cryted string, key string) string {
	// 转成字节数组
	crytedByte, _ := base64.StdEncoding.DecodeString(cryted)
	k := []byte(key)

	// 分组秘钥
	block, _ := aes.NewCipher(k)
	// 获取秘钥块的长度
	blockSize := block.BlockSize()
	// 加密模式
	blockMode := cipher.NewCBCDecrypter(block, k[:blockSize])
	// 创建数组
	orig := make([]byte, len(crytedByte))
	// 解密
	blockMode.CryptBlocks(orig, crytedByte)
	// 去补全码
	orig = PKCS7UnPadding(orig)
	return string(orig)
}

//补码
func PKCS7Padding(ciphertext []byte, blocksize int) []byte {
	padding := blocksize - len(ciphertext)%blocksize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

//去码
func PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}
