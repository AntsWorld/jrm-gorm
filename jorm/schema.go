package jorm

//数据库信息
type SchemaInfo struct {
	Tables []SchemaTable `json:"tables"` //数据库中所有的表
}

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
