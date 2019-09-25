package main

import (
	"flag"
	"github.com/redresseur/swagger/analyse"
	"github.com/redresseur/swagger/template"
	"github.com/redresseur/utils/ioutils"
	"path/filepath"
)

var (
	swagger = `swagger`
	output  = `output`
	authorityEnable = false
)

func main() {
	flag.StringVar(&swagger, swagger, "swagger.yaml", "swagger template, only support yaml")
	flag.StringVar(&output, output, "./", "output dir")
	flag.BoolVar(&authorityEnable, "author", authorityEnable, "api authority enable")
	flag.Parse()

	data, err := analyse.ReadYaml(swagger)
	if err != nil {
		panic(err)
	}

	defs, err := analyse.GetDefinition(data)
	if err != nil {
		panic(err)
	}

	if err := template.DefinitionComplete(defs); err != nil {
		panic(err)
	}

	apis, err := analyse.GetRestApi(data)
	if err != nil {
		panic(err)
	}

	if err := template.InterfaceComplete(apis); err != nil {
		panic(err)
	}

	if err := template.RouterComplete(apis); err != nil {
		panic(err)
	}

	if _, err := ioutils.CreateDirIfMissing(filepath.Join(output, "definitions")); err != nil {
		panic(err)
	}

	interfaceOut, err := ioutils.OpenFile(filepath.Join(output, "definitions", "interface.go"), "")
	if err != nil {
		panic(err)
	} else {
		defer interfaceOut.Close()
	}

	if err := template.OutputInterfaceCode(interfaceOut); err != nil {
		panic(err)
	}

	commonOut, err := ioutils.OpenFile(filepath.Join(output, "definitions", "common.go"), "")
	if err := template.OutputEnumsCode(commonOut); err != nil {
		panic(err)
	} else {
		defer commonOut.Close()
	}

	structureOut, err := ioutils.OpenFile(filepath.Join(output, "definitions", "structure.go"), "")
	if err != nil {
		panic(err)
	} else {
		defer structureOut.Close()
	}

	if err := template.OutputStructureCode(structureOut); err != nil {
		panic(err)
	}

	apisOutput, err := ioutils.OpenFile(filepath.Join(output, "apis.go"), "")
	if err != nil {
		panic(err)
	} else {
		defer apisOutput.Close()
	}

	if err := template.OutputRouterCode(apisOutput); err != nil {
		panic(err)
	}

	template.DescriptionComplete(analyse.GetHost(data), analyse.GetBasePath(data))
	descOutput, err := ioutils.OpenFile(filepath.Join(output, "descriptions.go"), "")
	if err != nil {
		panic(err)
	} else {
		defer descOutput.Close()
	}

	if err := template.OutputDescription(descOutput); err != nil {
		panic(err)
	}

	if authorityEnable {
		authorOutput, err := ioutils.OpenFile(filepath.Join(output, "authority.go"), "")
		if err != nil {
			panic(err)
		} else {
			defer authorOutput.Close()
		}

		if err := template.OutputAuthorityCode(authorOutput); err != nil {
			panic(err)
		}
	}
}
