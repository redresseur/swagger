package template

/*
	introduction: 生成结构体代码
	author: wangzhipengtest@163.com
*/

import (
	"errors"
	"fmt"
	"github.com/redresseur/swagger/analyse"
	"github.com/redresseur/utils/charset"
	"os"
	//"path/filepath"
	tt "text/template"
)

const(
	DefinitionPrefix = "#/definitions/"
	Reference = `$ref`

	Point = `*`
	Array = `[]`
	Interface = `interface{}`
)

const (
	ObjectType = `object`
	ArrayType = `array`
)

var globalStructures =  map[string] templateStructure{}
var globalDefs = map[string]*analyse.Definition{}

// *structName
func ptrto(name string) string {
	return Point + name
}

func arr(name string) string  {
	return Array + name
}

// 字段名称 类型 别名
type templateStructure map[string]string

// TODO: 格式化为驼峰命名
func inlineStructureName( structName, filedName string )(string, error)  {
	return charset.HumpFormat(structName, filedName)
}

func filed( s templateStructure, fieldName string,
	structName string, field *analyse.Field) error {
	// 引用不为空说明引用了其他结构体
	if field.Reference != ""{
		if def ,ok := globalDefs[field.Reference]; !ok {
			return fmt.Errorf("the reference %s is not valide", field.Reference)
		}else {
			s[fieldName] = ptrto(def.Name)
		}
	}else {
		// 引用为空说明是一个定义好的字段
		switch field.Type {
		case ObjectType: // 说明这是内置的默认结构体
			{
				if 0 != len(field.Properties){
					// 指向内置结构体的指针
					if inStName, err := charset.HumpFormat(structName, fieldName); err != nil{
						return fmt.Errorf("make strue name by [%s , %s]  : %v", structName, fieldName, err)
					}else {
						if err := structure(inStName, field.Properties); err != nil{
							return err
						}
						s[fieldName] = ptrto(inStName)
					}
				}else {
					s[fieldName] = Interface
				}
			}
		case ArrayType:
			{
				if field.Items == nil{
					return errors.New("the items is empty")
				}

				items := field.Items.(*analyse.Items)
				if items.Reference != ""{
					if def, ok := globalDefs[items.Reference]; !ok{
						return fmt.Errorf("the reference of %s is not valid in %s", items.Reference, structName)
					}else {
						s[fieldName] = arr(ptrto(def.Name))
					}
				}else if items.Type != "" {
					s[fieldName] = arr(items.Type)
				}else {
					return fmt.Errorf("the items of %s is not valid", structName)
				}
			}
		default:
			s[fieldName] = field.Type
		}
	}

	return nil
}

func structure(structName string, Properties map[string]interface{}) error {
	s := templateStructure{}
	for filedName, filedData := range Properties{
		f := filedData.(*analyse.Field)
		if err := filed(s, filedName, structName, f); err != nil{
			return err
		}
	}

	globalStructures[structName] = s
	return nil
}

// 拼装结构体定义 用于生成结构体代码
func definitionComplete(defs []*analyse.Definition) error {
	// 緩存所有的def, 用于做索引補全
	for _, def := range defs{
		path := DefinitionPrefix + def.Name
		globalDefs[path] = def
	}

	// 找到用于代替索引的部分
	for _, def := range defs{
		if err := structure(def.Name, def.Properties); err != nil{
			return err
		}
	}

	return nil
}

const template  = `{{$input:=.}}{{range $structureName, $structure := $input}}
type {{$structureName}}{{$filedNum := $structure|len}} struct{
    {{range $fieldName, $fieldType := $structure}}{{$fieldName|fieldNameFormat}} {{$fieldType}} `+ "`" + `json:"{{$fieldName}}"` + "`\n" +
`{{with $filedNum =  (sub $filedNum 1) }}{{if (gt $filedNum 0)}}    {{end}}{{end}}{{end}}}
{{end}}`

// 駝峰命名
func fieldNameFormat(name string) string {
	if charset.CheckSpecialCharacter(name){
		panic(fmt.Sprintf("filed name %s has included special charset.", name))
	}

	humpName, err := charset.HumpFormat(name)
	if err != nil{
		panic(err)
	}
	return humpName
}

func sub(a int, b int) int {
	return a - b
}

func outputStructureCode()error{
	funcMap := tt.FuncMap{
		"fieldNameFormat": fieldNameFormat,
		"sub": sub,
	}
	t, err := tt.New("structure").Funcs(funcMap).Parse(string(template))
	if  err != nil{
		return  err
	}

	return t.Execute(os.Stdout, globalStructures)
}