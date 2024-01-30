package library

import (
	"regexp"
	"strconv"
)

// Verify 关于验证的方法都写这里
type Verify struct {
}

// IsEmail 验证邮箱
func (r *Verify) IsEmail(email string) bool {
	match, _ := regexp.MatchString("([\\w\\-]+@[\\w\\-]+\\.[\\w\\-]+)", email)
	return match
}

// IsUname 用户名应为5-17个字母或数字字符
func (r *Verify) IsUname(uname string) bool {
	match, _ := regexp.MatchString("^[A-Za-z\\d_]{5,16}$", uname)
	return match
}

// IsPassword 密码需以6-16位字母和数字的组合
func (r *Verify) IsPassword(passwd string) bool {
	match, _ := regexp.MatchString("^(?=.*[a-z,A-Z])(?=.*\\d).[a-z,A-Z,\\d]{5,16}$", passwd)
	return match
}

// IsIdcard 校验身份证
func (r *Verify) IsIdcard(idcard string) bool {
	match, _ := regexp.MatchString("(^[1-9]\\d{5}(18|19|([23]\\d))\\d{2}((0[1-9])|(10|11|12))(([0-2][1-9])|10|20|30|31)\\d{3}[0-9Xx]$)|(^[1-9]\\d{5}\\d{2}((0[1-9])|(10|11|12))(([0-2][1-9])|10|20|30|31)\\d{3}$)", idcard)
	if !match {
		return false
	} else {
		//防止使用未来的身份证注册造成困扰
		var time1 = StringToInt64(Time())
		var time2 = StringToInt64(StrToTime(idcard[6:10] + "-" + idcard[10:12] + "-" + idcard[12:14] + " 00:00:00"))

		if time1-time2 > 432000 {
			return true
		} else {
			return false
		}
	}
}

// IsRealname 校验用户姓名
func (r *Verify) IsRealname(realname string) bool {
	match, _ := regexp.MatchString("^[\\x{4e00}-\\x{9fa5}]+[·•]?[\\x{4e00}-\\x{9fa5}]+$", realname)
	return match
}

// IsNumeric 判断字符串是否为数字类型
func (r *Verify) IsNumeric(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

// IsPhone 验证手机号
func (r *Verify) IsPhone(phone string) bool {
	if !r.IsNumeric(phone) {
		return false
	}

	match, _ := regexp.MatchString("^[1](([3][0-9])|([4][5-9])|([5][0-3,5-9])|([6][5,6])|([7][0-8])|([8][0-9])|([9][1,3,5,8,9]))[0-9]{8}$", phone)
	return match
}

// IsPhoneSimple 验证手机号（简单模式）
func (r *Verify) IsPhoneSimple(phone string) bool {
	if !r.IsNumeric(phone) {
		return false
	}

	match, _ := regexp.MatchString("^[1]([3-9])[0-9]{9}$", phone)
	return match
}

// IsTime 验证时间戳
func (r *Verify) IsTime(time string) bool {
	var ttt = StringToInt64(Time())
	var time2 = StringToInt64(time) // 统一制式string输入，统一输出
	return time2 < (ttt + 3600)
}
