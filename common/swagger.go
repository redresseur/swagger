package common

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/redresseur/auth_manager"
	"github.com/redresseur/auth_manager/namespace"
	"net/http"
	"strings"
)

type PathDescription struct {
	Url         string `json:"url"`
	Method      string `json:"method"`
	OperationId string `json:"operationId"`
}

type Descriptions struct {
	PathDescs []*PathDescription `json:"pathDescs" `
	BasePath  string             `json:"basePath"`
	Host      string             `json:"host"`
}

type SwaggerApiConf struct {
	*Descriptions
	Operations map[string]string
}

type BindOptions func()

var (
	globalMiddles []gin.HandlerFunc
	groupMiddles  []gin.HandlerFunc
	basePath      string
	host          string

	apiAuthority bool

	// 权限存储空间
	apiAuthorityRootSpace = namespace.NewNameSpace("/")

	apiAuthorityLinks = map[string]struct {
		Desc  *PathDescription
		Space *namespace.RestFulAuthorNamespace
	}{}
)

func WithApiAuthority() BindOptions {
	return func() {
		apiAuthority = true
	}
}

func WithGlobalMiddleWare(middle gin.HandlerFunc) BindOptions {
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

func ginMethod(m string) string {
	//m = strings.ToLower(m)
	return strings.ToUpper(m)
}

func urlParamTransfer(url string) string {
	strings.NewReplacer("{", ":").Replace(url)
	url = strings.NewReplacer("{", ":").Replace(url)
	url = strings.NewReplacer("}", "").Replace(url)
	return url
}

// TODO: 添加組管理
func RouterBind(engine *gin.Engine, description []byte, Operation func(operationId string) func(*gin.Context), bindOptions ...BindOptions) (err error) {
	apiDescs := Descriptions{}
	if err := json.Unmarshal(description, &apiDescs); err != nil {
		return err
	}

	basePath = apiDescs.BasePath
	host = apiDescs.Host
	for _, op := range bindOptions {
		op()
	}

	engine.Use(globalMiddles...)
	routerGroup := engine.Group(basePath, groupMiddles...)

	if apiAuthority {
		apiAuthorityRootSpace, err = namespace.AddSubNameSpace(apiAuthorityRootSpace, apiDescs.BasePath)
	}

	for _, desc := range apiDescs.PathDescs {
		if apiAuthority {
			// 开启认证
			sp, err := namespace.AddSubNameSpace(apiAuthorityRootSpace, desc.Url)
			if err != nil {
				return err
			}

			apiAuthorityLinks[desc.OperationId] = struct {
				Desc  *PathDescription
				Space *namespace.RestFulAuthorNamespace
			}{Desc: desc, Space: sp}
		}

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
	if strings.ToLower(ctx.Request.Method) == `options` {
		ctx.JSON(http.StatusOK, gin.H{})
		return
	}

	// 判断 SessionId
	ctx.Next()
}

func Authorization(ctx *gin.Context) {
	// 取到namespace
	sp, err := namespace.NameSpace(apiAuthorityRootSpace, ctx.Request.RequestURI)
	if err !=nil{
		ctx.JSON(http.StatusInternalServerError, gin.H{"errDesc": err.Error()})
		return
	}

	// 进行api权限测权限
	if err := auth_manager.CheckAuthority(ctx, sp); err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{"errDesc": err.Error()})
	} else {
		ctx.Next()
	}

	return
}
