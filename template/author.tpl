package apis

import (
    "github.com/redresseur/auth_manager/namespace"
    "github.com/redresseur/swagger/common"
)

{{$input:=.}}{{range $interfaceName, $methods := $input}}
func {{$interfaceName|fieldNameFormat}}GroupAuthority(ops ...namespace.CondsOp)error{
    return common.UpdateGroupAuthor("{{$interfaceName}}", ops...)
}
{{range $index, $method := $methods}}
func {{$method.Name}}ApiAuthority({{range $index, $param := $method.Parameters}}{{ if  (eq $param.IN  "path") }}{{$param.Name}} {{$param.Type}}, {{end}}{{end}}ops ...namespace.CondsOp)error{
    params := []string{
        {{range $index, $param := $method.Parameters}}{{ if  (eq $param.IN  "path") }}{{$param.Name}},{{end}}
        {{end}}
    }
    return common.UpdateApiAuthor("{{ $method.OperationId }}", params, ops...)
}
{{end}}{{end}}