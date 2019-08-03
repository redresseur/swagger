package template

const routerTemplate = `package apis

import (
    "./definitions"
    "encoding/json"
    "github.com/gin-gonic/gin"
    "net/http"
	"fmt"
)
{{$input:=.}}{{range $index, $router := $input}}{{$method := $router.Method}}
func {{$method.Name}}(ctx *gin.Context){
    if {{$router.Instance|instance}} == nil {
        ctx.Writer.WriteHeader(404)
        return
    }

    {{range $param := $method.Parameters}}
    {{ if  (eq $param.IN  "body") }}{{ if (eq $param.SwaggerType "object") }}{{$param.Name}} := &definitions.{{$param.Type|excludePtr}}{}
    if err := ctx.BindJSON(&{{$param.Name}}); err != nil{
        ctx.Writer.WriteHeader(http.StatusBadRequest)
        ctx.Writer.WriteString(err.Error())
        return
    }
    {{else}}
    {{$param.Name}} := ctx.Param("{{$param.Key}}")
    {{end}}{{end}}{{ if  (eq $param.IN  "path") }}{{$param.Name}} := ctx.Param("{{$param.Key}}"){{end}}{{ if  (eq $param.IN  "query") }}{{$param.Name}} := ctx.Query("{{$param.Key}}"){{end}}{{end}}
    {{$returnsNum := $method.Returns|len}}{{$returnsNum = (sub $returnsNum 1 ) }}{{range $index, $res := $method.Returns}}{{$res.Name}}{{if (gt $returnsNum $index)}}, {{else}}, statusCode{{end}}{{end}} := {{$router.Instance|instance}}.{{$method.Name}}({{$paramNum:=$method.Parameters|len}}{{$paramNum = (sub $paramNum 1 ) }}{{range $index, $param := $method.Parameters}}{{$param.Name}}{{if (gt $paramNum $index)}}, {{end}}{{end}})

    {{range $index, $res := $method.Returns}}{{$resCode := $res.StatusCode|atoi}}{{if (and (ge $resCode 100) (lt $resCode 200)) }}if {{$res.Name}} != nil {
        data, _ := json.Marshal({{$res.Name}})
        ctx.Writer.Header().Set("Content-Type", "application/json")
        ctx.Writer.Header().Set("Content-Length", fmt.Sprintf("%d", len(data)))
 		 ctx.Writer.Write(data)
        ctx.Writer.WriteHeader(statusCode)
        return
    }{{end}}{{if (and (ge $resCode 200) (lt $resCode 300)) }}if {{$res.Name}} != nil {
        data, _ := json.Marshal({{$res.Name}})
        ctx.Writer.Header().Set("Content-Type", "application/json")
        ctx.Writer.Header().Set("Content-Length", fmt.Sprintf("%d", len(data)))
		 ctx.Writer.Write(data)
        ctx.Writer.WriteHeader(statusCode)
        return
    }{{end}}{{if (and (ge $resCode 300) (lt $resCode 400)) }}if {{$res.Name}} != nil {
        data, _ := json.Marshal({{$res.Name}})
        ctx.Writer.Header().Set("Content-Type", "application/json")
        ctx.Writer.Header().Set("Content-Length", fmt.Sprintf("%d", len(data)))
        ctx.Writer.Write(data)
        ctx.Writer.WriteHeader(statusCode)
        return
    }{{end}}{{if (and (ge $resCode 400) (lt $resCode 500)) }}if {{$res.Name}} != nil {
        data, _ := json.Marshal({{$res.Name}})
        ctx.Writer.Header().Set("Content-Type", "application/json")
        ctx.Writer.Header().Set("Content-Length", fmt.Sprintf("%d", len(data)))
		 ctx.Writer.Write(data)
        ctx.Writer.WriteHeader(statusCode)
        return
    }{{end}}{{if (and (ge $resCode 500) (lt $resCode 600)) }}if {{$res.Name}} != nil {
        data, _ := json.Marshal({{$res.Name}})
        ctx.Writer.Header().Set("Content-Type", "application/json")
        ctx.Writer.Header().Set("Content-Length", fmt.Sprintf("%d", len(data)))
		 ctx.Writer.Write(data)
        ctx.Writer.WriteHeader(statusCode)
        return
    }{{end}}{{end}}
    {{range $index, $res := $method.Returns}}{{$resCode := $res.StatusCode|atoi}}{{if (eq $resCode -1) }}if {{$res.Name}} != nil {
        data, _ := json.Marshal({{$res.Name}})
        ctx.Writer.Header().Set("Content-Type", "application/json")
        ctx.Writer.Header().Set("Content-Length", fmt.Sprintf("%d", len(data)))
        ctx.Writer.Write(data)
        ctx.Writer.WriteHeader(statusCode)
        return
    }{{end}}{{end}}
}
{{end}}

{{$input|objects}}`
