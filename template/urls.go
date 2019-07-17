package template

import (
	"encoding/json"
	"fmt"
	"github.com/redresseur/swagger/common"
	"io"
	"io/ioutil"
	"os"
	tt "text/template"
)

var globalSwaggerConf = &common.SwaggerApiConf{Descriptions: &common.Descriptions{}, Operations: map[string]string{}}

func DescriptionComplete(host, basePath string){
	globalSwaggerConf.Host = host
	globalSwaggerConf.BasePath = basePath

	for _, ms := range globalMethods {
		for _, m := range ms{
			globalSwaggerConf.Operations[m.OperationId] = m.Name
			pd := common.PathDescription{}
			pd.OperationId = m.OperationId
			pd.Url = m.Url
			pd.Method = m.MethodType
			globalSwaggerConf.Descriptions.PathDescs = append(globalSwaggerConf.Descriptions.PathDescs, &pd)
		}
	}
}

const getDescriptionfunc  = `func ApisDescriptions()[]byte{
	return apiDescription
}`

func apiDescription( descriptions *common.Descriptions )string {
	data, err := json.Marshal(descriptions)
	if err != nil{
		panic(err)
	}

	res := "var apiDescription  = []byte{\n	"
	count := 0
	for _, b := range data{
		if count >= 16{
			res += "\n	"
			count = 0
		}

		res += fmt.Sprintf("0x%2x", b) + ","
		count++
	}
	res += "}\n\n"
	res += getDescriptionfunc
	return res
}

func OutputDescriptionWithTemplate(writer io.Writer, path string) error {
	if fd, err := os.Open(path); err != nil{
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
				"apiDescription": apiDescription,
			}

			if t, err := tt.New("swaggerConf").Funcs(funcMap).Parse(string(data));err != nil{
				return err
			}else {
				if err := t.Execute(writer, globalSwaggerConf); err != nil{
					return err
				}
			}
		}
	}

	return nil
}

func OutputDescription(writer io.Writer) error {
	funcMap := tt.FuncMap{
		"fieldNameFormat": fieldNameFormat,
		"sub": sub,
		"excludePtr": excludePtr,
		"atoi": atoi,
		"apiDescription": apiDescription,
	}

	if t, err := tt.New("swaggerConf").Funcs(funcMap).Parse(string(urlsTemplate));err != nil{
		return err
	}else {
		if err := t.Execute(writer, globalSwaggerConf); err != nil{
			return err
		}
	}

	return nil
}