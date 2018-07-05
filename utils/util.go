package utils

import (
	"strconv"
	"time"
	"strings"
)

// 判断某个条件是否被调用
func Judge(con bool, tvale, fvalue interface{}) interface{} {
	if con {
		return tvale
	}
	return fvalue

}

//  给定的参数是否在数组中（对输入参数合法性的验证）
func Has(original interface{}, target interface{}) bool {
	switch key := original.(type) {
	case string:
		for _, item := range target.([]string) {
			if key == item {
				return true
			}
		}
	case int:
		for _, item := range target.([] int) {
			if key == item {
				return true
			}
		}

	case int64:
		for _, item := range target.([] int64) {
			if key == item {
				return true
			}

		}
	default:
		return false
	}
	return false
}

//  给条件添加单引号
func AddSingleMark(data interface{}) string {
	return "'" + strings.Replace(TransStr(data), "'", `\'`, -1) + "'";
}

// 转换为string
func TransStr(data interface{}) string {
	switch data.(type) {
	case int:
		return strconv.Itoa(data.(int))
	case int64:
		return strconv.FormatInt(data.(int64), 10)
	case int32:
		return strconv.FormatInt(int64(data.(int32)), 10)
	case uint32:
		return strconv.FormatUint(uint64(data.(uint32)), 10)
	case uint64:
		return strconv.FormatUint(data.(uint64), 10)
	case float32:
		return strconv.FormatFloat(float64(data.(float32)), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(data.(float64), 'f', -1, 64)
	case string:
		return data.(string)
	case time.Time:
		return data.(time.Time).Format("2018-06-12 09:11:03")
	default:
		return ""
	}
}
func GetNowTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

//  这是一个测试
func Test() string {
	return "test"
}
