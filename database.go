package gosql

import (
	"fmt"
	"github.com/CourrageWang/gosql/utils"
	"errors"
	"strings"
	"database/sql"
	"strconv"
)

/**
  sql 核心处理类 （获取用户的输入参数，并将其解析为可执行的 sql ）
 */
var (
	judge = [] string{"=", "<", ">", "!="} //  查询参数合法性验证
	tx    *sql.Tx                          // 事务支持
)

// 数据库的结构映射
type Database struct {
	distinct     bool //  distinct
	table        string
	fileds       string           //  字段
	join         [][] interface{} // join
	avg          string           //  AVG()映射
	sum          string           // SUM ()映射
	max          string           //MAX()映射
	min          string           //MIN()映射
	limit        int              // Limit 映射
	offset       int              // offset 映射
	group        string           // group by映射
	order        string           // order by映射
	having       string           // having映射
	whers        [][] interface{}
	SqlLogs      [] string   // 所有sql日志
	LastSql      string      //  最新sql
	LastInsertId int         //最后插入数据的id
	RowsAffected int         // 插入受影响的行数
	data         interface{} //  更新或插入的数据
	trans        bool        //是否开启事务
}

//  查询的字段
func (dba *Database) Fileds(fileds string) *Database {
	dba.fileds = fileds
	return dba
}

// 插入或更新的数据
func (dba *Database) Data(data interface{}) *Database {
	dba.data = data
	return dba
}

//  查询的表
func (dba *Database) Table(table string) *Database {
	dba.table = table
	return dba
}

//  Where 查询条件 (and Where)
func (dba *Database) Where(args ...interface{}) *Database {
	wh := []interface{}{"and", args}  //  exp: [[and] [id > 2]] // 数组第一个保存条件(and/or) 第二个为参数[]interface{}
	dba.whers = append(dba.whers, wh) // exp: [[and] [id >2] [and] [id < 5]]
	return dba
}

//  or 关键字
func (dba *Database) OrWhere(args ...interface{}) *Database {
	wh := []  interface{}{"or", args}
	dba.whers = append(dba.whers, wh)
	return dba
}

// limit 关键字
func (dba *Database) Limit(limit int) *Database {
	dba.limit = limit
	return dba
}

// offset 关键字
func (dba *Database) Offset(offset int) *Database {
	dba.offset = offset
	return dba
}

//  group by 关键字
func (dba *Database) Group(group string) *Database {
	dba.group = group
	return dba
}

// order by 关键字
func (dba *Database) Order(order string) *Database {
	dba.order = order
	return dba
}

//  left join 关键字
func (dba *Database) Leftjon(args ...interface{}) *Database {
	dba.join = append(dba.join, []interface{}{"LEFT", args})
	return dba
}

// right join 关键字
func (dba *Database) Rightjoin(args ...interface{}) *Database {
	dba.join = append(dba.join, [] interface{}{"RIGHT", args})
	return dba
}

// having 关键字映射
func (dba *Database) Having(having string) *Database {
	dba.having = having
	return dba
}

//  Avg 函数
func (dba *Database) Avg(avg string) (interface{}, error) {
	return dba.buildUnion("avg", avg)
}

// sum 函数
func (dba *Database) Sum(sum string) (interface{}, error) {
	return dba.buildUnion("sum", sum)
}

// max函数
func (dba *Database) Max(max string) (interface{}, error) {
	return dba.buildUnion("max", max)
}

// min 函数
func (dba *Database) Min(min string) (interface{}, error) {
	return dba.buildUnion("min", min)
}

//https://www.jianshu.com/p/ee0d2e7bef54
// 建立联合查询
func (dba *Database) buildUnion(untion, fieds string) (interface{}, error) {
	// exp  :avg(age) as avg
	ustr := untion + "(" + fieds + ") as " + untion
	switch untion {
	case "avg":
		dba.avg = ustr
	case "sum":
		dba.sum = ustr
	case "max":
		dba.max = ustr
	case "min":
		dba.min = ustr
	}
	//  构建sql
	sqls, err := dba.buildQuery()
	if err != nil {
		return nil, err
	}
	//查询
	result, err := dba.Query(sqls)
	if err != nil {
		return nil, err
	}
	return result[0][untion], nil

}

// 获取所有的查询参数并将其转换为sqlStr
func (dba *Database) buildQuery() (string, error) {
	// Union
	union := [] string{
		dba.avg,
		dba.sum,
		dba.max,
		dba.min,
	}
	var uion string
	for _, item := range union {
		if item != "" {
			uion = item
			break
		}
	}

	//  distinct
	distinct := utils.Judge(dba.distinct, "distinct", "")
	//  table
	table := Connect.UseConfig["prefix"] + dba.table
	//  字段
	fileds := utils.Judge(dba.fileds == "", "*", dba.fileds).(string)
	// join
	pjoin, err := dba.ParseJoin()
	if err != nil {
		return "", err
	}
	join := pjoin
	parWhere, err := dba.pasreWhere()
	if err != nil {
		return "", err
	}
	//where 条件
	where := utils.Judge(parWhere == "", "", "WHERE"+parWhere).(string)
	// having 条件
	group := utils.Judge(dba.group == "", "", "GROUP BY "+dba.group).(string)
	// order by 条件
	order := utils.Judge(dba.order == "", "", "ORDER BY "+dba.order).(string)
	// limit
	limit := utils.Judge(dba.limit == 0, "", "LIMIT "+strconv.Itoa(dba.limit))
	// offset
	offset := utils.Judge(dba.offset == 0, "", "OFFSET "+strconv.Itoa(dba.offset))
	// having
	having := utils.Judge(dba.having == "", "", "HAVING "+dba.having).(string)

	sqlstr := fmt.Sprintf("SELECT %s%s FROM %s %s %s %s %s %s %s %s", distinct, utils.Judge(uion != "", uion, fileds), table, join, where, group, having, order, limit, offset)
	return sqlstr, nil //  生成最终的SQL字符串
}

// 解析where
func (dba *Database) pasreWhere() (string, error) {
	wheres := dba.whers // 获取所有的where
	var final [] string //  存放最终的执行语句
	for _, args := range wheres { //循环解析 获取条件以及参数
		var cond string = args[0].(string) // and/ or 条件
		params := args[1].([]interface{})  // 查询参数

		parmsLen := len(params)
		switch parmsLen {
		case 3: // exp {"id" ,">" "2"}
			res, err := dba.Parse(params)
			if err != nil {
				return res, err
			}
			final = append(final, cond+" "+res)
		case 1: // 参数为一维数组
			switch parmReal := params[0].(type) {
			case map[string]interface{}: // 一维数组

				var whereArr [] string
				for k, v := range parmReal {
					whereArr = append(whereArr, k+"="+utils.AddSingleMark(v))
				}
				final = append(final, cond+" ("+strings.Join(whereArr, " and ")+")")

			}

		}
	}
	return strings.TrimLeft(strings.Trim(strings.Join(final, " "), " "), "and"), nil
}

// 解析join
func (dba *Database) ParseJoin() (string, error) {
	var join [] interface{}
	var final [] string
	joinArr := dba.join
	for _, join = range joinArr {
		var w string
		var ok bool
		var args [] interface{}
		if len(join) != 2 { // 查询条件有两个 [0] 表示关键字，[1]表示查询的条件
			return "", errors.New("左查询条件有误")
		}
		//  获取条件
		if args, ok = join[1].([]interface{}); !ok {
			return "", errors.New("左查询条件有误")
		}
		argLen := len(args)
		switch argLen {
		case 4:
			w = args[0].(string) + " ON " + args[1].(string) + " " + args[2].(string) + " " + args[3].(string)
		default:
			return "", errors.New("格式化有误")

		}
		final = append(final, " "+join[0].(string)+" JOIN "+w)
	}
	return strings.Join(final, " "), nil
}

// 将条件转换为字符串
func (dba *Database) Parse(args []interface{}) (string, error) {
	argsLen := len(args)

	var parmsStore [] string //存储当前数据
	switch argsLen {
	case 3: //  常规参数 // [id > 2]
		if !utils.Has(args[1], judge) { //判断输入的条件是否合法
			return "", errors.New("输入参数不合法")
		}
	}

	parmsStore = append(parmsStore, args[0].(string))
	parmsStore = append(parmsStore, args[1].(string))
	// 最终 id >
	parmsStore = append(parmsStore, utils.AddSingleMark(args[2]))
	return strings.Join(parmsStore, " "), nil // 参数后面
}

// sql 查询实例调用系统实例
/**   sql 分为两种类型 (带有占位符以及不带占位符的)
   db.query ("select id from users where id = 2")
       or
   db.query("select id from users where id = ?" ,2)

 */
func (dba *Database) Query(args ...interface{}) ([]map[string]interface{}, error) {

	Data := make([]map[string]interface{}, 0)
	lenArgs := len(args)
	var vals [] interface{}    //  执行参数
	sqlstr := args[0].(string) // 查询参数
	if lenArgs > 1 { // 带有占位符的sql
		for k, v := range args {
			if k > 0 {
				vals = append(vals, v) // 获取占位符后的参数
			}
		}
	}
	//  待添加sql日志功能
	dba.LastSql = fmt.Sprintf(utils.GetNowTime()+" exectued: "+sqlstr, vals...)
	dba.SqlLogs = append(dba.SqlLogs, dba.LastSql)
	Connect.SqlLog = dba.SqlLogs
	//
	stmt, err := DB.Prepare(sqlstr)
	if err != nil {
		return Data, err
	}
	defer stmt.Close()
	rows, e := stmt.Query(vals ...) // 调用原生SQl执行查询
	if e != nil {
		return Data, err
	}
	defer rows.Close()
	columns, err := rows.Columns() // 返回列明
	if err != nil {
		return Data, err
	}
	count := len(columns)
	values := make([] interface{}, count)   // 保存返回结果
	ScanParm := make([] interface{}, count) // 查找的参数
	for rows.Next() {
		for i := 0; i < count; i++ {
			ScanParm[i] = &values[i]
		}
		rows.Scan(ScanParm...)              // exp ：rows.Scan (&id , &name)
		res := make(map[string]interface{}) //保存查询的结果
		for i, col := range columns {
			var v interface{}
			vals := values[i] //获取计算的结果
			if b, ok := vals.([] byte); ok {
				v = string(b)
			} else {
				v = vals
			}
			res[col] = v
		}
		Data = append(Data, res)
	}
	return Data, nil
}

// 插入数据
func (dba *Database) Insert() (int, error) {
	sqlstr, err := dba.buildExecute("insert")
	if err != nil {
		return 0, err
	}
	res, err := dba.Execute(sqlstr)
	if err != nil {
		return 0, err
	}
	return int(res), nil
}
func (dba *Database) Delete() (int, error) {
	sqlstr, err := dba.buildExecute("delete")
	if err != nil {
		return 0, err
	}
	res, errs := dba.Execute(sqlstr)
	if errs != nil {
		return 0, nil
	}
	return int(res), nil

}

//  更新数据
func (dba *Database) Update() (int, error) {
	// 生成sql
	sqlstr, err := dba.buildExecute("update")

	if err != nil {
		return 0, err
	}
	res, errs := dba.Execute(sqlstr)
	if errs != nil {
		return 0, err
	}
	return int(res), nil
}

/**
   添加日志功能
 */
//func (dba *Database) GetLastSql() string  {
//	if len(Connect.SqlLog)>0 {
//		//return Connect.SqlLog[]
//	}
//
//}
// 获取结果集 {[]map[string] interface{}}
func (dba *Database) Get() ([]map[string]interface{}, error) {
	// 构建sql
	sqls, err := dba.buildQuery()
	if err != nil {
		return nil, err
	}
	// 执行sql查询
	result, err := dba.Query(sqls)
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return nil, nil
	}
	return result, nil
}

// 执行SQl语句 支持sql语句 ，以及带有占位符的sql语句 调用 (增删改查需要考虑事务机制)
/**   一般的查询使用的db对象， 事务 ： sql.Tx,（其会从连接池中取一个空闲连接）直到回退或者
提交之后才会把连接释放到连接池。
 */

func (dba *Database) Execute(args ...interface{}) (int64, error) {
	lenargs := len(args)
	var sqlstr string
	var val [] interface{}
	sqlstr = args[0].(string) // 获取sql
	// sql 存在两种情况，exp: ("update users set id = ? where age = ? ", 100,2 ) or ("update user set age =2  where id =3")
	if lenargs > 1 {
		for k, v := range args {
			if k > 0 {
				val = append(val, v) // 追加所有的条件
			}
		}
	}
	// 添加日志
	dba.LastSql = fmt.Sprintf(utils.GetNowTime()+" exectued: "+sqlstr, val...)
	dba.SqlLogs = append(dba.SqlLogs, dba.LastSql)
	//  取出关键字
	var SqlType string = strings.ToLower(sqlstr[0:6]) // 取出exp : delete、update、insert、
	if SqlType == "SELECT" {
		return 0, errors.New("非法使用SELECT 关键字")
	}
	// 是否开启事务机制。
	if dba.trans == true { //  开启
		stmt, err := tx.Prepare(sqlstr)

		if err != nil {
			return 0, err
		}
		return dba.pasreExecute(stmt, SqlType, val);
	}
	//  未开启
	stmt, err := DB.Prepare(sqlstr)
	if err != nil {
		return 0, nil
	}
	return dba.pasreExecute(stmt, SqlType, val)
}

// 解析执行的条件 并调用 sql.Stmt 的Exec（）执行sql语句
func (dba *Database) pasreExecute(stmt *sql.Stmt, sqlType string, vals [] interface{}) (int64, error) {
	var rowAffected int64 //受影响的行数
	var err error
	result, errs := stmt.Exec(vals...)

	if errs != nil {
		return 0, errs
	}
	switch sqlType {
	case "insert":
		lastInseted, err := result.LastInsertId()
		if err == nil {
			dba.LastInsertId = int(lastInseted)
		}
		rowAffected, err = result.RowsAffected()
		dba.RowsAffected = int(rowAffected)

	case "update":
		rowAffected, err = result.RowsAffected()
	case "delete":
		rowAffected, err = result.RowsAffected()
	}
	return rowAffected, err
}

// 解析语句并生成可执行的sql
func (dba *Database) buildExecute(sqlType string) (interface{}, error) {
	//  exp : update
	var updatstr, insertkey, insertv, sqlstr string
	if sqlType != "delete" {
		updatstr, insertkey, insertv = dba.buiildData()
	}

	res, err := dba.pasreWhere()

	if err != nil {
		return res, err
	}

	where := utils.Judge(res == "", "", " WHERE "+res).(string)
	table := Connect.UseConfig["prefix"] + dba.table
	switch sqlType {
	case "update":
		sqlstr = fmt.Sprintf("update %s set %s%s", table, updatstr, where)
	case "insert":
		sqlstr = fmt.Sprintf("insert into %s(%s) values %s", table, insertkey, insertv)
	case "delete":
		sqlstr = fmt.Sprintf("delete from  %s %s", table, where)
	}
	return sqlstr, nil
}

// 生成更新的数据
func (dba *Database) buiildData() (string, string, string) {
	var datafileds [] string  // 插入的字段
	var datavalues []  string //  插入的值

	var dataObj [] string // 存储更新or删除的数据
	data := dba.data
	switch data.(type) {
	case map[string]interface{}:
		datas := make(map[string]string) // 存储信息
		// step --- 将所有数据中的值转换为string
		switch data.(type) {
		case map[string]interface{}:
			for key, V := range data.(map[string]interface{}) {
				datas[key] = utils.TransStr(V) //解析参数并将其值转换为string存储在map中
			}
		case map[string]int:
			for key, v := range data.(map[string]int) {
				datas[key] = utils.TransStr(v)
			}
		case map[string]string:
			for key, v := range data.(map[string]string) {
				datas[key] = v
			}
		}
		var tempStr [] string
		for key, value := range datas {
			// insert
			datafileds = append(datafileds, key)
			tempStr = append(tempStr, utils.AddSingleMark(value))
			// update  exp :[age='15' name='王勇琪15']
			dataObj = append(dataObj, key+"="+utils.AddSingleMark(value))
		}
		// 组合生成insert
		datavalues = append(datavalues, "("+strings.Join(tempStr, ",")+")")
	}
	return strings.Join(dataObj, ","), strings.Join(datafileds, ","), strings.Join(datavalues, "")
}

// 事务管理
// 开启事务(调用系统DB.Begin)
func (dba *Database) Begin() {
	tx, _ = DB.Begin()
	dba.trans = true
}

// 事务的提交
func (dba *Database) Commit() {
	tx.Commit()
	dba.trans = false // 提交后事务的生命周期结束
}
func (dba *Database) Rollback() {
	tx.Rollback()
	dba.trans = false // 提交后事务的声明周期结束
}

//  事务 (闭包执行)
func (dba *Database) Transaction(closure func() (error)) bool {
	dba.Begin() // 开启事务
	err := closure()
	if err != nil {
		dba.Rollback()
		return false
	}
	dba.Commit()
	return true
}
