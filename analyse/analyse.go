package analyse

/*
	introduce: 用来解析swagger.yaml
	author: wangzhipengtest@163.com
*/

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"gopkg.in/yaml.v2"
	"os"
	"reflect"
)

const(
	PATHS  = `paths`
	SWAGGER = `swagger`
	INFO = `info`
	SCHEMES = `schemes`
	HOST = `host`
	BASEPATH = `basePath`
	TAGS = `tags`
	DEFINITIONS = `definitions`
)

// 读取swagger.json
func ReadJson(path string)(map[string]interface{}, error) {
	fd, err := os.Open(path)
	if err != nil{
		return nil, err
	}
	defer fd.Close()

	swaggerSturct := map[string]interface{}{}
	decoder := json.NewDecoder(fd)
	decoder.Decode(&swaggerSturct)

	return swaggerSturct, err
}

// 读取swagger.yaml
func ReadYaml(path string)(map[string]interface{}, error) {
	fd, err := os.Open(path)
	if err != nil{
		return nil, err
	}

	yamlDecoder := yaml.NewDecoder(fd)
	swaggerSturct := map[string]interface{}{}
	yamlDecoder.Decode(&swaggerSturct)

	return swaggerSturct, err
}

const (
	GET  = `get`
	POST = `post`
)

type Parameter struct {
	In string `json:"in" yaml:"in" mapstructure:"in"`
	Name string `json:"name" yaml:"name" mapstructure:"name"`
	Type string `json:"type" yaml:"type" mapstructure:"type"`
	Schema interface{} `json:"schema" yaml:"schema" mapstructure:"schema"`
}

type RestApiDef struct {
	Tags []string `yaml:"tags" mapstructure:"tags" json:"tags"`
	Consumes []string `json:"consumes" yaml:"consumes" mapstructure:"consumes"`
	Produces []string `json:"produces" yaml:"produces" mapstructure:"produces"`
	Parameters []interface{} `json:"parameters" yaml:"parameters" mapstructure:"parameters"`
	Responses interface{} `json:"responses" yaml:"responses" mapstructure:"responses"`
} 

type RestApi struct {
	Url string
	Method string
	RestApiDef
}

type Items struct {
	Type string `json:"type" yaml:"type" mapstructure:"type"`
	Reference string `json:"$ref" yaml:"$ref" mapstructure:"$ref"`
}

type Field struct {
	Type string `json:"type" yaml:"type" mapstructure:"type"`
	Format string `json:"format" yaml:"format" mapstructure:"format"`
	Items interface{} `json:"items" yaml:"items" mapstructure:"items"`
	Reference string `json:"$ref" yaml:"$ref" mapstructure:"$ref"`
	Properties map[string]interface{} `json:"properties" yaml:"properties" mapstructure:"properties"`
}

type DefinitionDef struct {
	Type string `json:"type" yaml:"type" mapstructure:"type"`

	// 对应Field 结构体
	Properties map[string]interface{} `json:"properties" yaml:"properties" mapstructure:"properties"`
}

type Definition struct {
	Name string
	DefinitionDef
}

type Schema struct {
	// DefinitionDef
	Field
}

func GetRestApi(swaggerMap map[string]interface{})(apis []*RestApi, err error)  {
	paths, ok := swaggerMap[PATHS]
	if ! ok{
		return nil, errors.New("the paths are not existed")
	}

	for url, desc := range paths.(map[interface{}]interface{}){
		api := &RestApi{}
		descMap := desc.(map[interface{}]interface{})
		fmt.Println(reflect.TypeOf(url).Kind().String())
		if urlstr, ok :=  url.(string); ok{
			//fmt.Printf("%s\n", urlstr)
			api.Url = urlstr
		}

		for key, value := range descMap{
			//fmt.Printf("%s\n", key)
			if keystr, ok := key.(string); ok{
				api.Method = keystr
			}

			if err = mapstructure.Decode(value, &api.RestApiDef); err != nil{
				fmt.Printf("%v", err)
				return nil, err
			}

			for i, param := range api.RestApiDef.Parameters{
				p := &Parameter{}
				if err = mapstructure.Decode(param, p);err != nil{
					fmt.Printf("%v", err)
					return nil, err
				}else {
					if p.Schema != nil{
						s := &Schema{}
						if err = mapstructure.Decode(p.Schema.(map[interface{}]interface{}), &s.Field); err != nil{
							return nil, err
						}

						//if err = mapstructure.Decode(p.Schema.(map[interface{}]interface{}), &s.DefinitionDef); err != nil{
						//	return nil, err
						//}
						p.Schema = s
					}
				}

				api.RestApiDef.Parameters[i] = p
			}
		}

		apis = append(apis, api)
	}

	return
}

// 取出swagger中的结构体定义，
// 需要注意的是definitions有不存在的可能性
func GetDefinition(swaggerMap map[string]interface{})(defs []*Definition, err error){
	definitions, ok := swaggerMap[DEFINITIONS]
	if !ok{
		return nil, errors.New("the definitions are not exist")
	}

	for name, desc := range definitions.(map[interface{}]interface{}){
		definition := &Definition{}
		// fmt.Println(reflect.TypeOf(name).Kind().String())
		if namestr, ok :=  name.(string); ok{
			//fmt.Printf("%s\n", namestr)
			definition.Name = namestr
		}

		if err = mapstructure.Decode(desc, & definition.DefinitionDef); err != nil{
			fmt.Printf("%v", err)
			return
		}

		properties(definition.Properties)
		defs = append(defs, definition)
	}

	return
}

// 嵌套解析Properties
func properties(properties map[string]interface{} ) (err error) {
	for name, data := range properties{
		properties[name], err = filed(data)
		if err != nil{
			return
		}
	}

	return
}

func filed(filedData interface{}) (*Field, error) {
	f := &Field{}
	err := mapstructure.Decode(filedData, f)
	if err != nil{
		return nil, err
	}

	if f.Items != nil{
		items := &Items{}
		if err := mapstructure.Decode(f.Items, items); err != nil{
			return nil, err
		}

		f.Items = items
	}

	if 0 != len(f.Properties){
		err = properties(f.Properties)
	}

	return f, err
}

