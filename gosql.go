package gosql

import (
	"database/sql"
	"errors"
	"github.com/CourrageWang/gosql/drivers"
)

/*    数据库的初始化、
 */
var (
	// 原始的sql DB
	DB      *sql.DB
	Connect Connection
)

func init() {
	Connect.MaxCons = 0
	Connect.IdleCons = -1
}

// 实例化数据库
/**
    考虑到后期的扩展，多个数据库的问题
 */
func Open(args ...interface{}) (Connection, error) {
	if len(args) == 1 { // 只有一个参数使用配置中的默认数据库
	} else if len(args) == 2 { // 传入两个参数，使用第二个参数作为驱动的数据库
		if conf, ok := args[1].(string); ok {
			Connect.Default = conf //  得到数据库驱动
		} else { //   传入错误的数据库格式  exp： open（ config ，123）
			return Connect, errors.New("数据库的格式只能为字符串")
		}
	} else {
		return Connect, errors.New("传入参数有误...")
	}

	// 解析配置文件
	error := Connect.parseConfig(args[0])
	if error != nil {
		return Connect, error
	}
	// 加载驱动
	er := Connect.load()
	return Connect, er
}

// 数据库连接信息
type Connection struct {
	DbConfig  map[string]interface{} // 所有数据库配置信息
	UseConfig map[string]string      //  当前需要的配置
	MaxCons   int                    //  数据库连接池最大连接数 默认不受限制
	IdleCons  int                    // 闲置连接数
	Default   string                 //默认的数据库
	SqlLog    [] string              //所有的SQL信息
}

//  加载数据库驱动
func (conn *Connection) load() error {
	dbo := Connect.UseConfig
	var driver, dco string
	var err error
	switch dbo["driver"] {
	case "mysql":
		driver, dco = drivers.Mysql(dbo)
	}
	DB, err = sql.Open(driver, dco)
	//  驱动数据库
	DB.SetMaxOpenConns(conn.MaxCons)
	DB.SetMaxIdleConns(conn.IdleCons)
	if err != nil {
		return err
	}
	er := DB.Ping() //验证数据库连接是否存活
	return er
}

// 解析数据配置 (支持简单配置以及复杂配置)
func (conn *Connection) parseConfig(args interface{}) error {
	// 但数据库配置
	if conf, ok := args.(map[string]string); ok {
		conn.UseConfig = conf
	} else {
		return errors.New("format error in database config!")
	}
	return nil
}

//验证数据库连接是否存活
func (conn *Connection) Ping() error {
	return DB.Ping()
}
func (conn *Connection) Close() error {
	return DB.Close()
}
func (conn *Connection) GetInstance() *Database {
	return &Database{}
}