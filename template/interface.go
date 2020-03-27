package template

/*
	introduce: generate interface
	date: 2019/07/15
	author: wangzhipengtest@163.com
*/

import (
	"errors"
	"fmt"
	"github.com/redresseur/swagger/analyse"
	"github.com/redresseur/utils/charset"
	"io"
	"io/ioutil"
	"os"
	tt "text/template"
)

const (
	StringType  = `string`
	IntType     = `integer`
	BooleanType = `boolean`
)

type Param struct {
	IN          string
	SwaggerType string `desc:"在swagger中定义的类型"`
	Type        string `desc:"实际代码中定义的类型"`
	Key         string
	Name        string
}

type Result struct {
	StatusCode  string `desc:"返回對應的狀態碼"`
	Type        string `desc:"實際代碼中的類型"`
	SwaggerType string `desc:"在swagger中定義的類型"`
	// Key string
	Name string
}

//type Param map[string]string

type method struct {
	Name        string
	Parameters  []*Param
	Returns     []*Result
	MethodType  string `desc:"請求的類型：GET POST OPTIONS DELETE PUT TRACE HEADER CONNECT"`
	Url         string
	OperationId string
	Tags        []string
}

var (
	globalMethods = map[string][]*method{}
)

var (
	ErrTagsNotExist = errors.New("the list of tags was empty.")
)

func apiMethod(url, meth string, def *analyse.RestApiDef) (m *method, err error) {
	m = &method{}
	m.MethodType = meth
	m.Url = url
	m.OperationId = def.OperationId
	m.Tags = def.Tags

	if m.Name, err = charset.CamelCaseFormat(true, def.OperationId); err != nil {
		return nil, err
	}

	for _, param := range def.Parameters {
		p := param.(*analyse.Parameter)
		// 檢查Schema 是否為空
		if p.Schema != nil {
			s := p.Schema.(*analyse.Schema)
			// 檢查引用是否爲空
			if s.Reference != "" {
				if def, ok := globalDefs[s.Reference]; !ok {
					return nil, fmt.Errorf("the reference %s is not valide", s.Reference)
				} else {
					m.Parameters = append(m.Parameters, &Param{
						SwaggerType: ObjectType,
						Type:        ptrto(def.Name),
						IN:          p.In,
						Key:         p.Name,
						Name:        charset.CamelCaseFormatMust(false, p.Name),
					})
				}
			} else {
				// TODO: 支持除了引用之外的其他类型，例如： object 等
				return nil, fmt.Errorf("the schema is not support currently exclude reference.")
			}
		} else {
			// 目前只支持了，StringType 和 IntType
			switch p.Type {
			case StringType:
				{
					m.Parameters = append(m.Parameters, &Param{
						SwaggerType: StringType,
						Type:        StringType,
						IN:          p.In,
						Key:         p.Name,
						Name:        charset.CamelCaseFormatMust(false, p.Name),
					})
				}
			case IntType:
				{
					m.Parameters = append(m.Parameters, &Param{
						SwaggerType: IntType,
						Type:        p.Format,
						IN:          p.In,
						Key:         p.Name,
						Name:        charset.CamelCaseFormatMust(false, p.Name),
					})
				}
			case BooleanType:
				{
					m.Parameters = append(m.Parameters, &Param{
						SwaggerType: BooleanType,
						Type:        `bool`,
						IN:          p.In,
						Key:         p.Name,
						Name:        charset.CamelCaseFormatMust(false, p.Name),
					})
				}
			default:
				return nil, fmt.Errorf("the param type %s is not supported", p.Type)
			}
		}
		//m.Parameters = append(m.Parameters, tp)
	}

	// 提取responses
	responses, ok := def.Responses.(*analyse.Responses)
	if !ok {
		//return m, nil
		return nil, fmt.Errorf("the repsonses of %s is empty in [%s] case", m.Url, meth)
	}

	for statusCode, rspDef := range responses.RespDefinitions {
		s, ok := rspDef.Schema.(*analyse.Schema)
		if !ok {
			err = fmt.Errorf("the response definion is not valid in %s%s", m.Name, statusCode)
			return
		}

		if s.Reference != "" {
			if def, ok := globalDefs[s.Reference]; !ok {
				return nil, fmt.Errorf("the reference %s is not valide", s.Reference)
			} else {
				rspName, _ := charset.CamelCaseFormat(false, "rsp", statusCode)
				m.Returns = append(m.Returns, &Result{
					Type:        ptrto(def.Name),
					SwaggerType: ObjectType,
					StatusCode:  statusCode,
					Name:        rspName,
				})
				continue
			}
		}

		switch s.Type {
		case ObjectType:
			if len(s.Properties) != 0 {
				if rspObjectName, err := charset.CamelCaseFormat(true, m.Name, statusCode, "Rsp"); err != nil {
					return nil, err
				} else {
					if err = structure(rspObjectName, s.Properties); err != nil {
						return nil, err
					} else {
						rspName, _ := charset.CamelCaseFormat(false, "rsp", statusCode)
						m.Returns = append(m.Returns, &Result{
							Type:        ptrto(rspObjectName),
							SwaggerType: ObjectType,
							StatusCode:  statusCode,
							Name:        rspName,
						})
					}
				}
			} else {
				rspName, _ := charset.CamelCaseFormat(false, "rsp", statusCode)
				m.Returns = append(m.Returns, &Result{
					Type:        Interface,
					SwaggerType: ObjectType,
					StatusCode:  statusCode,
					Name:        rspName,
				})
			}
		case IntType:
			rspName, _ := charset.CamelCaseFormat(false, "rsp", statusCode)
			m.Returns = append(m.Returns, &Result{
				Type:        s.Format,
				SwaggerType: IntType,
				StatusCode:  statusCode,
				Name:        rspName,
			})
		case BooleanType:
			rspName, _ := charset.CamelCaseFormat(false, "rsp", statusCode)
			m.Returns = append(m.Returns, &Result{
				Type:        `bool`,
				SwaggerType: BooleanType,
				StatusCode:  statusCode,
				Name:        rspName,
			})
		case ArrayType:
			if s.Items == nil {
				err = errors.New("the items' definition are empty")
				return
			}

			items := s.Items.(*analyse.Items)
			if items.Reference != "" {
				if def, ok := globalDefs[items.Reference]; !ok {
					err = fmt.Errorf("the reference of %s is not valid", items.Reference)
					return
				} else {
					rspName, _ := charset.CamelCaseFormat(false, "rsp", statusCode)
					m.Returns = append(m.Returns, &Result{
						Type:        arr(ptrto(def.Name)),
						SwaggerType: ArrayType,
						StatusCode:  statusCode,
						Name:        rspName,
					})
				}
			} else if items.Type != "" {
				if items.Type == ObjectType {
					rspName, _ := charset.CamelCaseFormat(false, "rsp", statusCode)
					m.Returns = append(m.Returns, &Result{
						Type:        arr(Interface),
						SwaggerType: ArrayType,
						StatusCode:  statusCode,
						Name:        rspName,
					})
				} else {
					rspName, _ := charset.CamelCaseFormat(false, "rsp", statusCode)
					m.Returns = append(m.Returns, &Result{
						Type:        arr(items.Type),
						SwaggerType: ArrayType,
						StatusCode:  statusCode,
						Name:        rspName,
					})
				}
			} else {
				err = fmt.Errorf("the items is not valid")
				return
			}
		default:
			rspName, _ := charset.CamelCaseFormat(false, "rsp", statusCode)
			m.Returns = append(m.Returns, &Result{
				Type:        s.Type,
				SwaggerType: s.Type,
				StatusCode:  statusCode,
				Name:        rspName,
			})
		}
	}

	return
}

// api 的定义中必须要有tags
// 子tags要在母tags之上，
// 比如:
// 		- child
//		- parent
// 最终child 接口会被parent接口所包含
func InterfaceComplete(restFulApis []*analyse.RestApi) error {
	for _, api := range restFulApis {
		for meth, def := range api.RestApiDefs {
			if len(def.Tags) == 0 {
				return ErrTagsNotExist
			}

			// 用 tag[0] 作为interface 的名称
			//interfaceName := charset.CamelCaseFormatMust(true, def.Tags[0])
			// FixMe: 在模板中轉化
			interfaceName := def.Tags[0]
			if m, err := apiMethod(api.Url, meth, def); err != nil {
				return err
			} else {
				globalMethods[interfaceName] = append(globalMethods[interfaceName], m)
			}
		}
	}
	return nil
}

func OutputInterfaceCode(writer io.Writer) error {
	funcMap := tt.FuncMap{
		"fieldNameFormat": fieldNameFormat,
		"sub":             sub,
	}

	if t, err := tt.New("interface").Funcs(funcMap).Parse(string(interfaceTemplate)); err != nil {
		return err
	} else {
		if err := t.Execute(writer, globalMethods); err != nil {
			return err
		}
	}

	return nil
}

func OutputInterfaceCodeWithTemplate(writer io.Writer, Path string) error {
	if fd, err := os.Open(Path); err != nil {
		return err
	} else {
		if data, err := ioutil.ReadAll(fd); err != nil {
			return err
		} else {
			funcMap := tt.FuncMap{
				"fieldNameFormat": fieldNameFormat,
				"sub":             sub,
			}

			if t, err := tt.New("interface").Funcs(funcMap).Parse(string(data)); err != nil {
				return err
			} else {
				if err := t.Execute(writer, globalMethods); err != nil {
					return err
				}
			}
		}
	}

	return nil
}
