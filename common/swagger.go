package common

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type PathDescription struct {
	Url string `json:"url"`
	Method string `json:"method"`
	OperationId string `json:"operationId"`
}

type Descriptions struct {
	PathDescs []*PathDescription `json:"pathDescs" `
	BasePath string `json:"basePath"`
	Host string `json:"host"`
}

type SwaggerApiConf struct {
	*Descriptions
	Operations map[string]string
}

type BindOptions func ()

var (
	globalMiddles []gin.HandlerFunc
	groupMiddles []gin.HandlerFunc
	basePath string
	host string
)

func WithGlobalMiddleWare(middle gin.HandlerFunc ) BindOptions  {
	return func() {
		globalMiddles = append(globalMiddles, middle)
	}
}

func WithGroupMiddleWare(middle gin.HandlerFunc) BindOptions {
	return func() {
		groupMiddles = append(groupMiddles, middle)
	}
}

func WithBasePath(bp string) BindOptions {
	return func() {
		basePath = bp
	}
}

func WithHost(h string) BindOptions {
	return func() {
		host = h
	}
}


func ginMethod(m string)string  {
	//m = strings.ToLower(m)
	return strings.ToUpper(m)
}

func urlParamTransfer(url string)string{
	url = strings.ReplaceAll(url, "{", ":"  )
	url = strings.ReplaceAll(url, "}", ""  )
	return url
}

// TODO: 添加組管理
func RouterBind(engine *gin.Engine, description []byte, Operation func(operationId string)func(*gin.Context), bindOptions... BindOptions) error {
	apiDescs := Descriptions{}
	if err := json.Unmarshal(description, &apiDescs); err != nil{
		return err
	}

	basePath = apiDescs.BasePath
	host = apiDescs.Host
	for _, op := range bindOptions{
		op()
	}

	engine.Use(globalMiddles...)
	routerGroup := engine.Group(basePath, groupMiddles...)
	for _, desc := range apiDescs.PathDescs{
		routerGroup.Handle(ginMethod(desc.Method), urlParamTransfer(desc.Url), Operation(desc.OperationId))
	}

	Register(engine)

	return nil
}

func addHeader(ctx *gin.Context) {
	ctx.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	ctx.Writer.Header().Set("Access-Control-Allow-Headers", "Access-Control-Allow-Origin,Access-Control-Allow-Method,Content-Type")
	ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	ctx.Writer.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	ctx.Writer.Header().Set("content-type", "application/json")             //返回数据格式是json
}

func Cros(ctx *gin.Context) {
	// add header
	addHeader(ctx)

	// filter
	// 所有option全部跳过
	if strings.ToLower(ctx.Request.Method) == `options`{
		ctx.JSON(http.StatusOK, gin.H{})
		return
	}

	// 判断 SessionId
	ctx.Next()
}