package algorithm

import (
	"bytes"
	"io/ioutil"
	"sort"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

// UTF8ToGBK ...
func UTF8ToGBK(src string) ([]byte, error) {
	GB18030 := simplifiedchinese.All[0]
	return ioutil.ReadAll(transform.NewReader(bytes.NewReader([]byte(src)), GB18030.NewEncoder()))
}

// CompareString 比较两个中文字符串，用于排序
func CompareString(aa, bb string) bool {
	a, _ := UTF8ToGBK(aa)
	b, _ := UTF8ToGBK(bb)
	bLen := len(b)
	for idx, chr := range a {
		if idx >= bLen {
			return false
		}
		if chr != b[idx] {
			return chr < b[idx]
		}
	}
	return true
}

// SortStrings 按中文排序
func SortStrings(s []string) {
	sort.SliceStable(s, func(i, j int) bool {
		return CompareString(s[i], s[j])
	})
}
