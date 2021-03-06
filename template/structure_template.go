package template

const structureTemplate  = `package definitions
{{$input:=.}}{{range $structureName, $structure := $input}}
type {{$structureName}}{{$filedNum := $structure|len}} struct{
    {{range $fieldName, $fieldType := $structure}}{{$fieldName|fieldNameFormat}} {{$fieldType}} `+ "`" + `json:"{{$fieldName}}"` + "`\n" +
	`{{with $filedNum =  (sub $filedNum 1) }}{{if (gt $filedNum 0)}}    {{end}}{{end}}{{end}}}
{{end}}`
