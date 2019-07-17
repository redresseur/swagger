package template

import (
	"fmt"
	"github.com/redresseur/swagger/analyse"
	"github.com/redresseur/utils/charset"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	tt "text/template"
)

type Router struct {
	Method *method
	Instance string
}

var globalRouters = []*Router{}

const (
	Body  = `body`
	Path = `path`
	Query = `query`
	Default = `default`
)

func router( def *analyse.RestApiDef) (*Router, error) {
	res := &Router{}

	interfaceName, err := charset.CamelCaseFormat(true, def.Tags[0])
	if err != nil{
		return nil, err
	}

	res.Instance = interfaceName
	methods, ok := globalMethods[interfaceName]
	if  !ok {
		return nil, fmt.Errorf("the interface %s is valid", interfaceName)
	}

	methodName, err := charset.CamelCaseFormat(true, def.OperationId)
	if err != nil{
		return nil, err
	}

	for _, m := range methods{
		if m.Name == methodName {
			res.Method = m
			return res, nil
		}
	}

	return nil, fmt.Errorf("method %s not found in %s", methodName, interfaceName)
}

/*
	introduce: generate router
	date: 2019/07/15
	author: wangzhipengtest@163.com
*/

func RouterComplete(restFulApis []*analyse.RestApi) error {
	for _, api := range restFulApis{
		for _, def := range api.RestApiDefs{
			if len(def.Tags) == 0{
				return ErrTagsNotExist
			}

			if r, err := router(def); err != nil{
				return err
			}else {
				globalRouters = append(globalRouters, r)
			}
		}
	}
	return nil
}

func excludePtr(ptr string) string {
	return strings.NewReplacer("*", "").Replace(ptr)
}

func atoi(code string )int  {
	if code == Default {
		return -1
	}

	res, err := strconv.ParseInt(code, 10,0)
	if err !=nil{
		panic(err)
	}

	return int(res)
}

func instance(i string) string {
	return fmt.Sprintf("%sObject", charset.CamelCaseFormatMust(false, i))
}

const bindfuncTpl = `func Bind%sObject(obj definitions.%s){
	%sObject = obj
}`

func objects(routers []*Router)string{
	objs := map[string]bool{}
	for _, r := range routers{
		objs[r.Instance] = true
	}

	res := ""
	for o, _ := range objs{
		res += fmt.Sprintf("var %sObject definitions.%s = nil\n",
			charset.CamelCaseFormatMust(false, o), o)

		res += fmt.Sprintf(bindfuncTpl, o, o, charset.CamelCaseFormatMust(false, o))
		res += "\n\n"
	}

	return res
}

func OutputRouterCode(writer io.Writer)error  {
	funcMap := tt.FuncMap{
		"fieldNameFormat": fieldNameFormat,
		"sub": sub,
		"excludePtr": excludePtr,
		"atoi": atoi,
		"objects": objects,
		"instance":instance,
	}

	if t, err := tt.New("router").Funcs(funcMap).Parse(routerTemplate);err != nil{
		return err
	}else {
		if err := t.Execute(writer, globalRouters); err != nil{
			return err
		}
	}

	return nil
}

func OutputRouterCodeWithTemplate(writer io.Writer, Path string)error  {
	if fd, err := os.Open(Path); err != nil{
		return err
	}else {
		if data, err := ioutil.ReadAll(fd); err != nil{
			return err
		}else {
			funcMap := tt.FuncMap{
				"fieldNameFormat": fieldNameFormat,
				"sub": sub,
				"excludePtr": excludePtr,
				"atoi": atoi,
				"objects": objects,
				"instance":instance,
			}

			if t, err := tt.New("router").Funcs(funcMap).Parse(string(data));err != nil{
				return err
			}else {
				if err := t.Execute(writer, globalRouters); err != nil{
					return err
				}
			}
		}
	}

	return nil
}