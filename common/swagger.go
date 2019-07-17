package common

type PathDescription struct {
	Path string `json:"path"`
	Method string `json:"method"`
	OperationId string `json:"operationId"`
}

type Descriptions struct {
	PathDescs []*PathDescription `json:"pathDescs"`

}

