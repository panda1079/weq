package library

import (
	"encoding/hex"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
)

// Encryption 关于传输加密的方法都写这里
type Encryption struct {
}

// RandStr 产生随机字符串
func (r *Encryption) RandStr(length int, isInt bool) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	if isInt {
		letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	}
	b := make([]rune, length)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// StrToHex 字符串转16进制
func (r *Encryption) StrToHex(str string) string {
	return hex.EncodeToString([]byte(str))
}

// HexToStr 16进制转字符串
func (r *Encryption) HexToStr(hexStringData string) string {
	hexData, _ := hex.DecodeString(hexStringData)
	if Empty(string(hexData)) {
		return hexStringData
	}
	return string(hexData)
}

// ReverseString 字符串倒叙
func (r *Encryption) ReverseString(s string) string {
	runes := []rune(s)
	for from, to := 0, len(runes)-1; from < to; from, to = from+1, to-1 {
		runes[from], runes[to] = runes[to], runes[from]
	}
	return string(runes)
}

// Enhex 16进制转加密串（二次传输加密用）
func (r *Encryption) Enhex(hex string) (str string) {

	//函数执行完成后被调用,通过 recover() 函数捕获到之前发生的错误。
	defer func(hex string) {
		if r := recover(); r != nil {
			//fmt.Println(r)
			str = hex
		}
	}(hex)

	var ValNum = [4]map[string]string{
		{"1": "g", "2": "h", "3": "i", "4": "j", "5": "k", "6": "l", "7": "m", "8": "n", "9": "o", "0": "p"},
		{"1": "g", "2": "r", "3": "s", "4": "t", "5": "u", "6": "v", "7": "w", "8": "x", "9": "y", "0": "z"},
		{"1": "G", "2": "H", "3": "I", "4": "J", "5": "K", "6": "L", "7": "M", "8": "N", "9": "O", "0": "P"},
		{"1": "Q", "2": "R", "3": "S", "4": "T", "5": "U", "6": "V", "7": "W", "8": "X", "9": "Y", "0": "Z"},
	}

	// 字符串倒序
	hex = r.ReverseString(hex)

	//在字符串前面加一个随机数字
	str = strconv.Itoa(rand.Intn(9) + 1)

	//随机数位转大写
	HexS := []rune(hex)
	for i1 := 0; i1 < len(HexS); i1++ {
		Val := string(HexS[i1])

		if (HexS[i1] >= 97 && HexS[i1] <= 122) || (HexS[i1] >= 65 && HexS[i1] <= 90) {
			if rand.Intn(2) == 1 {
				str = str + Val
			} else {
				str = str + strings.ToUpper(Val)
			}
		} else {
			str = str + ValNum[rand.Intn(4)][Val]
		}

		if rand.Intn(20) < 5 {
			str = str + strconv.Itoa(rand.Intn(8)+1)
		}
	}

	//插入随机字母
	str = str + r.RandStr(1, false)

	//遵循base64规则，补充4倍位
	for i2 := 0; i2 < (len(str) % 4); i2++ {
		str = str + "="
	}

	return str
}

// Dehex 加密串转16进制（二次传输加密用）
func (r *Encryption) Dehex(hex string) (str string) {

	//函数执行完成后被调用,通过 recover() 函数捕获到之前发生的错误。
	defer func(hex string) {
		if r := recover(); r != nil {
			//fmt.Println(r)
			str = hex
		}
	}(hex)

	var ValNum = map[string]string{
		"g": "1", "h": "2", "i": "3", "j": "4", "k": "5", "l": "6", "m": "7", "n": "8", "o": "9", "p": "0",
		"r": "2", "s": "3", "t": "4", "u": "5", "v": "6", "w": "7", "x": "8", "y": "9", "z": "0",
		"G": "1", "H": "2", "I": "3", "J": "4", "K": "5", "L": "6", "M": "7", "N": "8", "O": "9", "P": "0",
		"Q": "1", "R": "2", "S": "3", "T": "4", "U": "5", "V": "6", "W": "7", "X": "8", "Y": "9", "Z": "0",
	}

	//去除等号
	hex = strings.Replace(hex, "=", "", -1)

	//去除数字(由于首位是数字，就顺带去除了)
	hex = regexp.MustCompile(`\d+`).ReplaceAllString(hex, "")

	//去除最后一个字符
	hex = hex[0 : len(hex)-1]

	// 字符串倒序
	hex = r.ReverseString(hex)

	//字母转数字
	HexS := []rune(hex)
	for i1 := 0; i1 < len(HexS); i1++ {
		Val := string(HexS[i1])
		if _, ok := ValNum[Val]; ok {
			//如果存在key，就是数字演变的字母，即换回数字
			str = str + ValNum[Val]
		} else {
			str = str + Val
		}
	}

	// 字符转小写
	str = strings.ToLower(str)

	return str
}
