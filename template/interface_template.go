package template

const interfaceTemplate  = `package definitions

{{$input:=.}}{{range $interfaceName, $methods := $input}}
type {{$interfaceName|fieldNameFormat}} interface{
    {{range $index, $method := $methods}}
    {{$method.Name}}{{$paramNum := $method.Parameters|len }}{{$paramNum = (sub $paramNum 1 ) }}({{range $index, $param := $method.Parameters }}{{$param.Name}} {{$param.Type}}{{if (gt $paramNum $index)}},{{end}} {{end}}){{$returnNum := $method.Returns|len }}{{$returnNum = (sub $returnNum 1 ) }}({{range $index, $res := $method.Returns}}{{$res.Name}} {{$res.Type}}{{if (gt $returnNum $index)}},{{else}}, statusCode int{{end}}{{end}})
    {{end}}
}
{{end}}`