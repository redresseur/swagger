package template

const urlsTemplate = `package apis

import(
    "github.com/gin-gonic/gin"
)

var operations = map[string]func(*gin.Context){
{{$input:=.}}{{range $operationId, $operationName := $input.Operations}}    "{{$operationId}}":{{$operationName}},
{{end}}}

func Operation(operationId string)func(*gin.Context){
   return operations[operationId]
}

{{$input.Descriptions|apiDescription}}`
