package template

import (
	"github.com/redresseur/swagger/analyse"
	"os"
	"testing"
)

func TestOutputDescriptionWithTemplate(t *testing.T) {
	res, err := analyse.ReadYaml(yamlPath)
	if err != nil {
		t.Fatalf("TestInterfaceComplete %v", err)
	}

	defs, err := analyse.GetDefinition(res)
	if err != nil {
		t.Fatalf("TestInterfaceComplete %v", err)
	}
	if err := DefinitionComplete(defs); err != nil {
		t.Fatalf("TestInterfaceComplete %v", err)
	}

	apis, err := analyse.GetRestApi(res)
	if err != nil {
		t.Fatalf("TestInterfaceComplete %v", err)
	}

	if err := InterfaceComplete(apis); err != nil {
		t.Fatalf("TestInterfaceComplete %v", err)
	}

	if err := RouterComplete(apis); err != nil {
		t.Fatalf("TestInterfaceComplete %v", err)
	}

	DescriptionComplete(analyse.GetHost(res), analyse.GetBasePath(res))

	OutputDescriptionWithTemplate(os.Stdout, "urls.tpl")
}
