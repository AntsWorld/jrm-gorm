package jdb

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
)

func init() {
	//日志输出样式
	log.SetFlags( /*log.Lshortfile |*/ log.Ldate | log.Ltime)
}

type SqlExecutor struct {
	DB                  *sql.DB
	StructureOrmTagName string
}

//执行非查询SQL语句(增、删、改)
func (sqlExecutor *SqlExecutor) ExecuteNonSelectSql(sqlStatement string, args ...interface{}) (sql.Result, error) {
	if sqlExecutor.DB == nil {
		return nil, errors.New("sqlExecutor.DB can't be nil")
	}
	if sqlStatement == "" {
		return nil, errors.New("sql statement can't be empty")
	}
	stmt, err := sqlExecutor.DB.Prepare(sqlStatement)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := stmt.Close(); err != nil {
			log.Println(err)
			return
		}
	}()
	result, err := stmt.Exec(args...)
	if err != nil {
		return nil, err
	}
	return result, nil
}

//执行查询SQL语句
func (sqlExecutor *SqlExecutor) ExecuteSelectSql(sqlStatement string, v interface{}, args ...interface{}) error {
	if sqlExecutor.StructureOrmTagName == "" {
		return errors.New("tag name can't be empty")
	}
	data, err := sqlExecutor.ExecuteSelectSqlForMapResult(sqlStatement, args...)
	if err != nil {
		log.Println(err)
		return err
	}
	//log.Println(data)
	if err := ReflectSelectSqlMapResultToStructure(data, v, sqlExecutor.StructureOrmTagName); err != nil {
		return err
	}
	return nil
}

//执行查询语句并将结果转换为Map切片
func (sqlExecutor *SqlExecutor) ExecuteSelectSqlForMapResult(sqlStatement string, args ...interface{}) ([]map[string]string, error) {
	rows, err := sqlExecutor.ExecuteSelectSqlForRows(sqlStatement, args...)
	if err != nil {
		return nil, err
	}
	//log.Println("rows = ", rows)
	defer func() {
		if err := rows.Close(); err != nil {
			log.Println(err)
			return
		}
	}()
	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	// Make a slice for the values
	values := make([]sql.RawBytes, len(columns))
	selectMapResult := make([]map[string]string, 0)
	// rows.Scan wants '[]interface{}' as an argument, so we must copy the references into such a slice
	// See http://code.google.com/p/go-wiki/wiki/InterfaceSlice for details
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	// Fetch rows
	for rows.Next() {
		// get RawBytes from data
		if err := rows.Scan(scanArgs...); err != nil {
			return nil, err
		}
		rowRet := make(map[string]string)
		for i, rawByte := range values {
			// Here we can check if the value is nil (NULL value)
			if rawByte == nil {
				rowRet[columns[i]] = ""
			} else {
				rowRet[columns[i]] = string(rawByte)
			}
			//fmt.Println(columns[i], ": ", value)
		}
		selectMapResult = append(selectMapResult, rowRet)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return selectMapResult, nil
}

//执行查询语句并返回查询结果
func (sqlExecutor *SqlExecutor) ExecuteSelectSqlForRows(sqlStatement string, args ...interface{}) (*sql.Rows, error) {
	if sqlExecutor.DB == nil {
		return nil, errors.New("sqlExecutor.DB can't be nil")
	}
	if sqlStatement == "" {
		return nil, errors.New("sql statement can't be empty")
	}
	stmt, err := sqlExecutor.DB.Prepare(sqlStatement)
	if err != nil {
		log.Println("Prepare err:", err)
		return nil, err
	}
	//log.Println("Prepare ok")
	defer func() {
		if err := stmt.Close(); err != nil {
			log.Println(err)
			return
		}
	}()
	rows, err := stmt.Query(args...)
	if err != nil {
		log.Println("Query err:", err)
		return nil, err
	}
	//log.Println("Query ok")
	return rows, nil

}

//将查询语句得到的[]Map结果反射到结构体中
func ReflectSelectSqlMapResultToStructure(data []map[string]string, v interface{}, tagName string) error {
	//判断数据是否为空,如果数据为空表示未查询到数据，不返回错误
	if nil == data {
		log.Println(errors.New("data can't be nil"))
		return nil
	}
	if len(data) == 0 {
		log.Println(errors.New("data can't be empty"))
		return nil
	}
	if v == nil {
		return errors.New("v can't be nil")
	}
	dataLen := len(data)
	vType := reflect.TypeOf(v)
	if vType.Kind() == reflect.Ptr {
		if vType.Elem().Kind() == reflect.Struct {
			//传入的是结构体地址,如果数据长度大于1,提示需要传入Slice才能接收所有数据
			if len(data) > 1 {
				return errors.New("data length more than one,type of v must be slice")
			}
			if err := ReflectMapToStructure(data[0], reflect.ValueOf(v), tagName); err != nil {
				return err
			}
		} else if vType.Elem().Kind() == reflect.Slice && vType.Elem().Elem().Kind() == reflect.Struct {
			srcPtrValue := reflect.ValueOf(v).Elem()
			srcPtrValueType := srcPtrValue.Type()
			newValue := reflect.MakeSlice(srcPtrValueType, 0, dataLen)
			srcPtrValue.Set(newValue)
			srcPtrValue.SetLen(dataLen)
			//log.Println("srcPtrValueType = ", srcPtrValueType.String())
			for i := 0; i < dataLen; i++ {
				newObj := reflect.New(srcPtrValueType.Elem())
				if err := ReflectMapToStructure(data[i], newObj, tagName); err != nil {
					return err
				}
				srcPtrValue.Index(i).Set(newObj.Elem())
			}
		} else {
			return errors.New("unSupport type")
		}
	} else {
		return errors.New("type must be Ptr")
	}

	return nil
}

//将map值映射到结构体
func ReflectMapToStructure(m map[string]string, v reflect.Value, tagName string) error {
	vType := v.Type()
	vElem := v.Elem()
	vTypeElem := vType.Elem()
	if !vElem.IsValid() {
		return errors.New("unValid value elem")
	}
	for i := 0; i < vElem.NumField(); i++ {
		value := vElem.Field(i)
		kind := value.Kind()
		tagValue := vTypeElem.Field(i).Tag.Get(tagName)
		if tagValue == "" {
			return errors.New("tag value can't be empty")
		}
		//对tagValue进行处理，获取第一个值
		tagValues := strings.Split(tagValue, ",")
		key := tagValues[0]
		if key != "" {
			meta, ok := m[key]
			if !ok {
				log.Println(fmt.Sprintf("map[%s] is not ok", key))
				continue
			}
			if !value.CanSet() {
				return errors.New("value can't set")
			}
			if meta == "" {
				continue
			}
			switch kind {
			case reflect.Bool:
				if result, err := strconv.ParseBool(meta); err != nil {
					return err
				} else {
					value.SetBool(result)
				}
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				if result, err := strconv.ParseInt(meta, 10, 64); err != nil {
					return err
				} else {
					value.SetInt(result)
				}
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				if result, err := strconv.ParseUint(meta, 10, 64); err != nil {
					return err
				} else {
					value.SetUint(result)
				}
			case reflect.Float32, reflect.Float64:
				if result, err := strconv.ParseFloat(meta, 64); err != nil {
					return err
				} else {
					value.SetFloat(result)
				}
			case reflect.String:
				value.SetString(meta)
			//case reflect.Uintptr:
			//case reflect.Complex64:
			//case reflect.Complex128:
			//case reflect.Array:
			//case reflect.Chan:
			//case reflect.Func:
			//case reflect.Interface:
			//case reflect.Map:
			//case reflect.Ptr:
			//case reflect.Slice:
			//case reflect.Struct:
			//	log.Println("value.Type() = ", value.Type(), ",value = ", meta)
			//	structValue := reflect.New(value.Type())
			//
			//	value.Set(structValue.Elem())
			default:
				return errors.New(fmt.Sprintf("type %s is not yet supported", kind.String()))
			}
		} else {

		}
	}
	return nil
}
