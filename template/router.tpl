{{$input:=.}}{{range $index, $router := $input}}
func {{$router.Name}}(ctx *gin.Context){
    {{range _, $object := $router.objects}}
    {{$object.VariableName}} := &definitions.{{$object.StructureName}}
    if err := ctx.BindJSON(&{{$object.VariableName}}); err != nil{
        return
    }{{end}}

    {{range _, $param := $router.UrlParameters}}
    {{$param.VariableName}} := ctx.Param("{{$param.Name}}")
    {{end}}

    {{range _, $param := $router.QueryParameters}}
    {{$param.VariableName}} := ctx.Param("{{$param.Name}}")
    {{end}}

    {{$router.Instance}}.{{$router.Method}}({{range $paramName, _ := $router.Parameters}}{{$paramName}},{{end}})
}{{end}}