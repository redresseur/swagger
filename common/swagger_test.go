package common

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"testing"
)

var (
	apisDesc = Descriptions{
		PathDescs: []*PathDescription{&PathDescription{
			Url: "/world/hello",
			Method: "GET",
			OperationId: "HelloWorld",
		},&PathDescription{
			Url: "/person/{id}/name",
			Method: "GET",
			OperationId: "PersonName",
		}},
		BasePath: "/v1",
		Host: "example.com",
	}

	ops = func(operationId string)func(*gin.Context){
		return func(context *gin.Context) {

		}
	}
)

func UpdateHelloWorldCond()  {

}

func TestRouterBind(t *testing.T) {
	_, engine := gin.CreateTestContext(nil)
	data, _ := json.Marshal(&apisDesc)


	RouterBind(engine, data, ops, WithApiAuthority())
}