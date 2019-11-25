package jdb

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

//######################################################################################################################
//数据库连接器
type DBConnector struct {
	DB         *sql.DB
	OrmTagName string
}

//新建数据库连接器
func NewDbConnector() *DBConnector {
	return &DBConnector{}
}

//打开数据库连接
func (connector *DBConnector) Open(driverName string, dataSourceName string) error {
	if db, err := sql.Open(driverName, dataSourceName); err != nil {
		return err
	} else if err := db.Ping(); err != nil {
		return err
	} else {
		connector.DB = db
	}
	return nil
}

//执行更新语句，包括：插入、删除、修改
func (connector *DBConnector) ExecuteUpdate(sqlStatement string, args ...interface{}) (sql.Result, error) {
	var sqlExecutor = SqlExecutor{DB: connector.DB}
	return sqlExecutor.ExecuteNonSelectSql(sqlStatement, args...)
}

//执行查询语句
func (connector *DBConnector) ExecuteQuery(sqlStatement string, v interface{}, args ...interface{}) error {
	var sqlExecutor = SqlExecutor{DB: connector.DB, StructureOrmTagName: connector.OrmTagName}
	return sqlExecutor.ExecuteSelectSql(sqlStatement, v, args...)
}

//执行查询语句并返回Rows
func (connector *DBConnector) ExecuteQueryForRows(sqlStatement string, args ...interface{}) (*sql.Rows, error) {
	var sqlExecutor = SqlExecutor{DB: connector.DB, StructureOrmTagName: connector.OrmTagName}
	return sqlExecutor.ExecuteSelectSqlForRows(sqlStatement, args...)
}

//执行查询语句并返回[]map[string]string
func (connector *DBConnector) ExecuteQueryForKeyValueMap(sqlStatement string, args ...interface{}) ([]map[string]string, error) {
	var sqlExecutor = SqlExecutor{DB: connector.DB, StructureOrmTagName: connector.OrmTagName}
	return sqlExecutor.ExecuteSelectSqlForMapResult(sqlStatement, args...)
}

//判断数据库是否已连接
func (connector *DBConnector) IsConnected() (bool, error) {
	if nil == connector.DB {
		return false, nil
	}
	if err := connector.DB.Ping(); err != nil {
		return false, err
	}
	return true, nil
}

//关闭数据库连接
func (connector *DBConnector) Close() error {
	if nil == connector.DB {
		return nil
	}
	return connector.DB.Close()
}
