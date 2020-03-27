package template

const enumsTemplate = `package definitions
{{$input:=.}}{{range $enumTypeName, $enums := $input}}
type {{$enumTypeName}} {{$enums.TypeName}}
const (
{{$elementNum := $enums.Elements|len }}{{$elementNum = (sub $elementNum 1 ) }}{{range $index, $element := $enums.Elements}}	{{$enumTypeName}}_{{$element}} {{$enumTypeName}} = "{{$element}}"
{{if (gt $elementNum $index)}}{{end}}{{else}}{{end}})
{{end}}`
