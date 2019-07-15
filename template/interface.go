package template
import (
	"errors"
	"fmt"
	"github.com/redresseur/swagger/analyse"
	"github.com/redresseur/utils/charset"
)

const (
	StringType  = `string`
)

/*
	introduce: generate interface
	date: 2019/07/15
	author: wangzhipengtest@163.com
*/

type Param map[string]string

type method struct {
	Name string
	Parameters []Param
	Returns []Param
}

var (
	globalMethods = map[string][]*method{}
)

var (
	ErrTagsNotExist = errors.New("the list of tags was empty.")
)

func apiMethod (api *analyse.RestApi)(*method, error) {
	m := &method{}
	for _, param := range  api.Parameters{
		p := param.(*analyse.Parameter)
		// 檢查Schema 是否為空
		tp := Param{}
		if p.Schema != nil{
			s := p.Schema.(*analyse.Schema)
			// 檢查引用是否爲空
			if s.Reference != ""{
				if def, ok :=  globalDefs[s.Reference]; ok{
					return nil, fmt.Errorf("the reference %s is not valide", s.Reference)
				}else {
					tp[p.Name] = ptrto(def.Name)
				}
			}else {
				// TODO: 支持除了引用之外的其他类型，例如： object 等
				return 	nil, fmt.Errorf("the schema is not support currently exclude reference.")
			}

			m.Parameters = append(m.Parameters, tp)
		}else {
			switch (p.Type){
			case StringType:
				{
					tp[p.Name] = StringType
				}
			default:
				return nil, fmt.Errorf("the param type %s is not supported", p.Type)
			}
		}
	}
	return nil, nil
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
		interfaceName, err := charset.HumpFormat(api.Tags[0])
		if err != nil{
			return err
		}

		globalMethods[interfaceName]= append(globalMethods[interfaceName], )
	}
	return nil	
}