{{$input:=.}}{{range $interfaceName, $methods := $input}}
type {{$interfaceName}} interface{
    {{range $index, $method := $methods}}
    {{$method.Name}}{{$filedNum := $method.Parameters|len }}({{range $paramName, $paramType := $method.Parameters }}{{$paramName}} {{$paramType}}{{with $filedNum =  (sub $filedNum 1) }}{{if (gt $filedNum 0)}},{{end}}{{end}} {{end}})({{range $paramName, $paramType := $method.Returns}}{{$paramName}} {{$paramType}}, {{end}})
    {{end}}
}{{end}}