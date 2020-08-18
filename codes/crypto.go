package codes

import (
	"bytes"
	"crypto/cipher"
	"crypto/des"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"sort"
	"strings"
)

func ToMd5Hex(str string) string {
	h := md5.New()
	h.Write([]byte(str)) // 需要加密的字符串
	cipherStr := h.Sum(nil)
	return hex.EncodeToString(cipherStr) // 输出加密结果
}

/**
  连接变量名和变量值，并按变量名进行排序
*/
func JoinAndSort(p map[string]string, contact, split string) string {
	var s, keys []string
	for k, _ := range p {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		s = append(s, k+contact+p[k])
	}

	return strings.Join(s, split)
}

//----------------------------3DES------------------------------//
func padding(src []byte, blocksize int) []byte {
	padnum := blocksize - len(src)%blocksize
	pad := bytes.Repeat([]byte{byte(padnum)}, padnum)
	return append(src, pad...)
}

func unpadding(src []byte) []byte {
	n := len(src)
	unpadnum := int(src[n-1])
	return src[:n-unpadnum]
}

type ThreeDES struct {
	block cipher.Block
	key   []byte
}

func (this *ThreeDES) SetKey(s string) error {
	var err error
	this.key = []byte(s)
	this.block, err = des.NewTripleDESCipher(this.key)
	if err != nil {
		return err
	}
	return nil
}

func (this *ThreeDES) Encrypt(src []byte) (string, error) {
	src = padding(src, this.block.BlockSize())
	blockmode := cipher.NewCBCEncrypter(this.block, this.key[:this.block.BlockSize()])
	crypted := make([]byte, len(src))
	blockmode.CryptBlocks(crypted, src)
	return base64.StdEncoding.EncodeToString(crypted), nil
}

func (this *ThreeDES) Decrypt(src string) ([]byte, error) {
	s, err := base64.StdEncoding.DecodeString(src)
	if err != nil {
		return nil, err
	}
	blockmode := cipher.NewCBCDecrypter(this.block, this.key[:this.block.BlockSize()])
	origData := make([]byte, len(s))
	blockmode.CryptBlocks(origData, s)
	return unpadding(origData), nil
}
