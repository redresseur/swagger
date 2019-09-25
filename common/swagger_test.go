package common

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/redresseur/auth_manager/namespace"
	"testing"
)

var (
	apisDesc = Descriptions{
		PathDescs: []*PathDescription{&PathDescription{
			Url:         "/world/hello",
			Method:      "GET",
			OperationId: "HelloWorld",
			Tags: []string{
				"world",
			},
		}, &PathDescription{
			Url:         "/person/{id}/name",
			Method:      "GET",
			OperationId: "PersonName",
			Tags: []string{
				"person",
			},
		}},
		BasePath: "/v1",
		Host:     "example.com",
	}

	ops = func(operationId string) func(*gin.Context) {
		return func(context *gin.Context) {

		}
	}
)

func TestUpdateApiAuthor(t *testing.T) {
	TestRouterBind(t)

	UpdateApiAuthor("HelloWorld", []string{}, namespace.WithDefaultCond(&namespace.EmptyCondition{}))
}

func TestUpdateApiAuthor2(t *testing.T) {
	TestRouterBind(t)

	UpdateApiAuthor("PersonName", []string{"limei"}, namespace.WithDefaultCond(&namespace.EmptyCondition{}))
}

func TestUpdateGroupAuthor(t *testing.T) {
	TestUpdateApiAuthor(t)

	UpdateGroupAuthor("world", namespace.WithDefaultCond(&namespace.EmptyCondition{}))
}

func TestRouterBind(t *testing.T) {
	_, engine := gin.CreateTestContext(nil)
	data, _ := json.Marshal(&apisDesc)

	RouterBind(engine, data, ops, WithApiAuthority())
}
