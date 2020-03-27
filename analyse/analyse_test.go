package analyse

import (
	"testing"
)

var (
	yamlPath = "swagger.yaml"
	jsonPath = "swagger.json"
)

func TestReadYaml(t *testing.T) {
	res, err := ReadJson(yamlPath)
	if err != nil {
		t.Fatalf("TestReadYaml %v", err)
	}

	for k := range res {
		t.Log(k)
	}
}

func TestReadJson(t *testing.T) {
	res, err := ReadYaml(jsonPath)
	if err != nil {
		t.Fatalf("TestReadJson %v", err)
	}

	for k := range res {
		t.Log(k)
	}

}

func TestGetRestApi(t *testing.T) {
	res, err := ReadYaml(yamlPath)
	if err != nil {
		t.Fatalf("TestGetPaths %v", err)
	}

	if _, err := GetRestApi(res); err != nil {
		t.Fatalf("TestGetRestApi %v", err)
	}
}

func TestGetDefinition(t *testing.T) {
	res, err := ReadYaml(yamlPath)
	if err != nil {
		t.Fatalf("TestGetPaths %v", err)
	}

	if _, err := GetDefinition(res); err != nil {
		t.Fatalf("TestGetDefinition %v", err)
	}
}
