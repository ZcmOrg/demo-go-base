/**
 * 这是 Go 提供的操作 SQL/SQL-Like 数据库的通用接口，但 Go 标准库并没有提供具体数据库的实现，需要结合第三方的驱动来使用该接口
 * 该包有一个子包：driver，它定义了一些接口供数据库驱动实现，一般业务代码中使用 database/sql 包即可，尽量避免使用 driver 这个子包
 */

package database

import (
	"database/sql"
	"fmt"
	"time"
)

func testOpenDB() (*sql.DB, err) {
	db, _ := sql.Open("mysql", "root:@tcp(localhost:3306)/test?charset=utf8") //一个数据库实例
	db.Driver()                                                               //Driver方法返回数据库下层驱动
	db.Ping()                                                                 //Ping检查与数据库的连接是否仍有效，如果需要会创建连接
	db.SetMaxIdleConns()                                                      //连接池最大空闲连接数
	db.SetMaxOpenConns()                                                      //连接池最多连接数
	db.Exec("", "")
	return db, nil
}

func testCloseDB(db *sql.DB) {
	db.Close() //Close关闭数据库，释放任何打开的资源。一般不会关闭 DB，因为 DB 句柄通常被多个 go 程共享，并长期活跃
}

func testPing() {
	db, _ := openDB()
	defer closeDB(db)
	if err := db.Ping(); err != nil {
		t.Errorf("err was %#v, expected nil", err)
		return
	}
}

func testQuery() {
	db, _ := openDB()
	defer closeDB(db)
	rows, err := db.Query("select * from table") //Query执行一次查询，返回多行结果（即Rows），一般用于执行 select 命令。参数 args 表示query中的占位参数
	if err != nil {
		t.Fatalf("Query: %v", err)
	}
	cName, err := rows.Columns() //Columns返回列名
	fmt.Println(cName)
	type row struct {
		age  int
		name string
	}
	got := []row{}
	for rows.Next() { //Next准备用于Scan方法的下一行结果
		var r row
		err = rows.Scan(&r.age, &r.name)
		if err != nil {
			t.Fatalf("Scan: %v", err)
		}
		got = append(got, r)
	}
	rows.Close()     //Close关闭Rows
	err = rows.Err() //Err返回可能的、在迭代时出现的错误
	if err != nil {
		t.Fatalf("Err: %v", err)
	}
}

func testQueryRow() {
	db, _ := openDB()
	defer closeDB(db)
	var name string
	var age int
	var birthday time.Time

	//QueryRow执行一次查询，并期望返回最多一行结果（即Row）
	db.QueryRow("select age from users where uid = ?", 3).Scan(&age) //Scan将该行查询结果各列分别保存进dest参数指定的值中
	db.QueryRow("select age from users where uid = ?", 3).Scan(&birthday)
}

func testStatement() {
	db, _ := openDB()
	defer closeDB(db)
	stmt, err := db.Prepare("SELECT|people|age|name=?") //Prepare创建一个准备好的状态用于之后的查询和命令
	if err != nil {
		t.Fatalf("Prepare: %v", err)
	}
	err = stmt.Close() //Close关闭状态
	if err != nil {
		t.Fatalf("Close: %v", err)
	}
	var name string
	err = stmt.QueryRow("foo").Scan(&name) //QueryRow使用提供的参数执行准备好的查询状态
	if err == nil {
		t.Errorf("expected error from QueryRow.Scan after Stmt.Close")
	}
}

func testExec() {
	db, _ := openDB()
	defer closeDB(db)
	exec(t, db, "CREATE|t1|name=string,age=int32,dead=bool") //Exec执行一次命令（包括查询、删除、更新、插入等），不返回任何执行结果。参数 args 表示 query 中的占位参数
	stmt, err := db.Prepare("INSERT|t1|name=?,age=?")
	if err != nil {
		t.Errorf("Stmt, err = %v, %v", stmt, err)
	}
	defer stmt.Close()
}

func testTx() {
	db, _ := openDB()
	defer closeDB(db)
	tx, _ := db.Begin() //Begin开始一个事务
	name := "test"
	age := 2
	tx.Exec("update users set name = ? where age = ?", name, age) //Exec执行命令，但不返回结果。例如执行insert和update
	rows, _ := tx.Query()                                         //Query执行查询并返回零到多行结果（Rows）
	row := tx.QueryRow()                                          //QueryRow执行查询并期望返回最多一行结果（Row）
	stmt, _ := tx.Prepare()                                       //Prepare准备一个专用于该事务的状态
	tx.Commit()                                                   //Commit递交事务
	tx.Rollback()                                                 //Rollback放弃并回滚事务
}
