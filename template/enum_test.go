package template

import (
	"github.com/redresseur/swagger/analyse"
	"os"
	"testing"
)

func TestAnalyseEnums(t *testing.T) {
	res, err := analyse.ReadYaml(yamlPath);
	if  err != nil{
		t.Fatalf("TestGetPaths %v", err)
	}

	defs, err := analyse.GetDefinition(res)
	if  err != nil{
		t.Fatalf("TestGetDefinition %v", err)
	}

	if err := DefinitionComplete(defs); err != nil{
		t.Fatalf("TestGetDefinition %v", err)
	}

	t.Logf("TestOutputStructureCode %v", OutputEnumsCode(os.Stdout))
}