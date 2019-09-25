package template

import (
	"io"
	"io/ioutil"
	"os"
	tt "text/template"
)

func OutputAuthorityCode(writer io.Writer)error  {
	funcMap := tt.FuncMap{
		"fieldNameFormat": fieldNameFormat,
		"sub": sub,
	}

	if t, err := tt.New("interface").Funcs(funcMap).Parse(string(authorTemplate));err != nil{
		return err
	}else {
		if err := t.Execute(writer, globalMethods); err != nil{
			return err
		}
	}

	return nil
}

func OutputAuthorityCodeWithTemplate(writer io.Writer, Path string)error  {
	if fd, err := os.Open(Path); err != nil{
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
				if err := t.Execute(writer, globalMethods); err != nil{
					return err
				}
			}
		}
	}

	return nil
}