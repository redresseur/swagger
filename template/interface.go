package template
import (
	"errors"
	"fmt"
	"github.com/redresseur/swagger/analyse"
	"github.com/redresseur/utils/charset"
	"io/ioutil"
	"os"
	tt "text/template"
)

const (
	StringType  = `string`
	IntType = `integer`
)

/*
	introduce: generate interface
	date: 2019/07/15
	author: wangzhipengtest@163.com
*/

type Param map[string]string

type method struct {
	Name string
	Parameters Param
	Returns Param
}

var (
	globalMethods = map[string][]*method{}
)

var (
	ErrTagsNotExist = errors.New("the list of tags was empty.")
)

func apiMethod (api *analyse.RestApi)(m *method, err error) {
	m = &method{Parameters:Param{}}
	if m.Name, err = charset.CamelCaseFormat(true, api.OperationId); err != nil{
		return nil, err
	}

	for _, param := range  api.Parameters{
		p := param.(*analyse.Parameter)
		// 檢查Schema 是否為空
		if p.Schema != nil{
			s := p.Schema.(*analyse.Schema)
			// 檢查引用是否爲空
			if s.Reference != ""{
				if def, ok :=  globalDefs[s.Reference]; !ok{
					return nil, fmt.Errorf("the reference %s is not valide", s.Reference)
				}else {
					m.Parameters[p.Name] = ptrto(def.Name)
				}
			}else {
				// TODO: 支持除了引用之外的其他类型，例如： object 等
				return 	nil, fmt.Errorf("the schema is not support currently exclude reference.")
			}
		}else {
			// 目前只支持了，StringType 和 IntType
			switch (p.Type){
			case StringType:
				{
					m.Parameters[p.Name] = StringType
				}
			case IntType:
				{
					m.Parameters[p.Name] = IntType
				}
			default:
				return nil, fmt.Errorf("the param type %s is not supported", p.Type)
			}
		}
		//m.Parameters = append(m.Parameters, tp)
	}

	// 提取responses
	responses, ok := api.Responses.(*analyse.Responses);
	if  !ok{
		return m, nil
	}

	m.Returns = Param{}
	for statusCode, rspDef := range responses.RespDefinitions{
		s, ok := rspDef.Schema.(*analyse.Schema);
		if  !ok{
			err = fmt.Errorf("the response definion is not valid in %s%s", m.Name, statusCode)
			return
		}

		if s.Reference != ""{
			if def, ok :=  globalDefs[s.Reference]; !ok{
				return nil, fmt.Errorf("the reference %s is not valide", s.Reference)
			}else {
				rspName, _ := charset.CamelCaseFormat(false, "rsp", statusCode)
				m.Returns[rspName] = ptrto(def.Name)
				continue
			}
		}

		switch (s.Type) {
		case ObjectType:
			if len(s.Properties) != 0 {
				if rspObjectName, err := charset.CamelCaseFormat(true, m.Name, statusCode, "Rsp"); err != nil{
					return nil, err
				}else {
					if err = structure(rspObjectName ,s.Properties); err != nil{
						return nil, err
					}else {
						rspName, _ := charset.CamelCaseFormat(false,"rsp", statusCode)
						m.Returns[rspName] = ptrto(rspObjectName)
					}
				}
			}else {
				rspName, _ := charset.CamelCaseFormat(false,"rsp", statusCode)
				m.Returns[rspName] = Interface
			}
		case IntType:
			m.Returns[statusCode] = s.Format
		case ArrayType:
			if s.Items == nil{
				err = errors.New("the items' definition are empty")
				return
			}

			items := s.Items.(*analyse.Items)
			if items.Reference != ""{
				if def, ok := globalDefs[items.Reference]; !ok{
					err = fmt.Errorf("the reference of %s is not valid", items.Reference)
					return
				}else {
					rspName, _ := charset.CamelCaseFormat(false,"rsp", statusCode)
					m.Returns[rspName] = arr(ptrto(def.Name))
				}
			}else if items.Type != "" {
				if items.Type == ObjectType {
					rspName, _ := charset.CamelCaseFormat(false,"rsp", statusCode)
					m.Returns[rspName] = arr(Interface)
				} else {
					rspName, _ := charset.CamelCaseFormat(false,"rsp", statusCode)
					m.Returns[rspName] = arr(items.Type)
				}
			}else {
				err = fmt.Errorf("the items is not valid" )
				return
			}
		default:
			rspName, _ := charset.CamelCaseFormat(false,"rsp", statusCode)
			m.Returns[rspName] = s.Type
		}
	}

	return m, nil
}

// api 的定义中必须要有tags
// 子tags要在母tags之上，
// 比如:
// 		- child
//		- parent
// 最终child 接口会被parent接口所包含
func interfaceComplete(restFulApis []*analyse.RestApi) error {
	for _, api := range restFulApis{
		if len(api.Tags) == 0{
			return ErrTagsNotExist
		}

		// 用 tag[0] 作为interface 的名称
		interfaceName, err := charset.CamelCaseFormat(true, api.Tags[0])
		if err != nil{
			return err
		}

		if m, err := apiMethod(api); err != nil{
			return err
		} else {
			globalMethods[interfaceName]= append(globalMethods[interfaceName], m)
		}
	}
	return nil	
}

func outputInterfaceCode()error  {
	if fd, err := os.Open("interface.tpl"); err != nil{
		return err
	}else {
		if data, err := ioutil.ReadAll(fd); err != nil{
			return err
		}else {
			funcMap := tt.FuncMap{
				"fieldNameFormat": fieldNameFormat,
				"sub": sub,
			}

			if t, err := tt.New("interface").Funcs(funcMap).Parse(string(data));err != nil{
				return err
			}else {
				if err := t.Execute(os.Stdout, globalMethods); err != nil{
					return err
				}
			}
		}
	}

	return nil
}