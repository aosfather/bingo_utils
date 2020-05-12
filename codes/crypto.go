package codes

import (
	"crypto/md5"
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
