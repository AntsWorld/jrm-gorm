package jorm

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"strings"
)

func init() {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime)
}

//数据库表中列的元数据
type SchemaTableColumns struct {
	TableCatalog           string `json:"tableCatalog" orm:"TABLE_CATALOG"`                      //表限定符
	TableSchema            string `json:"tableSchema" orm:"TABLE_SCHEMA"`                        //表所有者
	TableName              string `json:"tableName" orm:"TABLE_NAME"`                            //表名
	ColumnName             string `json:"columnName" orm:"COLUMN_NAME"`                          //列名
	OrdinalPosition        int64  `json:"ordinalPosition" orm:"ORDINAL_POSITION"`                //该列在该表中的顺序
	ColumnDefault          string `json:"columnDefault" orm:"COLUMN_DEFAULT"`                    //列的默认值,sql.NullString
	IsNullable             string `json:"isNullable" orm:"IS_NULLABLE"`                          //是否可以为null
	DataType               string `json:"dataType" orm:"DATA_TYPE"`                              //数据类型
	CharacterMaximumLength int64  `json:"characterMaximumLength" orm:"CHARACTER_MAXIMUM_LENGTH"` //数据的长度,sql.NullInt64
	CharacterOctetLength   int64  `json:"characterOctetLength" orm:"CHARACTER_OCTET_LENGTH"`     //数据的存储长度,sql.NullInt64
	NumericPrecision       int64  `json:"numericPrecision" orm:"NUMERIC_PRECISION"`              //精度,sql.NullInt64
	NumericScale           int64  `json:"numericScale" orm:"NUMERIC_SCALE"`                      //小数位数,sql.NullInt64
	DatetimePrecision      int64  `json:"datetimePrecision" orm:"DATETIME_PRECISION"`            //如果列是字符数据或 text 数据类型，那么返回 master，指明字符集所在的数据库,sql.NullInt64
	CharacterSetName       string `json:"characterSetName" orm:"CHARACTER_SET_NAME"`             //如果列是字符数据或 text 数据类型，那么返回 dbo，指明字符集的所有者名称,sql.NullString
	CollationName          string `json:"collationName" orm:"COLLATION_NAME"`                    //如果该列是字符数据或 text 数据类型，那么为字符集返回唯一的名称。否则，返回 null,sql.NullString
	ColumnType             string `json:"columnType" orm:"COLUMN_TYPE"`                          //列的类型，例如varchar(20)
	ColumnKey              string `json:"columnKey" orm:"COLUMN_KEY"`                            //如果等于pri，表示是主键
	EXTRA                  string `json:"extra" orm:"EXTRA"`                                     //义列的时候的其他信息，例如自增，主键
	PRIVILEGES             string `json:"privileges" orm:"PRIVILEGES"`                           //操作权限有：select,insert,update,references
	ColumnComment          string `json:"columnComment" orm:"COLUMN_COMMENT"`                    //列的备注
	GenerationExpression   string `json:"generationExpression" orm:"GENERATION_EXPRESSION"`
}

//表中列信息描述
type LocalSchemaTableColumns struct {
	SchemaTableColumns SchemaTableColumns
	CamelColumnName    string `json:"camelColumnName"`    //列名的驼峰格式
	ColumnDefaultValue string `json:"columnDefaultValue"` //默认值
	Nullable           bool   `json:"nullable"`           //是否可以为空
	GoDataType         string `json:"goDataType"`         //数据类型
}

func (tableInfo *SchemaTableInfo) ConvertTableColumnsToLocalTableColumns(tableColumns []SchemaTableColumns) ([]LocalSchemaTableColumns, error) {
	if tableInfo == nil {
		return nil, errors.New("schemaTableInfo can't be nil")
	}
	if tableColumns == nil || len(tableColumns) == 0 {
		return nil, errors.New("tableColumns can't be nil")
	}
	localSchemaTableColumns := make([]LocalSchemaTableColumns, 0)
	for _, tableColumn := range tableColumns {
		localTableColumn := LocalSchemaTableColumns{}
		localTableColumn.SchemaTableColumns = tableColumn
		//获取
		localTableColumn.CamelColumnName = tableInfo.FmtColumnsNameToCamelName(tableColumn.ColumnName)

		if tableColumn.IsNullable == "NO" {
			localTableColumn.Nullable = false
		} else {
			localTableColumn.Nullable = true
		}
		localTableColumn.GoDataType = tableInfo.MappingMysqlDataTypeToGo(tableColumn.DataType, tableColumn.ColumnType)
		switch localTableColumn.GoDataType {
		case "int", "int64", "uint", "uint64", "float32", "float64":
			if tableColumn.ColumnDefault != "" {
				localTableColumn.ColumnDefaultValue = tableColumn.ColumnDefault
			} else {
				localTableColumn.ColumnDefaultValue = "0"
			}
		case "string":
			if tableColumn.ColumnDefault != "" {
				localTableColumn.ColumnDefaultValue = fmt.Sprintf("\"%s\"", tableColumn.ColumnDefault)
			} else {
				localTableColumn.ColumnDefaultValue = fmt.Sprintf("\"\"")
			}
		case "bool":
			if tableColumn.ColumnDefault == "1" {
				localTableColumn.ColumnDefaultValue = "true"
			} else {
				localTableColumn.ColumnDefaultValue = "false"
			}
		default:
			return nil, errors.New(fmt.Sprintf("unsupported datatype: %s", localTableColumn.GoDataType))
		}
		localSchemaTableColumns = append(localSchemaTableColumns, localTableColumn)
	}

	return localSchemaTableColumns, nil
}

//数据库表信息
type SchemaTableInfo struct {
	SchemaTable  SchemaTable          `json:"schemaTable"`
	TableColumns []SchemaTableColumns `json:"tableColumns"`
}

//基础表模板数据
type BaseTableTemplate struct {
	DatabaseName               string               //数据库名称
	CamelDatabaseName          string               //数据库名的驼峰格式
	TableName                  string               //表名称
	AIColumns                  []SchemaTableColumns //表中所有自增字段
	PriColumns                 []SchemaTableColumns //表中所有PRI字段
	UniqueColumns              []SchemaTableColumns //表中所有Unique字段
	EditableColumns            []SchemaTableColumns //表中所有可编辑字段,插入和更新SQL语句使用
	TableStructureName         string               //表结构体名称,默认使用表名称的驼峰格式
	InsertRequestStructureName string               //Insert请求结构体名称
	UpdateRequestStructureName string               //Update请求结构体名称
}

//获取数据库名称
func (tableInfo *SchemaTableInfo) GetDatabaseName() (string, error) {
	if tableInfo == nil {
		return "", errors.New("tableInfo ptr can't be nil")
	}
	if tableInfo.TableColumns == nil || len(tableInfo.TableColumns) == 0 {
		return "", errors.New("TableColumns can't be nil or empty")
	}
	firstColumn := tableInfo.TableColumns[0]
	if firstColumn.TableSchema == "" {
		return "", errors.New("TableSchema is nil")
	}
	return firstColumn.TableSchema, nil
}

//获取表名称
func (tableInfo *SchemaTableInfo) GetTableName() (string, error) {
	if tableInfo == nil {
		return "", errors.New("tableInfo ptr can't be nil")
	}
	if tableInfo.TableColumns == nil || len(tableInfo.TableColumns) == 0 {
		return "", errors.New("TableColumns can't be nil or empty")
	}
	firstColumn := tableInfo.TableColumns[0]
	if firstColumn.TableName == "" {
		return "", errors.New("TableName is nil")
	}
	return firstColumn.TableName, nil
}

//获取表中自增字段
func (tableInfo *SchemaTableInfo) GetAutoIncrementColumns() ([]SchemaTableColumns, error) {
	if tableInfo == nil {
		return nil, errors.New("tableInfo ptr can't be nil")
	}
	if tableInfo.TableColumns == nil || len(tableInfo.TableColumns) == 0 {
		return nil, errors.New("TableColumns can't be nil or empty")
	}
	autoIncrementColumns := make([]SchemaTableColumns, 0)
	for _, v := range tableInfo.TableColumns {
		if strings.Contains(v.EXTRA, "auto_increment") {
			autoIncrementColumns = append(autoIncrementColumns, v)
		}
	}
	return autoIncrementColumns, nil
}

//获取表中主键字段
func (tableInfo *SchemaTableInfo) GetPrimaryKeyColumns() ([]SchemaTableColumns, error) {
	if tableInfo == nil {
		return nil, errors.New("tableInfo ptr can't be nil")
	}
	if tableInfo.TableColumns == nil || len(tableInfo.TableColumns) == 0 {
		return nil, errors.New("TableColumns can't be nil or empty")
	}
	priColumns := make([]SchemaTableColumns, 0)
	for _, v := range tableInfo.TableColumns {
		if strings.Contains(v.ColumnKey, "PRI") {
			priColumns = append(priColumns, v)
		}
	}
	return priColumns, nil
}

//获取表表中唯一键字段
func (tableInfo *SchemaTableInfo) GetUniqueColumns() ([]SchemaTableColumns, error) {
	if tableInfo == nil {
		return nil, errors.New("tableInfo ptr can't be nil")
	}
	if tableInfo.TableColumns == nil || len(tableInfo.TableColumns) == 0 {
		return nil, errors.New("TableColumns can't be nil or empty")
	}
	uniqueColumns := make([]SchemaTableColumns, 0)
	for _, v := range tableInfo.TableColumns {
		if strings.Contains(v.ColumnKey, "UNI") {
			uniqueColumns = append(uniqueColumns, v)
		}
	}
	return uniqueColumns, nil
}

//获取表中可编辑字段
func (tableInfo *SchemaTableInfo) GetEditableColumns() ([]SchemaTableColumns, error) {
	if tableInfo == nil {
		return nil, errors.New("tableInfo ptr can't be nil")
	}
	if tableInfo.TableColumns == nil || len(tableInfo.TableColumns) == 0 {
		return nil, errors.New("TableColumns can't be nil or empty")
	}
	editableColumns := make([]SchemaTableColumns, 0)
	//除自增、DATETIME外所有字段
	for _, v := range tableInfo.TableColumns {
		if strings.Contains(v.EXTRA, "auto_increment") ||
			strings.Contains(v.ColumnDefault, "CURRENT_TIMESTAMP") ||
			strings.Contains(v.EXTRA, "on update CURRENT_TIMESTAMP") {
			continue
		}
		editableColumns = append(editableColumns, v)
	}
	return editableColumns, nil
}

//根据mysql数据类型获取golang数据类型
func (tableInfo *SchemaTableInfo) MappingMysqlDataTypeToGo(columnDataType string, columnType string) string {
	dt := strings.ToUpper(columnDataType)
	//如果列类型中包含unsigned，表示是无符号类型
	if strings.Contains(columnType, "unsigned") {
		dt = strings.ToUpper("unsigned ") + dt
	}
	vType, exists := MysqlDataTypeMapToGoDataType[dt]
	if !exists {
		log.Println(fmt.Sprintf("数据类型[%s]还未纳入数据类型转换列表中,默认设置为string", columnDataType))
		vType = "string"
	}
	//如果类型为tinyint(1)，golang为bool类型
	if strings.Contains(columnType, "tinyint(1)") {
		vType = "bool"
	}

	return vType
}

//将数据库名称转换为驼峰格式
func (tableInfo *SchemaTableInfo) FmtDatabaseNameToCamelName(value string) string {
	return tableInfo.FmtColumnsNameToCamelName(value)
}

//将表名称转换为结构体名需要的格式
func (tableInfo *SchemaTableInfo) FmtTableNameToStructureName(value string) string {
	return tableInfo.FmtColumnsNameToCamelName(value)
}

//将表名称转换为文件名需要的格式
func (tableInfo *SchemaTableInfo) FmtTableNameToFileName(value string) string {
	return tableInfo.FmtColumnsNameToJsonTagValueName(value)
}

//将表字段名称格式化为结构体驼峰命名方式
func (tableInfo *SchemaTableInfo) FmtColumnsNameToCamelName(value string) string {
	temp := strings.Split(value, "_")
	var str string
	for i := 0; i < len(temp); i++ {
		b := []rune(temp[i])
		for j := 0; j < len(b); j++ {
			if j == 0 {
				if b[j] >= 'a' && b[j] <= 'z' {
					//如果首字母是小写字母，转换为大写字母
					b[j] -= 32
					str += string(b[j])
				} else {
					str += string(b[j])
				}
			} else {
				str += string(b[j])
			}
		}
	}

	return str
}

//表字段名转换为结构体中json tag对应的名称
func (tableInfo *SchemaTableInfo) FmtColumnsNameToJsonTagValueName(value string) string {
	temp := strings.Split(value, "_")
	var str string
	for i := 0; i < len(temp); i++ {
		b := []rune(temp[i])
		//log.Println("b = ", string(b))
		for j := 0; j < len(b); j++ {
			if j == 0 {
				if i == 0 {
					if b[j] >= 'A' && b[j] <= 'Z' {
						//如果首字母是小写字母，转换为大写字母
						b[j] += 32
						str += string(b[j])
					} else {
						str += string(b[j])
					}
				} else {
					if b[j] >= 'a' && b[j] <= 'z' {
						b[j] -= 32
						str += string(b[j])
					} else {
						str += string(b[j])
					}
				}

			} else {
				str += string(b[j])
			}
		}
	}

	return str
}

//生成表结构体
func (tableInfo *SchemaTableInfo) GenerateTableStructure() (string, error) {
	//生成表模板数据
	tableName, err := tableInfo.GetTableName()
	if err != nil {
		return "", err
	}
	maxColumnNameLen := 0
	maxColumnTypeLen := 0
	for _, value := range tableInfo.TableColumns {
		columnName := tableInfo.FmtColumnsNameToCamelName(value.ColumnName)
		columnType := tableInfo.MappingMysqlDataTypeToGo(value.DataType, value.ColumnType)
		if len(columnName) > maxColumnNameLen {
			maxColumnNameLen = len(columnName)
		}
		if len(columnType) > maxColumnTypeLen {
			maxColumnTypeLen = len(columnType)
		}
	}
	structureName := tableInfo.FmtColumnsNameToCamelName(tableName)
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("type %s struct {\n", structureName))

	for _, value := range tableInfo.TableColumns {
		//结构体字段结构
		// Name dataType `json:"table_catalog,defaultValue" isNullAble:"NO" comment:""`
		//DeviceId   string `json:"deviceId" orm:"device_id,NO,varchar(255),UNI" isNullAble:"NO" comment:""`
		//变量名称，将列名称以驼峰方式命名
		columnName := tableInfo.FmtColumnsNameToCamelName(value.ColumnName)
		//数据类型，将数据库类型映射为golang数据类型
		dataType := tableInfo.MappingMysqlDataTypeToGo(value.DataType, value.ColumnType)
		//名称后需要添加的对齐空格
		spaceString := ""
		for i := 0; i < (maxColumnNameLen - len(columnName)); i++ {
			spaceString += " "
		}
		//类型后需要添加的对齐空格
		typeSpaceString := ""
		for i := 0; i < (maxColumnTypeLen - len(dataType)); i++ {
			typeSpaceString += " "
		}
		//json tag值,将列名称首字母小写后作为tag名称
		jsonTagValue := tableInfo.FmtColumnsNameToJsonTagValueName(columnName)
		//`orm:""`
		ormTagValue := value.ColumnName
		//默认值
		if value.ColumnDefault != "" {
			ormTagValue += fmt.Sprintf(",%s", value.ColumnDefault)
		}
		//是否可以为空
		if value.IsNullable == "NO" {
			ormTagValue += fmt.Sprintf(",%s", "NotNullable") //不能为空
		} else {
			ormTagValue += fmt.Sprintf(",%s", "Nullable") //可以为空
		}
		//值类型
		if "" != value.ColumnType {
			ormTagValue += fmt.Sprintf(",%s", value.ColumnType) //列类型
		}
		if "" != value.ColumnKey {
			ormTagValue += fmt.Sprintf(",%s", value.ColumnKey) //ColumnKey
		}
		if "" != value.EXTRA {
			ormTagValue += fmt.Sprintf(",%s", value.EXTRA) //EXTRA
		}
		//是否允许为空
		isNullableValue := fmt.Sprintf("%s", value.IsNullable)
		//字段描述信息
		commentValue := fmt.Sprintf("%s", value.ColumnComment)

		//组装后的数据
		columnStr := fmt.Sprintf("\t%s%s %s%s `json:\"%s\" orm:\"%s\" nullAble:\"%s\" comment:\"%s\"`\n",
			columnName, spaceString, dataType, typeSpaceString, jsonTagValue, ormTagValue, isNullableValue, commentValue)
		buffer.WriteString(columnStr)
	}
	buffer.WriteString("}")

	return buffer.String(), nil
}

//生成基础表模板数据
func (tableInfo *SchemaTableInfo) GenerateBaseTableTemplate() (*BaseTableTemplate, error) {
	//获取数据库名称
	databaseName, err := tableInfo.GetDatabaseName()
	if err != nil {
		return nil, err
	}
	//获取表名称
	tableName, err := tableInfo.GetTableName()
	if err != nil {
		return nil, err
	}
	//获取自增字段
	autoIncrementColumns, err := tableInfo.GetAutoIncrementColumns()
	if err != nil {
		return nil, err
	}
	//获取主键字段
	primaryKeyColumns, err := tableInfo.GetPrimaryKeyColumns()
	if err != nil {
		return nil, err
	}
	//获取唯一键字段
	uniqueColumns, err := tableInfo.GetUniqueColumns()
	if err != nil {
		return nil, err
	}
	//获取唯一键字段
	editableColumns, err := tableInfo.GetEditableColumns()
	if err != nil {
		return nil, err
	}

	baseTableTemplate := BaseTableTemplate{}
	baseTableTemplate.DatabaseName = databaseName
	baseTableTemplate.CamelDatabaseName = tableInfo.FmtColumnsNameToCamelName(databaseName)
	baseTableTemplate.TableName = tableName
	baseTableTemplate.AIColumns = append(baseTableTemplate.AIColumns, autoIncrementColumns...)
	baseTableTemplate.PriColumns = append(baseTableTemplate.PriColumns, primaryKeyColumns...)
	baseTableTemplate.UniqueColumns = append(baseTableTemplate.UniqueColumns, uniqueColumns...)
	baseTableTemplate.EditableColumns = append(baseTableTemplate.EditableColumns, editableColumns...)
	baseTableTemplate.TableStructureName = tableInfo.FmtColumnsNameToCamelName(baseTableTemplate.TableName)
	baseTableTemplate.InsertRequestStructureName = fmt.Sprintf("Insert%sRequest", baseTableTemplate.TableStructureName)
	baseTableTemplate.UpdateRequestStructureName = fmt.Sprintf("Update%sRequest", baseTableTemplate.TableStructureName)
	////生成文件需要的字段
	//baseTableTemplate.PackageName = fmt.Sprintf("%sdb", strings.ToLower(tableInfo.FmtColumnsNameToCamelName(databaseName)))
	//baseTableTemplate.GoFileName = fmt.Sprintf("%s.go", strings.ToLower(tableInfo.FmtTableNameToFileName(tableName)))
	//baseTableTemplate.GoTestFileName = fmt.Sprintf("%s_test.go", strings.ToLower(tableInfo.FmtTableNameToFileName(tableName)))

	return &baseTableTemplate, nil
}
