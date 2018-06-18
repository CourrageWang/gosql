package main

import (
	"github.com/CourrageWang/gosql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

var co gosql.Connection
var e error

func main() {
	co, e = gosql.Open(map[string]string{
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
	defer co.Close()
	db := co.GetInstance()
	//res, err := db.Table("users").Fileds("id,age").Where("id", ">", 2).Where("id", "<", 4).Avg("age")
	//res, err := db.Table("users").Fileds("id, age").Where("id", ">", 2).OrWhere("id","=",1).Get()
	//res, err := db.Table("users  a ").Rightjoin("card b ", "a.id", "=", "b.user_id").Get()
	//res, err := db.Execute("update users set age = ? where id = ? ", 100,4)
	//res, err := db.Execute("update users set age = 100 where id = 4 ")
	//res, err := db.Execute("delete  from users  where id = 3")
	/*res, err := db.Execute("insert into users (id ,age) values (8,9)")
	fmt.Println("----->" ,db.RowsAffected)
	if err != nil {

		fmt.Println("FAIL: test failed.", err)
		return
	}

	if res == 0 {
		fmt.Println("FAIL: test failed.", err)
	} else {
		fmt.Println("PASS: avg=", res)
	}*/
	//data := map[string]interface{}{
	//	"age":  15,
	//	"name": "王勇琪12",
	//}
	//where := map[string]interface{}{
	//	"id": 4,
	//}
	//
	//res, err := db.Table("users").Data(data).Where(where).Update()
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//
	//fmt.Println(res)

	/*data := map[string]interface{}{
		"id":  9,
		"age":  12,
		"name": "z",
	}
	res, err := db.Table("users").Data(data).Insert()
	fmt.Println(res)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("RowsAffected: %d \n", db.RowsAffected)
	fmt.Printf("LastInsertId: %d", db.LastInsertId)

*/
	//where := map[string]interface{}{
	//	"id": 17,
	//}
	//res, err := db.Table("users").Where(where).Delete()
	//res, err := db.Table("users").Where("id", ">", 2).Sum("age")
	//res, err := db.Table("users").Where("id", ">", 2).Min("age")
	res, err := db.Table("users").Fileds("age,id,count(age) as sum").Where("id", ">", 0).OrWhere("age", ">", 0).
		Group("age").Order("id desc").Having("sum > 1").Limit(3).Offset(1).Get()
	fmt.Println(db.LastSql) //打印sql日志

	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res)

	//user, err := db.Query("select id,age from users where id>? limit ?", 1, 3)
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//
	//
	//fmt.Println("user" ,user)
}
