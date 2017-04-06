package sql

import (
	"database/sql"
	"errors"
	"fmt"
)

import _ "github.com/go-sql-driver/mysql"

const (
	SelectSql = iota
	InsertSql
	UpdateSql
	DeleteSql
)

type DbInfo struct {
	DbType     string
	DbHost     string
	DbName     string
	DbUser     string
	DbPassWd   string
	DbCharset  string
	DbMaxConns int
	DbMinConns int
}

type SqlAlchemy struct {
	sqlStr    string
	values    []interface{}
	first     bool
	wFirst    bool
	varNum    int
	tableName string
	sqlType   int
	dbPool    *sql.DB
}

func (this *SqlAlchemy) Init(dbInfo *DbInfo) (int, error) {
	dbConnStr := dbInfo.DbUser + ":" + dbInfo.DbPassWd + "@tcp(" + dbInfo.DbHost + ")/" + dbInfo.DbName + "?charset=" + dbInfo.DbCharset
	fmt.Println(dbConnStr)
	var err error
	this.dbPool, err = sql.Open(dbInfo.DbType, dbConnStr)
	if err != nil {
		return -1, err
	}
	this.dbPool.SetMaxOpenConns(dbInfo.DbMaxConns)
	this.dbPool.SetMaxIdleConns(dbInfo.DbMinConns)
	err = this.dbPool.Ping()
	if err != nil {
		return -1, err
	}
	return 0, nil
}

func (this *SqlAlchemy) Select(tableName string) *SqlAlchemy {
	this.sqlStr = "select "
	this.tableName = tableName
	this.setProperties(SelectSql)
	return this
}

func (this *SqlAlchemy) Insert(tableName string) *SqlAlchemy {
	this.sqlStr = "insert into " + tableName + "("
	this.setProperties(InsertSql)
	return this
}

func (this *SqlAlchemy) Update(tableName string) *SqlAlchemy {
	this.sqlStr = "update " + tableName + " set "
	this.setProperties(UpdateSql)
	return this
}

func (this *SqlAlchemy) Delete(tableName string) *SqlAlchemy {
	this.sqlStr = "delete from " + tableName
	this.setProperties(DeleteSql)
	return this
}

func (this *SqlAlchemy) ExecQuery() (*sql.Rows, error) {
	if this.sqlType != SelectSql {
		return nil, errors.New("sqlAlchemy:Is not select sql")
	}
	if this.first == true {
		this.sqlStr += " from " + this.tableName
	}
	fmt.Println(this.sqlStr)
	fmt.Println(this.values)
	fmt.Println(this.varNum)
	return nil, nil
	//return this.dbPool.Query(this.sqlStr, this.values)
}

func (this *SqlAlchemy) Execute() (int64, error) {
	if this.sqlType == InsertSql {
		this.sqlStr += ") values("
		for i := 0; i < this.varNum; i++ {
			if i != 0 {
				this.sqlStr += ","
			}
			this.sqlStr += "?"
		}
		this.sqlStr += ")"
		fmt.Println(this.sqlStr)
		fmt.Println(this.values)
		fmt.Println(this.varNum)
		return 0, nil
		/*res, err := this.stmtExecute()
		if err != nil {
			return -1, err
		}
		return res.LastInsertId(), nil*/
	} else if this.sqlType == UpdateSql || this.sqlType == DeleteSql {
		fmt.Println(this.sqlStr)
		fmt.Println(this.values)
		fmt.Println(this.varNum)
		return 0, nil
		/*res, err := this.stmtExecute()
		if err != nil {
			return -1, err
		}
		return res.RowsAffected(), nil*/
	}
	return -1, errors.New("SqlAlchemy:Is not valid sql")
}

func (this *SqlAlchemy) S(field string, value interface{}) *SqlAlchemy {
	if this.first != true {
		this.sqlStr += ","
	}
	this.first = false
	this.sqlStr += field + "=?"
	this.values = append(this.values, value)
	this.varNum++
	return this
}

func (this *SqlAlchemy) V(field string, value interface{}) *SqlAlchemy {
	if this.first != true {
		this.sqlStr += ","
	}
	this.first = false
	this.sqlStr += field
	this.values = append(this.values, value)
	this.varNum++
	return this
}

func (this *SqlAlchemy) F(fields ...string) *SqlAlchemy {
	for i, name := range fields {
		if i != 0 {
			this.sqlStr += ","
		}
		this.sqlStr += name
	}

	return this
}

func (this *SqlAlchemy) W(fieldName string, value interface{}, sign string) *SqlAlchemy {
	this.wFirst = true
	return this.condition(fieldName, value, sign, "")
}

func (this *SqlAlchemy) And(fieldName string, value interface{}, sign string) *SqlAlchemy {
	return this.condition(fieldName, value, sign, "and")
}

func (this *SqlAlchemy) Or(fieldName string, value interface{}, sign string) *SqlAlchemy {
	return this.condition(fieldName, value, sign, "or")
}

func (this *SqlAlchemy) condition(fieldName string, value interface{}, sign string, linkStr string) *SqlAlchemy {
	if this.wFirst == true {
		if this.sqlType == SelectSql {
			this.sqlStr += " from " + this.tableName + " where "
		} else if this.sqlType != InsertSql {
			this.sqlStr += " where "
		}
		this.wFirst = false
	} else {
		this.sqlStr += " " + linkStr + " "
	}
	this.sqlStr += fieldName + sign + "?"
	this.values = append(this.values, value)
	this.varNum++

	return this
}

func (this *SqlAlchemy) stmtExecute() (sql.Result, error) {
	stmt, err := this.dbPool.Prepare(this.sqlStr)
	if err != nil {
		return nil, err
	}
	return stmt.Exec(this.values)
}

func (this *SqlAlchemy) setProperties(sqlType int) {
	this.sqlType = sqlType
	this.values = nil
	this.first = true
	this.wFirst = true
	this.varNum = 0
}
