package template

import (
	"github.com/redresseur/swagger/analyse"
	"testing"
)

var (
	yamlPath = "swagger.yaml"
	jsonPath = "swagger.json"
)

func TestDefinitionComplete(t *testing.T)  {
	res, err := analyse.ReadYaml(yamlPath);
	if  err != nil{
		t.Fatalf("TestGetPaths %v", err)
	}

	defs, err := analyse.GetDefinition(res)
	if  err != nil{
		t.Fatalf("TestGetDefinition %v", err)
	}

	definitionComplete(defs)
}

func TestOutputStructureCode(t *testing.T)  {
	TestDefinitionComplete(t)
	t.Logf("TestOutputStructureCode %v", outputStructureCode())
}