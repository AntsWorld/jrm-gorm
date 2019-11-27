package jorm

import (
	"encoding/json"
	"errors"
)

//数据库表元数据
type SchemaTable struct {
	TableCatalog   string `json:"tableCatalog" orm:"TABLE_CATALOG"`     //
	TableSchema    string `json:"tableSchema" orm:"TABLE_SCHEMA"`       //
	TableName      string `json:"tableName" orm:"TABLE_NAME"`           //
	TableType      string `json:"tableType" orm:"TABLE_TYPE"`           //
	ENGINE         string `json:"engine" orm:"ENGINE"`                  //sql.NullString
	VERSION        int64  `json:"version" orm:"VERSION"`                //sql.NullInt64
	RowFormat      string `json:"rowFormat" orm:"ROW_FORMAT"`           //sql.NullString
	TableRows      int64  `json:"tableRows" orm:"TABLE_ROWS"`           //sql.NullInt64
	AvgRowLength   int64  `json:"avgRowLength" orm:"AVG_ROW_LENGTH"`    //sql.NullInt64
	DataLength     int64  `json:"dataLength" orm:"DATA_LENGTH"`         //sql.NullInt64
	MaxDataLength  int64  `json:"maxDataLength" orm:"MAX_DATA_LENGTH"`  //sql.NullInt64
	IndexLength    int64  `json:"indexLength" orm:"INDEX_LENGTH"`       //sql.NullInt64
	DataFree       int64  `json:"dataFree" orm:"DATA_FREE"`             //sql.NullInt64
	AutoIncrement  int64  `json:"autoIncrement" orm:"AUTO_INCREMENT"`   //sql.NullInt64
	CreateTime     string `json:"createTime" orm:"CREATE_TIME"`         //sql.NullString
	UpdateTime     string `json:"updateTime" orm:"UPDATE_TIME"`         //sql.NullString
	CheckTime      string `json:"checkTime" orm:"CHECK_TIME"`           //sql.NullString
	TableCollation string `json:"tableCollation" orm:"TABLE_COLLATION"` //sql.NullString
	CHECKSUM       int64  `json:"checksum" orm:"CHECKSUM"`              //sql.NullInt64
	CreateOptions  string `json:"createOptions" orm:"CREATE_OPTIONS"`   //sql.NullString
	TableComment   string `json:"tableComment" orm:"TABLE_COMMENT"`
}

//数据库信息
type SchemaInfo struct {
	Tables []SchemaTable `json:"tables"` //数据库中所有的表
}

//查询数据库中某一张表的信息
func (schemaInfo *SchemaInfo) QueryTableInfoByPosition(orm *JrOrm, position int) ([]SchemaTableInfo, error) {
	if orm == nil {
		return nil, errors.New("orm param can't be nil")
	}
	if schemaInfo == nil {
		return nil, errors.New("schemaInfo can't be nil")
	}
	if schemaInfo.Tables == nil || len(schemaInfo.Tables) == 0 {
		return nil, errors.New("schemaInfo Tables can't be nil or empty")
	}
	if position < 0 || position > len(schemaInfo.Tables)-1 {
		return nil, errors.New("position out of range")
	}
	table := schemaInfo.Tables[position]
	schemaTableInfos := make([]SchemaTableInfo,0)
	schemaTableInfo:=SchemaTableInfo{}
	if err := orm.QuerySchemaTableInfo(table, &schemaTableInfo); err != nil {
		return nil, err
	}
	schemaTableInfo.SchemaTable = table
	schemaTableInfos = append(schemaTableInfos,schemaTableInfo)
	return schemaTableInfos, nil
}

//查询数据库中某一张表的信息
func (schemaInfo *SchemaInfo) QueryOneTableInfo(orm *JrOrm, table SchemaTable) ([]SchemaTableInfo, error) {
	if orm == nil {
		return nil, errors.New("orm param can't be nil")
	}
	if schemaInfo == nil {
		return nil, errors.New("schemaInfo can't be nil")
	}
	if schemaInfo.Tables == nil || len(schemaInfo.Tables) == 0 {
		return nil, errors.New("schemaInfo Tables can't be nil or empty")
	}
	contains := false
	for _, v := range schemaInfo.Tables {
		vData, err := json.Marshal(&v)
		if err != nil {
			return nil, err
		}
		tableData, err := json.Marshal(&table)
		if err != nil {
			return nil, err
		}
		if string(vData) == string(tableData) {
			contains = true
			break
		}
	}
	if !contains {
		return nil, errors.New("table not belong this database")
	}

	schemaTableInfos := make([]SchemaTableInfo,0)
	schemaTableInfo:=SchemaTableInfo{}
	if err := orm.QuerySchemaTableInfo(table, &schemaTableInfo); err != nil {
		return nil, err
	}
	schemaTableInfo.SchemaTable = table
	schemaTableInfos = append(schemaTableInfos,schemaTableInfo)
	return schemaTableInfos, nil
}

//查询数据库中所有表的信息
func (schemaInfo *SchemaInfo) QueryAllTableInfo(orm *JrOrm) ([]SchemaTableInfo, error) {
	if orm == nil {
		return nil, errors.New("orm param can't be nil")
	}
	if schemaInfo == nil {
		return nil, errors.New("schemaInfo can't be nil")
	}
	if schemaInfo.Tables == nil || len(schemaInfo.Tables) == 0 {
		return nil, errors.New("schemaInfo Tables can't be nil or empty")
	}
	schemaTableInfos := make([]SchemaTableInfo, 0)
	for _, table := range schemaInfo.Tables {
		//tableName := table.TableName
		schemaTableInfo := SchemaTableInfo{}
		if err := orm.QuerySchemaTableInfo(table, &schemaTableInfo); err != nil {
			return nil, err
		}
		schemaTableInfo.SchemaTable = table
		schemaTableInfos = append(schemaTableInfos, schemaTableInfo)
	}
	return schemaTableInfos, nil
}
