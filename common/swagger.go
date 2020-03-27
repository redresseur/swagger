package common

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/sessions"

	"github.com/gin-gonic/gin"
	"github.com/redresseur/auth_manager"
	"github.com/redresseur/auth_manager/namespace"
	"net/http"
	"regexp"
	"strings"
)

type PathDescription struct {
	Url         string   `json:"url"`
	Method      string   `json:"method"`
	OperationId string   `json:"operationId"`
	Tags        []string `json:"tags"`
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

	apiAuthorityEnable bool

	// 权限存储空间
	apiAuthorityRootSpace = namespace.NewNameSpace("/", nil)

	apiAuthorityBaseSpace *namespace.RestFulAuthorNamespace

	apiAuthorityLinks = map[string]struct {
		Desc  *PathDescription
		Space *namespace.RestFulAuthorNamespace
	}{}

	apiAuthorityGroup = map[string]*namespace.RestFulAuthorNamespace{}
)

func UpdateGroupAuthor(name string, ops ...namespace.CondsOp) error {
	sp, ok := apiAuthorityGroup[name]
	if !ok {
		return errors.New("the group is not found")
	}

	return namespace.UpdateCondition(sp, ops...)
}

func Replace(url string, params ...string) (string, error) {
	rc, _ := regexp.Compile(`(\{.[a-zA-Z\_0-9]+\})`)
	urlParams := rc.FindAllString(url, -1)

	if len(urlParams) > len(params) {
		return "", errors.New("the number of params is not enough")
	}

	for i, up := range urlParams {
		url = strings.Replace(url, up, params[i], 1)
	}

	return url, nil
}

func UpdateApiAuthor(operationId string, param []string, ops ...namespace.CondsOp) (err error) {
	desc, ok := apiAuthorityLinks[operationId]
	if !ok {
		return errors.New("the operation id is not found")
	}

	uri := desc.Desc.Url
	if len(param) > 0 {
		if uri, err = Replace(uri, param...); err != nil {
			return err
		}
	}

	sp, err := namespace.AddSubNameSpace(apiAuthorityBaseSpace, uri)
	if err != nil {
		return err
	}

	if err := namespace.UpdateCondition(sp, ops...); err != nil {
		return err
	}

	// 复制组策略
	for _, g := range desc.Desc.Tags {
		g_sp, ok := apiAuthorityGroup[g]
		if !ok {
			continue
		}

		sp, err = namespace.ReverseFind(sp, g)
		if err != nil {
			break
		}

		// 此处使用浅拷贝
		// 所以当组策略发生变化的时候
		// 所有关联的策略都会变化
		sp.SrcNamespace = g_sp.SrcNamespace
	}

	return nil
}

func WithApiAuthority() BindOptions {
	return func() {
		apiAuthorityEnable = true
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

	if apiAuthorityEnable {
		apiAuthorityBaseSpace, err = namespace.AddSubNameSpace(apiAuthorityRootSpace, apiDescs.BasePath)
	}

	for _, desc := range apiDescs.PathDescs {
		if apiAuthorityEnable {
			// 开启认证
			sp, err := namespace.AddSubNameSpace(apiAuthorityBaseSpace, desc.Url)
			if err != nil {
				return err
			}

			for _, g := range desc.Tags {
				g_sp, _ := namespace.ReverseFind(sp, g)
				if g_sp == nil {
					continue
				}

				g_sp_g, ok := apiAuthorityGroup[g]
				if !ok {
					apiAuthorityGroup[g] = namespace.NewNameSpace(g, nil)
					g_sp_g = apiAuthorityGroup[g]
				}

				g_sp.SrcNamespace = g_sp_g.SrcNamespace
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

func Cros() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// add header
		addHeader(ctx)

		// filter
		// 所有option全部跳过
		if strings.ToLower(ctx.Request.Method) == `options` {
			ctx.JSON(http.StatusOK, gin.H{})
			ctx.Abort()
			return
		}

		// 判断 SessionId
		// ctx.Next()
	}
}

func Authorization() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if !apiAuthorityEnable {
			//ctx.Next()
			return
		}

		// 取到namespace
		sp, err := namespace.NameSpace(apiAuthorityRootSpace, ctx.Request.RequestURI)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"errDesc": err.Error()})
			ctx.Abort()
			return
		}

		// 进行api权限测权限
		if err := auth_manager.CheckAuthority(ctx, sp); err != nil {
			ctx.JSON(http.StatusForbidden, gin.H{"errDesc": err.Error()})
			ctx.Abort()
		} else {
			ctx.Next()
		}

		return
	}
}

func Sessions(store sessions.Store, appName string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//sid := ctx.GetHeader(X_Session_UUID)
		//if sid == "" {
		//	return
		//}
		//
		ss, err := store.Get(ctx.Request, appName)
		if err != nil {
			return
		}

		ctx.Set(X_Session, ss)
	}
}

const (
	X_Session_UUID = "X-Session-UUID"
	X_Permission   = "X-Permission"
	X_Session      = "X-S-Session"
)

func PermissionFromSessions(ctx *gin.Context) (interface{}, error) {
	s, ok := ctx.Get(X_Session)
	if !ok {
		// 此时可能尚未登陆
		return nil, nil
	}

	S, ok := s.(sessions.Session)
	if !ok {
		return nil, errors.New("Session Broken")
	}

	return S.Values, nil
}
