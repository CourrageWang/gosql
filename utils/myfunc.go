package utils

import "strings"

//添加执行函数
func Where(args string) (str string) {
	if args != "" {
		str = strings.Split(args, ".")[0]
		return
	}
	return ""
}
