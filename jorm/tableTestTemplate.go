package jorm

//######################################################################################################################
//数据库操作函数测试模板
//######################################################################################################################
var DefaultTableTestTemplateText = `package {{.BaseTableTemplate.PackageName}}

import "testing"

//Insert操作测试
{{range $k,$v := .InsertTemplates}}
	func Test{{$v.FunctionName}}(t *testing.T) {
		
	}
{{end}}
//测试Delete操作
{{range $k,$v := .DeleteTemplates}}
	func Test{{$v.FunctionName}}(t *testing.T) {
		
	}
{{end}}
//测试Query操作
{{range $k,$v := .QueryTemplates}}
	func Test{{$v.FunctionName}}(t *testing.T) {
		
	}
{{end}}
//测试Update操作
{{range $k,$v := .UpdateTemplates}}
	func Test{{$v.FunctionName}}(t *testing.T) {
		
	}
{{end}}
`
