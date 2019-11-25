package jorm

import (
	"encoding/json"
	"fmt"
	"github.com/AntsWorld/jormsample/jdb"
	"testing"
)

func TestJrOrm_QuerySchemaInfo(t *testing.T) {
	//根据数据库连接信息生成DataSourceName
	dataSourceName, err := jdb.DBConnectInfo.ToDataSourceName()
	if err != nil {
		t.Log(err)
		return
	}
	t.Log(fmt.Sprintf("drivername:%s,dataSourceName:%s", jdb.DriverName, dataSourceName))
	dbConnector := jdb.NewDbConnector()
	if err := dbConnector.Open(jdb.DriverName, dataSourceName); err != nil {
		t.Log(err)
		return
	}
	t.Log("connect database success")
	defer func() {
		if err := dbConnector.Close(); err != nil {
			t.Log(err)
			return
		}
		t.Log("close database connect success")
	}()
	dbConnector.OrmTagName = "orm"
	jrOrm := New(dbConnector)
	//查询数据库信息
	schemaInfo := SchemaInfo{}
	if err := jrOrm.QuerySchemaInfo("hotel_config", &schemaInfo); err != nil {
		t.Log(err)
		return
	}
	data, _ := json.Marshal(&schemaInfo)
	t.Log("SchemaInfo = ", string(data))
	//查询表信息
	if len(schemaInfo.Tables) == 0 {
		return
	}
	schemaTableInfo := SchemaTableInfo{}
	if err := jrOrm.QuerySchemaTableInfo(schemaInfo.Tables[0].TableName, &schemaTableInfo); err != nil {
		t.Log(err)
		return
	}
	schemaTableInfoData, _ := json.Marshal(&schemaTableInfo)
	t.Log("SchemaTableInfo = ", string(schemaTableInfoData))
	//生成结构体
	if value, err := schemaTableInfo.GenerateTableStructure(); err != nil {
		t.Log(err)
		return
	} else {
		t.Log(fmt.Sprintf("结构体：\n%s\n", value))
	}
	//生成插入函数请求参数结构体
	if value, _, _, err := schemaTableInfo.GenerateInsertRequestStructure(); err != nil {
		t.Log(err)
		return
	} else {
		t.Log(fmt.Sprintf("InsertRequest结构体：\n%s\n", value))
	}
	//生成更新函数请求参数结构体
	if value, _, _, err := schemaTableInfo.GenerateUpdateRequestStructure(); err != nil {
		t.Log(err)
		return
	} else {
		t.Log(fmt.Sprintf("UpdateRequest结构体：\n%s\n", value))
	}
}
