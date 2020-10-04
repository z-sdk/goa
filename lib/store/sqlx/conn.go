package sqlx

import (
	"database/sql"
	"errors"
	"time"
)

const (
	// tagName 结构体字段中，数据库字段的标记名称
	tagName = "db"

	// 数据库慢日志阈值，用于记录慢查询和慢执行
	slowThreshold = 500 * time.Millisecond
)

var (
	ErrNotFound             = errors.New("没有结果集")
	ErrNotSettable          = errors.New("扫描目标不可设置")
	ErrUnsupportedValueType = errors.New("不支持的扫描目标类型")
	ErrNotReadableValue     = errors.New("无法读取的值，检查结构字段是否大写开头")
)

type (
	// StmtConn 语句执行和查询接口
	//StmtConn interface {
	//	Query(dest interface{}, args ...interface{}) error
	//	Exec(args ...interface{}) (sql.Result, error)
	//	Close() error
	//}

	// stmtConn 预编译连接
	//stmtConn interface {
	//	Query(args ...interface{}) (*sql.Rows, error)
	//	Exec(args ...interface{}) (sql.Result, error)
	//}

	// statement 预编译语句会话：将查询封装为预编译语句，供底层查询和执行
	//statement struct {
	//	stmt *sql.Stmt
	//}

	// Session 提供外部查询和执行的会话接口
	Session interface {
		Query(dest interface{}, query string, args ...interface{}) error
		Exec(query string, args ...interface{}) (sql.Result, error)
		//Prepare(query string) (StmtConn, error)
	}

	// session 提供内部查询和执行的会话接口
	session interface {
		Query(query string, args ...interface{}) (*sql.Rows, error)
		Exec(query string, args ...interface{}) (sql.Result, error)
	}

	// TransactFn 事务内执行函数，传入事务会话
	TransactFn func(session Session) error

	// Conn 提供外部数据库会话和事务的接口
	Conn interface {
		Session
		Transact(fn TransactFn) error
	}

	// conn 包内连接实例，封装查询、执行、事务及断路器支持
	conn struct {
		driverName     string    // 驱动名称，支持 mysql/postgres/clickhouse 等 sql-like
		dataSourceName string    // 数据源名称 Data Source Name，既数据库连接字符串
		beginTx        beginTxFn // 可开始事务
	}

	// Option 是一个可选的数据库增强函数
	Option func(c *conn)
)

// NewConn 新建指定数据库驱动和DSN的连接
func NewConn(driverName, dataSourceName string, opts ...Option) Conn {
	prefectDSN(&dataSourceName)

	c := &conn{
		driverName:     driverName,
		dataSourceName: dataSourceName,
		beginTx:        beginTx,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// ----------------- conn 实现方法 ↓ ----------------- //

// Query 执行数据库查询并将结果扫描至结果。
// 如果 dest 字段不写tag的话，系统按顺序配对，此时需要与sql中的查询字段顺序一致
// 如果 dest 字段写了tag的话，系统按名称配对，此时可以和sql中的查询字段顺序不同
func (c *conn) Query(dest interface{}, query string, args ...interface{}) error {
	db, err := getConn(c.driverName, c.dataSourceName)
	if err != nil {
		logConnError(c.dataSourceName, err)
		return err
	}
	return doQuery(db, func(rows *sql.Rows) error {
		return scan(rows, dest)
	}, query, args...)
}

func (c *conn) Exec(query string, args ...interface{}) (sql.Result, error) {
	db, err := getConn(c.driverName, c.dataSourceName)
	if err != nil {
		logConnError(c.dataSourceName, err)
		return nil, err
	}
	return doExec(db, query, args...)
}

func (c *conn) Transact(fn TransactFn) error {
	return doTx(c, c.beginTx, fn)
}

// Prepare 创建一个稍后查询或执行的预编译语句
//func (c *conn) Prepare(query string) (stmt StmtConn, err error) {
//	db, err := getConn(c.driverName, c.dataSourceName)
//	if err != nil {
//		logConnError(c.dataSourceName, err)
//		return nil, err
//	}
//	if st, err := db.Prepare(query); err != nil {
//		return nil, err
//	} else {
//		stmt = statement{stmt: st}
//		return stmt, nil
//	}
//}

//func (s statement) Query(dest interface{}, args ...interface{}) error {
//	panic("implement me")
//}
//
//func (s statement) Exec(args ...interface{}) (sql.Result, error) {
//	panic("implement me")
//}
//
//func (s statement) Close() error {
//	panic("implement me")
//}