package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strings"
)

// Md5Encode 小写
func Md5Encode(data string) string {
	// md5流程：将数据写入md5哈希对象 -> 计算哈希值 -> 将哈希值转换成十六进制的字符串并返回
	h := md5.New()                     // 创建一个md5哈希对象
	h.Write([]byte(data))              // 将数据写入到哈希对象中
	tempStr := h.Sum(nil)              // 计算哈希值，nil表示使用默认的摘要长度
	return hex.EncodeToString(tempStr) // 将md5哈希值转换成十六进制的字符串
}

// MD5Encode 大写
func MD5Encode(data string) string {
	return strings.ToUpper(Md5Encode(data))
} // 加密md5并转化为大写

// MakePassword 加密 （生成密码）
func MakePassword(plainpwd, salt string) string {
	return Md5Encode(plainpwd + salt)
} // 将明文密码和盐值相加后再用md5加密

// ValidPassword 解密 （验证密码）
func ValidPassword(plainpwd, salt string, password string) bool {
	// 验证密码流程：将明文密码和盐值相加后再用md5加密 -> 将加密后的密码和用户提供的密码进行比较
	md := Md5Encode(plainpwd + salt)
	fmt.Println(md + "    " + password)
	return md == password
}
