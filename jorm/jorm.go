package jorm

import (
	"fmt"
	"github.com/AntsWorld/jrm-gorm/jdb"
	"log"
	"os"
	"path/filepath"
	"text/template"
)

type JrOrm struct {
	DBConnector *jdb.DBConnector
}

//new JrOrm
func New(dbConnector *jdb.DBConnector) *JrOrm {
	return &JrOrm{DBConnector: dbConnector}
}

//close db connect
func (orm *JrOrm) Close() error {
	if orm.DBConnector != nil {
		return orm.DBConnector.Close()
	}
	return nil
}

//查询数据库信息
func (orm *JrOrm) QuerySchemaInfo(databaseName string, schemaInfo *SchemaInfo) error {
	//查询数据库中所有的表
	sqlStatement := "select * from information_schema.tables where table_schema=? and table_type='base table';"
	var schemaTables []SchemaTable
	if err := orm.DBConnector.ExecuteQuery(sqlStatement, &schemaTables, databaseName); err != nil {
		return err
	}
	schemaInfo.Tables = schemaTables
	return nil
}

//查询数据库表信息
func (orm *JrOrm) QuerySchemaTableInfo(schemaTable SchemaTable, schemaTableInfo *SchemaTableInfo) error {
	tableName := schemaTable.TableName
	sqlStatement := "select * from information_schema.columns where table_name =?;"
	var schemaTableColumns []SchemaTableColumns
	if err := orm.DBConnector.ExecuteQuery(sqlStatement, &schemaTableColumns, tableName); err != nil {
		return err
	}
	schemaTableInfo.TableColumns = schemaTableColumns
	schemaTableInfo.SchemaTable = schemaTable
	return nil
}

//根据模板字符串和模板数据生成文件
func (orm *JrOrm) GenerateTemplateFile(text string, templateData interface{}, filePath string, fileName string) error {
	//判断filePath是否存在
	if file, err := os.Open(filePath); err != nil {
		if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
			log.Println("MkdirAll error:", err)
			return err
		}
	} else {
		if err := file.Close(); err != nil {
			return err
		}
	}

	fileAbsolutePath := fmt.Sprintf("%s%c%s", filePath, filepath.Separator, fileName)
	file, err := os.Create(fileAbsolutePath)
	if err != nil {
		return err
	}
	//解析模板
	tmpl, err := template.New("apiDoc").Parse(text)
	if err != nil {
		return err
	}
	//写到文件
	return tmpl.Execute(file, templateData)
}
