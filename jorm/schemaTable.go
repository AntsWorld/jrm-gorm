package jorm

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"strings"
)

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

//数据库表信息
type SchemaTableInfo struct {
	SchemaTable          SchemaTable          `json:"schemaTable"`
	TableColumns         []SchemaTableColumns `json:"tableColumns"`
	BaseTableTemplate    BaseTableTemplate    `json:"baseTableTemplate,omitempty"`
	DefaultTableTemplate DefaultTableTemplate `json:"defaultTableTemplate,omitempty"`
}

//删除函数模板参数
type DeleteTemplate struct {
	FunctionName           string
	SqlStatement           string
	WhereKeyDataType       string
	WhereKeyValueParamName string
}

//查询模板参数
type QueryTemplate struct {
	UniQuery               bool //是否是唯一键查询
	FunctionName           string
	SqlStatement           string
	WhereKeyDataType       string
	WhereKeyValueParamName string
	QueryResultDataType    string
}

//将数据库数据类型转换为GO数据类型
func (tableInfo *SchemaTableInfo) MappingMysqlDataTypeToGo(dataType string, columnType string) string {
	dt := strings.ToUpper(dataType)
	//如果列类型中包含unsigned，表示是无符号类型
	if strings.Contains(columnType, "unsigned") {
		dt = strings.ToUpper("unsigned ") + dt
	}
	vType, exists := MysqlDataTypeMapToGoDataType[dt]
	if !exists {
		log.Println(fmt.Sprintf("数据类型[%s]还未纳入数据类型转换列表中,默认设置为string", dataType))
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

//生成基础表模板数据
func (tableInfo *SchemaTableInfo) GenerateBaseTableTemplate() (*BaseTableTemplate, error) {
	if len(tableInfo.TableColumns) == 0 {
		return nil, errors.New("table columns can't be nil")
	}
	firstColumn := tableInfo.TableColumns[0]
	baseTableTemplate := &BaseTableTemplate{}
	//表基础字段
	baseTableTemplate.DatabaseName = firstColumn.TableSchema
	baseTableTemplate.TableName = firstColumn.TableName
	//主键
	for _, v := range tableInfo.TableColumns {
		if strings.Contains(v.ColumnKey, "PRI") {
			baseTableTemplate.PriColumns = append(baseTableTemplate.PriColumns, v)
		}
	}
	//唯一键
	for _, v := range tableInfo.TableColumns {
		if strings.Contains(v.ColumnKey, "UNI") {
			baseTableTemplate.UniqueColumns = append(baseTableTemplate.UniqueColumns, v)
		}
	}
	//可编辑Columns,除去自增字段、创建时间、更新时间外都认为是可编辑的
	for _, v := range tableInfo.TableColumns {
		if strings.Contains(v.EXTRA, "auto_increment") ||
			strings.Contains(v.ColumnDefault, "CURRENT_TIMESTAMP") ||
			strings.Contains(v.EXTRA, "on update CURRENT_TIMESTAMP") {
			continue
		}
		baseTableTemplate.EditableColumns = append(baseTableTemplate.EditableColumns, v)
	}
	//生成文件需要的字段
	baseTableTemplate.PackageName = strings.ToLower(tableInfo.FmtDatabaseNameToCamelName(firstColumn.TableSchema))
	baseTableTemplate.GoFileName = fmt.Sprintf("%s.go", tableInfo.FmtTableNameToFileName(firstColumn.TableName))
	baseTableTemplate.GoTestFileName = fmt.Sprintf("%s_test.go",
		tableInfo.FmtTableNameToFileName(firstColumn.TableName))
	baseTableTemplate.TableStructureName = tableInfo.FmtTableNameToStructureName(baseTableTemplate.TableName)

	tableInfo.BaseTableTemplate = *baseTableTemplate
	return baseTableTemplate, nil
}

//生成表模板数据
func (tableInfo *SchemaTableInfo) GenerateTableTemplate() (*DefaultTableTemplate, error) {
	//生成基础表模板数据
	baseTableTemplate, err := tableInfo.GenerateBaseTableTemplate()
	if err != nil {
		return nil, err
	}
	tableTemplate := &DefaultTableTemplate{}
	//基础表模板数据
	tableTemplate.BaseTableTemplate = *baseTableTemplate
	//生成表结构体需要的字段
	if structure, err := tableInfo.GenerateTableStructure(); err != nil {
		return nil, err
	} else {
		tableTemplate.TableStructure = structure
	}
	//插入函数模板数据
	if templates, err := tableInfo.GenerateInsertTemplates(); err != nil {
		return nil, err
	} else {
		tableTemplate.InsertTemplates = templates
	}
	//删除函数模板数据
	if templates, err := tableInfo.GenerateDeleteTemplates(); err != nil {
		return nil, err
	} else {
		tableTemplate.DeleteTemplates = templates
	}
	//查询函数模板数据
	if templates, err := tableInfo.GenerateQueryTemplates(); err != nil {
		return nil, err
	} else {
		tableTemplate.QueryTemplates = templates
	}
	//更新函数模板数据
	if templates, err := tableInfo.GenerateUpdateTemplates(); err != nil {
		return nil, err
	} else {
		tableTemplate.UpdateTemplates = templates
	}

	tableInfo.DefaultTableTemplate = *tableTemplate
	return tableTemplate, nil
}

//生成表结构体
func (tableInfo *SchemaTableInfo) GenerateTableStructure() (string, error) {
	//生成表模板数据
	tableTemplate, err := tableInfo.GenerateBaseTableTemplate()
	if err != nil {
		return "", err
	}
	//基础参数
	structureName := tableTemplate.TableStructureName

	var buffer bytes.Buffer
	//buffer.WriteString(fmt.Sprintf("package %s\n", packageName))
	//buffer.WriteString(fmt.Sprintf("var dbNameOf%s = \"%s\"\n", structureName, databaseName))
	//buffer.WriteString(fmt.Sprintf("var tableNameOf%s = \"%s\"\n", structureName, tableName))
	buffer.WriteString(fmt.Sprintf("type %s struct {\n", structureName))
	//获取字段名称最大长度,做格式对齐使用
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
		columnStr := fmt.Sprintf("\t%s%s %s%s `json:\"%s\" orm:\"%s\" isNullAble:\"%s\" comment:\"%s\"`\n",
			columnName, spaceString, dataType, typeSpaceString, jsonTagValue, ormTagValue, isNullableValue, commentValue)
		buffer.WriteString(columnStr)
	}
	buffer.WriteString("}")

	return buffer.String(), nil
}

//生成插入函数请求结构体
func (tableInfo *SchemaTableInfo) GenerateInsertRequestStructure() (string, []string, []string, error) {
	//生成表模板数据
	baseTableTemplate, err := tableInfo.GenerateBaseTableTemplate()
	if err != nil {
		return "", nil, nil, err
	}
	//结构体名称
	requestStructureName := CreateInsertRequestStructureName(baseTableTemplate.TableStructureName)
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("type %s struct {\n", requestStructureName))
	//获取字段名称最大长度,做格式对齐使用
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
	tableColumns := make([]string, 0)
	columns := make([]string, 0)
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
		columnStr := fmt.Sprintf("\t%s%s %s%s `json:\"%s\" orm:\"%s\" isNullAble:\"%s\" comment:\"%s\"`\n",
			columnName, spaceString, dataType, typeSpaceString, jsonTagValue, ormTagValue, isNullableValue, commentValue)
		//自增字段、datetime字段不需要包含到Insert请求参数中
		if strings.Contains(ormTagValue, "auto_increment") ||
			strings.Contains(ormTagValue, "CURRENT_TIMESTAMP") {
			continue
		}
		tableColumns = append(tableColumns, value.ColumnName)
		columns = append(columns, columnName)

		buffer.WriteString(columnStr)
	}
	buffer.WriteString("}\n")

	return buffer.String(), tableColumns, columns, nil
}

//生成更新函数请求结构体
func (tableInfo *SchemaTableInfo) GenerateUpdateRequestStructure() (string,
	[]string, []string, error) {
	//生成表模板数据
	baseTableTemplate, err := tableInfo.GenerateBaseTableTemplate()
	if err != nil {
		return "", nil, nil, err
	}
	//结构体名称
	requestStructureName := CreateUpdateRequestStructureName(baseTableTemplate.TableStructureName)
	//基础参数
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("type %s struct {\n", requestStructureName))
	//获取字段名称最大长度,做格式对齐使用
	maxColumnNameLen := 0
	maxColumnTypeLen := 0
	for _, value := range baseTableTemplate.EditableColumns {
		columnName := tableInfo.FmtColumnsNameToCamelName(value.ColumnName)
		columnType := tableInfo.MappingMysqlDataTypeToGo(value.DataType, value.ColumnType)
		if len(columnName) > maxColumnNameLen {
			maxColumnNameLen = len(columnName)
		}
		if len(columnType) > maxColumnTypeLen {
			maxColumnTypeLen = len(columnType)
		}
	}
	tableColumns := make([]string, 0)
	columns := make([]string, 0)
	for _, value := range baseTableTemplate.EditableColumns {
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
		columnStr := fmt.Sprintf("\t%s%s %s%s `json:\"%s\" orm:\"%s\" isNullAble:\"%s\" comment:\"%s\"`\n",
			columnName, spaceString, dataType, typeSpaceString, jsonTagValue, ormTagValue, isNullableValue, commentValue)

		tableColumns = append(tableColumns, value.ColumnName)
		columns = append(columns, columnName)

		buffer.WriteString(columnStr)
	}
	buffer.WriteString("}")

	return buffer.String(), tableColumns, columns, nil
}

//生成插入函数请求参数结构体名称
func CreateInsertRequestStructureName(tableStructureName string) string {
	return fmt.Sprintf("InsertNew%sRowRequest", tableStructureName)
}

//生成更新函数请求参数结构体名称
func CreateUpdateRequestStructureName(tableStructureName string) string {
	return fmt.Sprintf("Update%sRequest", tableStructureName)
}

//生成插入函数模板数据
func (tableInfo *SchemaTableInfo) GenerateInsertTemplates() ([]InsertTemplate, error) {
	//生成基础表模板数据
	baseTableTemplate, err := tableInfo.GenerateBaseTableTemplate()
	if err != nil {
		return nil, err
	}
	templates := make([]InsertTemplate, 0)
	//生成Insert请求结构体
	requestStructureName := CreateInsertRequestStructureName(baseTableTemplate.TableStructureName)
	requestStructure, tableColumns, structureColumns, err := tableInfo.GenerateInsertRequestStructure()
	if err != nil {
		return nil, err
	}
	//生成Insert SQL语句
	//"INSERT INTO hotel_config.devices (device_id,device_name,hotel_name,config) VALUES (?,?,?,?);"
	insertColumns := ""
	insertValuesPlace := ""
	for _, value := range tableColumns {
		if insertColumns == "" {
			insertColumns += value
		} else {
			insertColumns += "," + value
		}
		if insertValuesPlace == "" {
			insertValuesPlace += "?"
		} else {
			insertValuesPlace += ",?"
		}
	}
	sqlStatement := fmt.Sprintf("INSERT INTO %s.%s (%s) VALUES (%s);", baseTableTemplate.DatabaseName,
		baseTableTemplate.TableName, insertColumns, insertValuesPlace)

	//函数名
	functionName := fmt.Sprintf("InsertNewRowInto%s", baseTableTemplate.TableStructureName)

	insertTemplate := InsertTemplate{}
	insertTemplate.InsertRequestStructureName = requestStructureName
	insertTemplate.InsertRequestStructure = requestStructure
	insertTemplate.FunctionName = functionName
	insertTemplate.InsertSqlStatement = sqlStatement
	insertTemplate.InsertSqlTableColumns = tableColumns
	insertTemplate.InsertSqlStructureColumns = structureColumns

	templates = append(templates, insertTemplate)

	return templates, nil
}

//生成删除函数模板数据
func (tableInfo *SchemaTableInfo) GenerateDeleteTemplates() ([]DeleteTemplate, error) {
	//生成基础表模板数据
	baseTableTemplate, err := tableInfo.GenerateBaseTableTemplate()
	if err != nil {
		return nil, err
	}
	templates := make([]DeleteTemplate, 0)
	//生成删除模板数据
	if len(baseTableTemplate.PriColumns) == 0 && len(baseTableTemplate.UniqueColumns) == 0 {
		return nil, errors.New("table must contains pri or uni columns")
	}
	//获取可用户删除判断的列信息
	columns := make([]SchemaTableColumns, 0)
	columns = append(columns, baseTableTemplate.PriColumns...)
	columns = append(columns, baseTableTemplate.EditableColumns...)
	//遍历列切片
	for _, column := range columns {
		//生成UpdateSQL语句
		//"DELETE FROM {{$v.DatabaseName}}.{{$v.TableName}} WHERE {{$v.TableColumnName}}=?;"
		whereKeyColumnName := column.ColumnName                                                          // 条件字段
		whereKeyColumnDataType := tableInfo.MappingMysqlDataTypeToGo(column.DataType, column.ColumnType) //条件字段数据类型
		whereKeyDataTypeParamName := tableInfo.FmtColumnsNameToJsonTagValueName(whereKeyColumnName)      //条件字段变量名
		sqlStatement := fmt.Sprintf("DELETE FROM %s.%s WHERE %s=?;", baseTableTemplate.DatabaseName,
			baseTableTemplate.TableName, whereKeyColumnName)
		//函数名
		functionName := fmt.Sprintf("Delete%sRowsBy%s", baseTableTemplate.TableStructureName,
			tableInfo.FmtColumnsNameToCamelName(whereKeyColumnName))
		deleteTemplate := DeleteTemplate{}
		deleteTemplate.FunctionName = functionName
		deleteTemplate.SqlStatement = sqlStatement
		deleteTemplate.WhereKeyDataType = whereKeyColumnDataType
		deleteTemplate.WhereKeyValueParamName = whereKeyDataTypeParamName
		templates = append(templates, deleteTemplate)
	}

	return templates, nil
}

//生成查询函数模板数据
func (tableInfo *SchemaTableInfo) GenerateQueryTemplates() ([]QueryTemplate, error) {
	//生成基础表模板数据
	baseTableTemplate, err := tableInfo.GenerateBaseTableTemplate()
	if err != nil {
		return nil, err
	}
	templates := make([]QueryTemplate, 0)
	//生成删除模板数据
	if len(baseTableTemplate.PriColumns) == 0 && len(baseTableTemplate.UniqueColumns) == 0 {
		return nil, errors.New("table must contains pri or uni columns")
	}
	//获取可用户删除判断的列信息
	columns := make([]SchemaTableColumns, 0)
	columns = append(columns, baseTableTemplate.PriColumns...)
	columns = append(columns, baseTableTemplate.EditableColumns...)
	//遍历列切片
	for _, column := range columns {

		//生成Select SQL语句
		//"SELECT * FROM {{$v.DatabaseName}}.{{$v.TableName}} WHERE {{$v.TableColumnName}}=?;"
		whereKeyColumnName := column.ColumnName                                                          // 条件字段
		whereKeyColumnDataType := tableInfo.MappingMysqlDataTypeToGo(column.DataType, column.ColumnType) //条件字段数据类型
		whereKeyDataTypeParamName := tableInfo.FmtColumnsNameToJsonTagValueName(whereKeyColumnName)      //条件字段变量名
		sqlStatement := fmt.Sprintf("SELECT * FROM %s.%s WHERE %s=?;", baseTableTemplate.DatabaseName,
			baseTableTemplate.TableName, whereKeyColumnName)
		//函数名
		functionName := fmt.Sprintf("Query%sBy%s", baseTableTemplate.TableStructureName,
			tableInfo.FmtColumnsNameToCamelName(whereKeyColumnName))
		queryTemplate := QueryTemplate{}
		queryTemplate.FunctionName = functionName
		queryTemplate.SqlStatement = sqlStatement
		queryTemplate.WhereKeyDataType = whereKeyColumnDataType
		queryTemplate.WhereKeyValueParamName = whereKeyDataTypeParamName
		queryTemplate.QueryResultDataType = baseTableTemplate.TableStructureName
		if strings.Contains(column.ColumnKey, "UNI") {
			queryTemplate.UniQuery = true
		}
		templates = append(templates, queryTemplate)
	}

	return templates, nil
}

//生成更新函数
func (tableInfo *SchemaTableInfo) GenerateUpdateTemplates() ([]UpdateTemplate, error) {
	//生成基础表模板数据
	baseTableTemplate, err := tableInfo.GenerateBaseTableTemplate()
	if err != nil {
		return nil, err
	}
	templates := make([]UpdateTemplate, 0)
	//##################################################################################################################
	//1.生成一次更新所有可更新字段函数模板数据,使用主键做更新条件判断,如果没有主键则使用唯一键，如果没有唯一健则提示错误
	if len(baseTableTemplate.PriColumns) == 0 && len(baseTableTemplate.UniqueColumns) == 0 {
		return nil, errors.New("table must contains pri or uni columns")
	}
	//更新请求结构体名称
	requestStructureName := CreateUpdateRequestStructureName(baseTableTemplate.TableStructureName)
	//生成更新请求结构体
	requestStructure, tableColumns, structureColumns, err := tableInfo.GenerateUpdateRequestStructure()
	if err != nil {
		return nil, err
	}
	//使用主键和唯一键作为where条件更新表信息
	priAndUniColumns := make([]SchemaTableColumns, 0)
	priAndUniColumns = append(priAndUniColumns, baseTableTemplate.PriColumns...)
	priAndUniColumns = append(priAndUniColumns, baseTableTemplate.UniqueColumns...)
	for _, v := range priAndUniColumns {
		sqlColumns := "" //需要更新哪些字段
		for _, value := range tableColumns {
			if sqlColumns == "" {
				sqlColumns += fmt.Sprintf("%s=?", value)
			} else {
				sqlColumns += fmt.Sprintf(",%s=?", value)
			}
		}

		//生成UpdateSQL语句
		//"UPDATE `jrm_manager`.`tb_user` SET user_name=?,phone_number=?,password=?,email=?,avatar=?,sex=?,enable=? WHERE uid=?;"
		whereKeyColumnName := v.ColumnName                                                          // 条件字段
		whereKeyColumnDataType := tableInfo.MappingMysqlDataTypeToGo(v.DataType, v.ColumnType)      //条件字段数据类型
		whereKeyDataTypeParamName := tableInfo.FmtColumnsNameToJsonTagValueName(whereKeyColumnName) //条件字段变量名
		sqlStatement := fmt.Sprintf("UPDATE %s.%s SET %s WHERE %s=?;", baseTableTemplate.DatabaseName,
			baseTableTemplate.TableName, sqlColumns, whereKeyColumnName)
		//函数名
		functionName := fmt.Sprintf("Update%sBy%s", baseTableTemplate.TableStructureName,
			tableInfo.FmtColumnsNameToCamelName(whereKeyColumnName))
		//组装数据
		updateAllTemplate := UpdateTemplate{}
		updateAllTemplate.RequestStructureName = requestStructureName
		updateAllTemplate.RequestStructure = requestStructure
		updateAllTemplate.FunctionName = functionName
		updateAllTemplate.UpdateSqlStatement = sqlStatement
		updateAllTemplate.UpdateSqlTableColumns = tableColumns
		updateAllTemplate.UpdateSqlStructureColumns = structureColumns
		updateAllTemplate.WhereKeyName = whereKeyColumnName
		updateAllTemplate.WhereKeyDataType = whereKeyColumnDataType
		updateAllTemplate.WhereKeyValueParamName = whereKeyDataTypeParamName
		templates = append(templates, updateAllTemplate)
	}
	//使用主键和唯一键更新单个字段
	for _, priColumn := range priAndUniColumns {
		for _, v := range baseTableTemplate.EditableColumns {
			//生成UpdateSQL语句
			//"UPDATE `jrm_manager`.`tb_user` SET user_name=?,phone_number=?,password=?,email=?,avatar=?,sex=?,enable=? WHERE uid=?;"
			whereKeyName := priColumn.ColumnName                                                             // 条件字段
			whereKeyDataType := tableInfo.MappingMysqlDataTypeToGo(priColumn.DataType, priColumn.ColumnType) //条件字段数据类型
			whereKeyDataTypeParamName := tableInfo.FmtColumnsNameToJsonTagValueName(whereKeyName)            //条件字段变量名
			sqlStatement := fmt.Sprintf("UPDATE %s.%s SET %s WHERE %s=?;", baseTableTemplate.DatabaseName,
				baseTableTemplate.TableName, v.ColumnName+"=?", whereKeyName)
			//函数名
			functionName := fmt.Sprintf("Update%s%sBy%s", baseTableTemplate.TableStructureName,
				tableInfo.FmtColumnsNameToCamelName(v.ColumnName),
				tableInfo.FmtColumnsNameToCamelName(whereKeyName))
			structureColumns := make([]string, 0)
			structureColumns = append(structureColumns, tableInfo.FmtColumnsNameToCamelName(v.ColumnName))
			//组装数据
			updateAllTemplate := UpdateTemplate{}
			//updateAllTemplate.RequestStructureName = requestStructureName
			//updateAllTemplate.RequestStructure = requestStructure
			updateAllTemplate.FunctionName = functionName
			updateAllTemplate.UpdateSqlStatement = sqlStatement
			updateAllTemplate.UpdateSqlTableColumns = tableColumns
			updateAllTemplate.UpdateSqlStructureColumns = structureColumns
			updateAllTemplate.WhereKeyName = whereKeyName
			updateAllTemplate.WhereKeyDataType = whereKeyDataType
			updateAllTemplate.WhereKeyValueParamName = whereKeyDataTypeParamName
			templates = append(templates, updateAllTemplate)
		}
	}

	return templates, nil
}
