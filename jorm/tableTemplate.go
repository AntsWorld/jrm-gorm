package jorm

//######################################################################################################################
//数据库表模板
//######################################################################################################################
//基础表模板数据
type BaseTableTemplate struct {
	DatabaseName       string               //数据库名称
	TableName          string               //表名称
	PriColumns         []SchemaTableColumns //表中所有PRI字段
	UniqueColumns      []SchemaTableColumns //表中所有Unique字段
	EditableColumns    []SchemaTableColumns //表中所有可编辑字段,插入和更新SQL语句使用
	PackageName        string               //包名,默认用数据库名称转驼峰式命名
	GoFileName         string               //生成的Go文件名,默认使用表名称的驼峰格式并且首字母小写
	GoTestFileName     string               //生成的Go测试文件名
	TableStructureName string               //表结构体名称,默认使用表名称的驼峰格式
}

//默认表模板参数
type DefaultTableTemplate struct {
	BaseTableTemplate BaseTableTemplate //基础模板数据
	TableStructure    string            //表结构体
	InsertTemplates   []InsertTemplate  //Insert模板数据
	DeleteTemplates   []DeleteTemplate  //Delete模板数据
	QueryTemplates    []QueryTemplate   //Query模板数据
	UpdateTemplates   []UpdateTemplate  //Update模板数据
}

//插入函数模板参数
//函数名、函数参数列表(参数名、参数类型)、返回值列表(返回值类型)
//SQL语句中Columns列表、Columns值列表
type InsertTemplate struct {
	FunctionName               string   //函数名,生成测试代码时需要使用
	InsertRequestStructure     string   //插入函数请求结构体
	InsertRequestStructureName string   //插入函数请求结构体名称
	InsertSqlStatement         string   //SQL语句
	InsertSqlTableColumns      []string //SQL语句中表字段名
	InsertSqlStructureColumns  []string //SQL语句中结构体字段名
}

//更新模板参数，需要满足以下条件：一次更新所有可更新字段、单独更新每个字段
//表结构体名称、表可更新字段结构体、表可更新字段结构体名称、数据库名、表名、需要更新字段列表、需要更新字段列表值、更新条件Key、更新条件Value
type UpdateTemplate struct {
	RequestStructure          string   //更新请求结构体
	RequestStructureName      string   //更新请求结构体名称，数据类型
	FunctionName              string   //函数名,生成测试代码时需要使用
	UpdateSqlStatement        string   //SQL语句
	UpdateSqlTableColumns     []string //SQL语句中表字段名
	UpdateSqlStructureColumns []string //SQL语句中结构体字段名
	WhereKeyName              string
	WhereKeyDataType          string
	WhereKeyValueParamName    string
}

//######################################################################################################################
//数据库操作基础模板
var DefaultTableTemplateText = `package {{.BaseTableTemplate.PackageName}}

import (
	"database/sql"
	"github.com/AntsWorld/jrm-gorm/jdb"
)

//数据库名
var DBNameOf{{.BaseTableTemplate.TableStructureName}} = "{{.BaseTableTemplate.DatabaseName}}"

//表名
var TableNameOf{{.BaseTableTemplate.TableStructureName}} = "{{.BaseTableTemplate.TableName}}"

//Table Structure
{{.TableStructure}}
//Insert操作
{{range $k,$v := .InsertTemplates}}
	{{if ne $v.InsertRequestStructure ""}}
		{{if eq $k 0}}{{$v.InsertRequestStructure}}{{end}}
		func {{$v.FunctionName}}(dbConnector *jdb.DBConnector,request *{{$v.InsertRequestStructureName}}) (sql.Result, error) {
			sqlStatement := "{{$v.InsertSqlStatement}}"
			return dbConnector.ExecuteUpdate(sqlStatement,{{range $k,$v := $v.InsertSqlStructureColumns}}{{if  eq $k 0}}request.{{$v}}{{else}}, request.{{$v}}{{end}}{{end}})
		}
	{{end}}
{{end}}
//Delete操作,根据主键和唯一键可以做精确删除；根据普通键可以做批量删除。
{{range $k,$v := .DeleteTemplates}}
	func {{$v.FunctionName}}(dbConnector *jdb.DBConnector,{{$v.WhereKeyValueParamName}} {{$v.WhereKeyDataType}}) (sql.Result, error) {
			sqlStatement := "{{$v.SqlStatement}}"
			return dbConnector.ExecuteUpdate(sqlStatement, {{$v.WhereKeyValueParamName}})
	}
{{end}}
//Query操作，根据主键和唯一键可以做精确查询；根据普通键可以做批量查询；查询需要默认启用分页；需要支持模糊查询；
{{range $k,$v := .QueryTemplates}}
	func {{$v.FunctionName}}(dbConnector *jdb.DBConnector,{{$v.WhereKeyValueParamName}} {{$v.WhereKeyDataType}}) ([]{{$v.QueryResultDataType}}, error) {
		queryResults := make([]{{$v.QueryResultDataType}}, 0)
		sqlStatement := "{{$v.SqlStatement}}"
		if err := dbConnector.ExecuteQuery(sqlStatement,&queryResults,{{$v.WhereKeyValueParamName}});err!=nil{
			return nil,err
		}
		return queryResults,nil
	}
{{end}}
//Update操作
{{range $k,$v := .UpdateTemplates}}
	{{if ne $v.RequestStructure ""}}
		{{if eq $k 0}}{{$v.RequestStructure}}{{end}}
		func {{$v.FunctionName}}(dbConnector *jdb.DBConnector,{{$v.WhereKeyValueParamName}} {{$v.WhereKeyDataType}},request *{{$v.RequestStructureName}}) (sql.Result, error) {
			sqlStatement := "{{$v.UpdateSqlStatement}}"
			return dbConnector.ExecuteUpdate(sqlStatement,{{range $k,$v := $v.UpdateSqlStructureColumns}}{{if  eq $k 0}}request.{{$v}}{{else}}, request.{{$v}}{{end}}{{end}},{{$v.WhereKeyValueParamName}})
		}
	{{else}}
		func {{$v.FunctionName}}(dbConnector *jdb.DBConnector,newValue interface{},{{$v.WhereKeyValueParamName}} {{$v.WhereKeyDataType}}) (sql.Result, error) {
			sqlStatement := "{{$v.UpdateSqlStatement}}"
			return dbConnector.ExecuteUpdate(sqlStatement, newValue, {{$v.WhereKeyValueParamName}})
		}
	{{end}}
{{end}}
`

////######################################################################################################################
//Query操作，根据主键和唯一键可以做精确查询；根据普通键可以做批量查询；查询需要默认启用分页；需要支持模糊查询；
//{{range $k,$v := .QueryTemplates}}
//{{if eq $v.UniQuery true}}
//func {{$v.FunctionName}}(dbConnector *jdb.DBConnector,{{$v.WhereKeyValueParamName}} {{$v.WhereKeyDataType}}) (*{{$v.QueryResultDataType}}, error) {
//queryResult := {{$v.QueryResultDataType}}{}
//sqlStatement := "{{$v.SqlStatement}}"
//if err := dbConnector.ExecuteQuery(sqlStatement,&queryResult,{{$v.WhereKeyValueParamName}});err!=nil{
//return nil,err
//}
//return &queryResult,nil
//}
//{{else}}
//func {{$v.FunctionName}}(dbConnector *jdb.DBConnector,{{$v.WhereKeyValueParamName}} {{$v.WhereKeyDataType}}) ([]{{$v.QueryResultDataType}}, error) {
//queryResults := make([]{{$v.QueryResultDataType}}, 0)
//sqlStatement := "{{$v.SqlStatement}}"
//if err := dbConnector.ExecuteQuery(sqlStatement,&queryResults,{{$v.WhereKeyValueParamName}});err!=nil{
//return nil,err
//}
//return queryResults,nil
//}
//{{end}}
//{{end}}
////######################################################################################################################
