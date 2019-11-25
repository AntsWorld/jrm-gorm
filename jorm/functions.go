package jorm

////######################################################################################################################
////查询数据库中所有的的表
//func (orm *JrOrm) SelectSchemaTables(dbName string) ([]SchemaTable, error) {
//	sqlCmd := "select * from information_schema.tables where table_schema=? and table_type='base table';"
//	//db := orm.DBConnector.DB
//	//stmt, err := db.Prepare(sqlCmd)
//	//if err != nil {
//	//	return nil, err
//	//}
//	//defer func() {
//	//	if err := stmt.Close(); err != nil {
//	//		log.Println(err)
//	//	}
//	//}()
//	//
//	////rows, err := stmt.Query(dbName)
//	rows, err := orm.DBConnector.ExecuteQueryForRows(sqlCmd, dbName)
//	if err != nil {
//		return nil, err
//	}
//
//	defer func() {
//		if err := rows.Close(); err != nil {
//			log.Println(err)
//		}
//	}()
//
//	var tables []SchemaTable
//
//	for rows.Next() {
//		table := SchemaTable{}
//		err := rows.Scan(
//			&table.TableCatalog,
//			&table.TableSchema,
//			&table.TableName,
//			&table.TableType,
//			&table.ENGINE,
//			&table.VERSION,
//			&table.RowFormat,
//			&table.TableRows,
//			&table.AvgRowLength,
//			&table.DataLength,
//			&table.MaxDataLength,
//			&table.IndexLength,
//			&table.DataFree,
//			&table.AutoIncrement,
//			&table.CreateTime,
//			&table.UpdateTime,
//			&table.CheckTime,
//			&table.TableCollation,
//			&table.CHECKSUM,
//			&table.CreateOptions,
//			&table.TableComment)
//		if err != nil {
//			log.Println(err)
//			continue
//		}
//		//log.Println(fmt.Sprintf("%+v", table))
//		tables = append(tables, table)
//	}
//	return tables, nil
//}
//
////查询表中所有的列信息
//func (orm *JrOrm) SelectSchemaTableColumns(tbName string) ([]SchemaTableColumns, error) {
//	sqlCmd := "select * from information_schema.columns where table_name =?;"
//	//db := orm.DBConnector.DB
//	//stmt, err := db.Prepare(sqlCmd)
//	//if err != nil {
//	//	return nil, err
//	//}
//	//defer func() {
//	//	if err := stmt.Close(); err != nil {
//	//		log.Println(err)
//	//	}
//	//}()
//	//rows, err := stmt.Query(tbName)
//	rows, err := orm.DBConnector.ExecuteQueryForRows(sqlCmd, tbName)
//	if err != nil {
//		return nil, err
//	}
//
//	defer func() {
//		if err := rows.Close(); err != nil {
//			log.Println(err)
//		}
//	}()
//
//	var tableColumns []SchemaTableColumns
//
//	for rows.Next() {
//		column := SchemaTableColumns{}
//		err := rows.Scan(
//			&column.TableCatalog,
//			&column.TableSchema,
//			&column.TableName,
//			&column.ColumnName,
//			&column.OrdinalPosition,
//			&column.ColumnDefault,
//			&column.IsNullable,
//			&column.DataType,
//			&column.CharacterMaximumLength,
//			&column.CharacterOctetLength,
//			&column.NumericPrecision,
//			&column.NumericScale,
//			&column.DatetimePrecision,
//			&column.CharacterSetName,
//			&column.CollationName,
//			&column.ColumnType,
//			&column.ColumnKey,
//			&column.EXTRA,
//			&column.PRIVILEGES,
//			&column.ColumnComment,
//			&column.GenerationExpression)
//		if err != nil {
//			log.Println(err)
//			continue
//		}
//		//log.Println(fmt.Sprintf("%+v", column))
//		tableColumns = append(tableColumns, column)
//	}
//	return tableColumns, nil
//}
//
////根据数据库表生成Go结构体
//func (orm *JrOrm) GenerateStructureBySchemaTable(schemaTable *SchemaTable, tableColumns []SchemaTableColumns) (string, error) {
//	if nil == schemaTable {
//		return "", errors.New("SchemaTable is nil")
//	}
//	if nil == tableColumns || len(tableColumns) == 0 {
//		return "", errors.New("SchemaTableColumns is nil")
//	}
//
//	var buffer bytes.Buffer
//	////包名,以数据库名称命名
//	//packageName := strings.ToLower(FmtColumnsNameToCamelName(schemaTable.TableSchema))
//	////表名称
//	//tableName := tableColumns[0].TableName
//	structureName := FmtColumnsNameToCamelName(tableColumns[0].TableName)
//	//buffer.WriteString(fmt.Sprintf("package %s\n", packageName))
//	//buffer.WriteString(fmt.Sprintf("var dbNameOf%s = \"%s\"\n", structureName, schemaTable.TableSchema))
//	//buffer.WriteString(fmt.Sprintf("var tableNameOf%s = \"%s\"\n", structureName, tableName))
//	//buffer.WriteString(fmt.Sprintf("//%s\n", schemaTable.TableComment))
//	buffer.WriteString(fmt.Sprintf("type %s struct {\n", structureName))
//	//for _, value := range tableColumns {
//	//	//结构体字段结构：Name dataType `json:"table_catalog,defaultValue" isNullAble:"NO" comment:""`
//	//	//变量名称，将列名称以驼峰方式命名
//	//	columnName := FmtColumnsNameToCamelName(value.ColumnName)
//	//	//数据类型，将数据库类型映射为golang数据类型
//	//	dataType := MappingMysqlDataTypeToGo(value.DataType, value.ColumnType)
//	//	//json tag值,将列名称首字母小写后作为tag名称
//	//	jsonTagValue := fmt.Sprintf("%s", FmtColumnsNameToJsonTagValueName(columnName))
//	//	//`orm:""`
//	//	ormTagValue := fmt.Sprintf("%s", value.ColumnName)
//	//	if value.ColumnDefault.Valid {
//	//		columnDefault := value.ColumnDefault.String //默认值
//	//		ormTagValue += fmt.Sprintf(",%s", columnDefault)
//	//	}
//	//	ormTagValue += fmt.Sprintf(",%s", value.IsNullable) //是否可为空
//	//	if "" != value.ColumnType {
//	//		ormTagValue += fmt.Sprintf(",%s", value.ColumnType) //列类型
//	//	}
//	//	if "" != value.ColumnKey {
//	//		ormTagValue += fmt.Sprintf(",%s", value.ColumnKey) //ColumnKey
//	//	}
//	//	if "" != value.EXTRA {
//	//		ormTagValue += fmt.Sprintf(",%s", value.EXTRA) //EXTRA
//	//	}
//	//	//是否允许为空
//	//	isNullableValue := fmt.Sprintf("%s", value.IsNullable)
//	//	//字段描述信息
//	//	commentValue := fmt.Sprintf("%s", value.ColumnComment)
//	//
//	//	//组装后的数据
//	//	columnStr := fmt.Sprintf("\t%s %s `json:\"%s\" orm:\"%s\" isNullAble:\"%s\" comment:\"%s\"`\n",
//	//		columnName, dataType, jsonTagValue, ormTagValue, isNullableValue, commentValue)
//	//	buffer.WriteString(columnStr)
//	//}
//	buffer.WriteString("}\n")
//
//	return buffer.String(), nil
//}
//
////根据数据库表生成数据库基础操作需要的模板数据
//func (orm *JrOrm) GenerateInsertTemplateData(schemaTable *SchemaTable, tableColumns []SchemaTableColumns) (*InsertTemplate, error) {
//	if nil == schemaTable {
//		return nil, errors.New("SchemaTable is nil")
//	}
//	if nil == tableColumns || len(tableColumns) == 0 {
//		return nil, errors.New("SchemaTableColumns is nil")
//	}
//	//Insert操作数据
//	var insertDataBuffer bytes.Buffer
//	//插入语句请求参数
//	insertRequestStructureName := fmt.Sprintf("Insert%sRequest", FmtColumnsNameToCamelName(tableColumns[0].TableName))
//	insertFuncRequestParamName := "item"
//	insertSqlColumnNames := ""
//	insertSqlColumnPlaces := ""
//	insertSqlColumnValues := ""
//
//	insertDataBuffer.WriteString(fmt.Sprintf("type %s struct {\n", insertRequestStructureName))
//	for _, value := range tableColumns {
//		//结构体字段结构：Name dataType `json:"table_catalog,defaultValue" isNullAble:"NO" comment:""`
//		//变量名称，将列名称以驼峰方式命名
//		columnName := FmtColumnsNameToCamelName(value.ColumnName)
//		//数据类型，将数据库类型映射为golang数据类型
//		dataType := MappingMysqlDataTypeToGo(value.DataType, value.ColumnType)
//		//json tag值,将列名称首字母小写后作为tag名称
//		jsonTagValue := fmt.Sprintf("%s", FmtColumnsNameToJsonTagValueName(columnName))
//		//`orm:""`
//		ormTagValue := fmt.Sprintf("%s", value.ColumnName)
//		//if value.ColumnDefault.Valid {
//		//	columnDefault := value.ColumnDefault.String //默认值
//		//	ormTagValue += fmt.Sprintf(",%s", columnDefault)
//		//}
//		ormTagValue += fmt.Sprintf(",%s", value.IsNullable) //是否可为空
//		if "" != value.ColumnType {
//			ormTagValue += fmt.Sprintf(",%s", value.ColumnType) //列类型
//		}
//		if "" != value.ColumnKey {
//			ormTagValue += fmt.Sprintf(",%s", value.ColumnKey) //ColumnKey
//		}
//		if "" != value.EXTRA {
//			ormTagValue += fmt.Sprintf(",%s", value.EXTRA) //EXTRA
//		}
//		//是否允许为空
//		isNullableValue := fmt.Sprintf("%s", value.IsNullable)
//		//字段描述信息
//		commentValue := fmt.Sprintf("%s", value.ColumnComment)
//
//		//组装后的数据
//		columnStr := fmt.Sprintf("\t%s %s `json:\"%s\" isNullAble:\"%s\" comment:\"%s\"`\n",
//			columnName, dataType, jsonTagValue, isNullableValue, commentValue)
//		//自增字段、datetime字段不需要包含到Insert请求参数中
//		if strings.Contains(ormTagValue, "auto_increment") ||
//			strings.Contains(ormTagValue, "CURRENT_TIMESTAMP") {
//			continue
//		}
//		if "" == insertSqlColumnNames {
//			insertSqlColumnNames += value.ColumnName
//		} else {
//			insertSqlColumnNames += "," + value.ColumnName
//		}
//		if "" == insertSqlColumnPlaces {
//			insertSqlColumnPlaces += "?"
//		} else {
//			insertSqlColumnPlaces += "," + "?"
//		}
//		if "" == insertSqlColumnValues {
//			insertSqlColumnValues += fmt.Sprintf("%s.%s", insertFuncRequestParamName, columnName)
//		} else {
//			insertSqlColumnValues += "," + fmt.Sprintf("%s.%s", insertFuncRequestParamName, columnName)
//		}
//
//		insertDataBuffer.WriteString(columnStr)
//	}
//	insertDataBuffer.WriteString("}\n")
//
//	//结构体赋值
//	insertTemplateData := InsertTemplate{}
//	insertTemplateData.InsertRequestStructureName = insertRequestStructureName
//	insertTemplateData.InsertFuncRequestParamName = insertFuncRequestParamName
//	insertTemplateData.NewItemRequestStructureStrValue = insertDataBuffer.String()
//	insertTemplateData.InsertSqlColumnNames = insertSqlColumnNames
//	insertTemplateData.InsertSqlColumnPlaces = insertSqlColumnPlaces
//	insertTemplateData.InsertSqlColumnValues = insertSqlColumnValues
//
//	return &insertTemplateData, nil
//}
//
////根据数据库表生成数据库基础操作需要的模板数据
//func (orm *JrOrm) GenerateUniqueTemplateData(schemaTable *SchemaTable, tableColumns []SchemaTableColumns) (*DeleteTemplateData, error) {
//	if nil == schemaTable {
//		return nil, errors.New("SchemaTable is nil")
//	}
//	if nil == tableColumns || len(tableColumns) == 0 {
//		return nil, errors.New("SchemaTableColumns is nil")
//	}
//	//结构体赋值
//	templateData := DeleteTemplateData{}
//	//Insert操作数据
//	for _, value := range tableColumns {
//		columnName := FmtColumnsNameToCamelName(value.ColumnName)
//		//数据类型，将数据库类型映射为golang数据类型
//		dataType := MappingMysqlDataTypeToGo(value.DataType, value.ColumnType)
//		//json tag值,将列名称首字母小写后作为tag名称
//		jsonTagValue := fmt.Sprintf("%s", FmtColumnsNameToJsonTagValueName(columnName))
//		//`orm:""`
//		ormTagValue := fmt.Sprintf("%s", value.ColumnName)
//		//if value.ColumnDefault.Valid {
//		//	columnDefault := value.ColumnDefault.String //默认值
//		//	ormTagValue += fmt.Sprintf(",%s", columnDefault)
//		//}
//		ormTagValue += fmt.Sprintf(",%s", value.IsNullable) //是否可为空
//		if "" != value.ColumnType {
//			ormTagValue += fmt.Sprintf(",%s", value.ColumnType) //列类型
//		}
//		if "" != value.ColumnKey {
//			ormTagValue += fmt.Sprintf(",%s", value.ColumnKey) //ColumnKey
//		}
//		if "" != value.EXTRA {
//			ormTagValue += fmt.Sprintf(",%s", value.EXTRA) //EXTRA
//		}
//
//		//判断条件
//		//||
//		//(strings.Contains(ormTagValue, "PRI") && (strings.Contains(ormTagValue, "auto_increment")))
//		if strings.Contains(ormTagValue, "UNI") {
//			data := DeleteByUniQKeyTLData{}
//			data.StructureNameOfTable = FmtColumnsNameToCamelName(schemaTable.TableName)
//			data.DatabaseName = schemaTable.TableSchema
//			data.TableName = schemaTable.TableName
//			data.TableColumnName = value.ColumnName
//			data.ColumnName = columnName
//			data.JsonTagName = jsonTagValue
//			data.DataType = dataType
//			templateData.DeleteByUniQKeyTLData = append(templateData.DeleteByUniQKeyTLData, data)
//		}
//	}
//
//	return &templateData, nil
//}
//
////获取数据库表模板数据,用户根据这些数据并结合自己定义的模板生成文件
//func (orm *JrOrm) GenerateTableTemplateData(schemaTable *SchemaTable) (*DefaultTableTemplate, error) {
//	//参数校验
//	if nil == schemaTable {
//		return nil, errors.New("SchemaTable is nil")
//	}
//	//查询表信息
//	//tableColumns, err := orm.SelectSchemaTableColumns(schemaTable.TableName)
//	//if err != nil {
//	//	return nil, err
//	//}
//	////生成结构体数据
//	//structureStringValue, err := orm.GenerateStructureBySchemaTable(schemaTable, tableColumns)
//	//if err != nil {
//	//	return nil, err
//	//}
//	////生成数据库基础操作需要的模板数据
//	//insertTemplateData, err := orm.GenerateInsertTemplateData(schemaTable, tableColumns)
//	//if err != nil {
//	//	return nil, err
//	//}
//	////生成数据库基础操作需要的模板数据
//	//deleteTemplateData, err := orm.GenerateUniqueTemplateData(schemaTable, tableColumns)
//	//if err != nil {
//	//	return nil, err
//	//}
//
//	//组装返回数据
//	templateData := DefaultTableTemplate{}
//	templateData.DatabaseName = schemaTable.TableSchema
//	templateData.TableName = schemaTable.TableName
//	//templateData.TableComment = schemaTable.TableComment
//	templateData.PackageName = strings.ToLower(FmtColumnsNameToCamelName(schemaTable.TableSchema))
//	templateData.StructureName = FmtColumnsNameToCamelName(schemaTable.TableName)
//	//templateData.StructureStringValue = structureStringValue
//	//templateData.InsertTemplate = *insertTemplateData
//	//templateData.DeleteTemplateData = *deleteTemplateData
//	//templateData.SelectTemplateData = *deleteTemplateData
//
//	return &templateData, nil
//}
