{{$input:=.}}{{range $interfaceName, $methods := $input}}
type {{$interfaceName}} interface{
    {{range $index, $method := $methods}}
    {{$method.Name}}({{range $paramName, $paramType := $methods.Parameters }}{{$paramName}} {{$paramType}}, {{end}}){{range $paramName, $paramType := $methods.Returns}}
    {{end}}
}
{{end}}