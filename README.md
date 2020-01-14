# gosql
用golang实现的自己的一个数据库框架
使用简单，只需简单配置即可
支持两种调用方式的调用
1、直接书写sql语句通过
db.Execute("insert into users (id ,age) values (8,9)")
Execute 方式支持 ： 具有占位符的sql 语句
exp：
res, err := db.Execute("update users set age = ? where id = ? ", 100,4)
res, err := db.Execute("update users set age = 100 where id = 4 ")

2、方法调用
res, err := db.Table("users").Data(data).Insert()
	fmt.Println(res)
	if err != nil {
		fmt.Println(err)
		return
	}
3 、提供sql日志管理 
fmt.Println(db.LastSql) //打印sql日志
4、支持简单的事务管理机制