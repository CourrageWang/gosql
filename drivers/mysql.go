package drivers

import "fmt"

//  mysql 驱动

/**  xx:xx@tcp(127.0.0.1:3306)/mvbox?charset=utf8
     需要信息 ：  用户名 、密码、协议、 地址、 端口 、 数据库名 、 字符集
     dirver : 驱动类型 、 dco :数据库连接对象 ,Dbo : 数据库对象
 */
func Mysql(dbO map[string]string) (driver string, dco string) {
	driver = "mysql"
	dco = fmt.Sprintf("%s:%s@%s(%s:%s)/%s?charset=%s", dbO["username"], dbO["password"], dbO["protocol"],
		dbO["host"], dbO["port"], dbO["database"], dbO["charset"])
	return
}