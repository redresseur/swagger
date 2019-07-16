package template

import (
	"github.com/redresseur/swagger/analyse"
	"testing"
)

func TestInterfaceComplete(t *testing.T)  {
	res, err := analyse.ReadYaml(yamlPath);
	if  err != nil{
		t.Fatalf("TestInterfaceComplete %v", err)
	}

	defs, err := analyse.GetDefinition(res)
	if  err != nil{
		t.Fatalf("TestInterfaceComplete %v", err)
	}
	if err := definitionComplete(defs); err != nil{
		t.Fatalf("TestInterfaceComplete %v", err)
	}

	apis, err := analyse.GetRestApi(res)
	if err != nil{
		t.Fatalf("TestInterfaceComplete %v", err)
	}

	if err := interfaceComplete(apis); err != nil{
		t.Fatalf("TestInterfaceComplete %v", err)
	}
}

func TestOutputInterfaceCode(t *testing.T)  {
	TestInterfaceComplete(t)

	if err := outputInterfaceCode(); err != nil{
		t.Fatalf("TestOutputInterfaceCode %v", err)
	}

	if err := outputStructureCode(); err != nil{
		t.Fatalf("TestOutputInterfaceCode %v", err)
	}
}

