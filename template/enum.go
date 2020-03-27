package template

import (
	"github.com/redresseur/utils/charset"
	"io"
	"log"
	tt "text/template"
)

var globalEnums map[string]*Enums = map[string]*Enums{}

type Enums struct {
	TypeName string
	EnumName string
	Elements []string
}

func AnalyseEnums(structName, filedName, typeName string, elements []string) string {
	var err error
	if len(elements) == 0 {
		return typeName
	}

	enums := Enums{
		TypeName: typeName,
	}

	enums.EnumName, err = charset.CamelCaseFormat(true, structName, filedName)
	if err != nil {
		log.Fatalf("StructName %s, filedName %s: %v", structName, filedName, err)
		return ""
	}

	if _, ok := globalEnums[enums.EnumName]; ok {
		log.Fatalf("%s is confilict", enums.EnumName)
	}

	enums.Elements = elements
	globalEnums[enums.EnumName] = &enums
	return enums.EnumName
}

func OutputEnumsCode(writer io.Writer) error {
	funcMap := tt.FuncMap{
		"fieldNameFormat": fieldNameFormat,
		"sub":             sub,
		"excludePtr":      excludePtr,
		"atoi":            atoi,
	}

	if t, err := tt.New("enums").Funcs(funcMap).Parse(enumsTemplate); err != nil {
		return err
	} else {
		if err := t.Execute(writer, globalEnums); err != nil {
			return err
		}
	}

	return nil
}
