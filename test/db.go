package main

import (
	"github.com/CourrageWang/gosql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

var connection gosql.Connection
var err error

func init() {
	_, e := gosql.Open(map[string]string{
		"host":     "127.0.0.1", // 数据库地址
		"username": "root",      // 数据库用户名
		"password": "123456",    // 数据库密码
		"port":     "3306",      // 端口
		"database": "wyq",       // 链接的数据库名字
		"charset":  "utf8",      // 字符集
		"protocol": "tcp",       // 链接协议
		"driver":   "mysql",     // 数据库驱动(mysql,sqlite,postgres,oracle,mssql)
	})
	if e != nil {
		fmt.Println(e)
	}
}
func main()  {
	
}

func Databses_Avg() {
	db := connection.GetInstance()
	db.Query("select id,age from users where id>? limit ?", 1, 2)
}
