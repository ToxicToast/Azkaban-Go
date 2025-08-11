package routerclient

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/ToxicToast/Azkaban-Go/libs/shared/grpcclient"
	"github.com/ToxicToast/Azkaban-Go/libs/shared/healthmon"
	"github.com/ToxicToast/Azkaban-Go/libs/shared/helper"
	"github.com/ToxicToast/Azkaban-Go/libs/shared/registryclient"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
)

type Client struct {
	router *gin.Engine
	port   string
	pool   *grpcclient.Client
}

var pathParamRe = regexp.MustCompile(`\{([a-zA-Z0-9_]+)\}`)

func yamlToGinPath(p string) string {
	return pathParamRe.ReplaceAllString(p, ":$1")
}

func setField(dst map[string]any, dotted string, val any) {
	parts := strings.Split(dotted, ".")
	m := dst
	for i, p := range parts {
		if i == len(parts)-1 {
			m[p] = val
			return
		}
		if _, ok := m[p]; !ok {
			m[p] = map[string]any{}
		}
		m = m[p].(map[string]any)
	}
}

func getQueryMaybe(c *gin.Context, key string) (string, bool) {
	v := c.Query(key)
	if v == "" {
		vs, ok := c.Request.URL.Query()[key]
		if ok && len(vs) > 0 {
			return vs[0], true
		}
		return "", false
	}
	return v, true
}

func jsonMarshal(v any) ([]byte, error) {
	// kleines Extra: true/false/"123" → bleiben Strings; Proto rückt es i.d.R. gerade.
	return json.Marshal(v)
}

func normalizeTarget(t string) string {
	i := strings.LastIndex(t, ".")
	if i == -1 {
		return t // oder error
	}
	service := t[:i]  // "warcraft.character.WarcraftCharacterService"
	method := t[i+1:] // "GetCharacters"
	return service + "/" + method
}

func mapGrpcError(err error) (int, string) {
	st, ok := status.FromError(err)
	if !ok {
		return http.StatusBadGateway, "upstream error"
	}
	switch st.Code() {
	case codes.NotFound:
		return http.StatusNotFound, st.Message()
	case codes.InvalidArgument, codes.FailedPrecondition, codes.OutOfRange:
		return http.StatusBadRequest, st.Message()
	case codes.Unauthenticated:
		return http.StatusUnauthorized, st.Message()
	case codes.PermissionDenied:
		return http.StatusForbidden, st.Message()
	case codes.DeadlineExceeded:
		return http.StatusGatewayTimeout, st.Message()
	case codes.Unavailable:
		return http.StatusBadGateway, st.Message()
	default:
		return http.StatusInternalServerError, st.Message()
	}
}

func coercePathParam(s string) any {
	// erst int, dann bool, zuletzt string
	if n, err := strconv.ParseInt(s, 10, 64); err == nil {
		return n
	}
	if b, err := strconv.ParseBool(s); err == nil {
		return b
	}
	return s
}

func setMode(envName string) {
	switch envName {
	case "prod", "staging":
		gin.SetMode(gin.ReleaseMode)
	case "dev":
		gin.SetMode(gin.DebugMode)
	default:
		gin.SetMode(gin.TestMode)
	}
}

type httpError struct {
	Code    int
	Message string
}

func newHTTPError(code int, msg string) *httpError {
	return &httpError{Code: code, Message: msg}
}

func applyPathParams(ctx *gin.Context, payload map[string]any, pathMap map[string]string) *httpError {
	for from, to := range pathMap {
		val := ctx.Param(from)
		if val == "" {
			return newHTTPError(http.StatusBadRequest, fmt.Sprintf("missing path param: %s", from))
		}
		setField(payload, to, val)
	}
	return nil
}

func applyBody(ctx *gin.Context, payload map[string]any, bodyMap map[string]string) *httpError {
	var body map[string]any
	if err := ctx.ShouldBindJSON(&body); err != nil {
		return newHTTPError(http.StatusBadRequest, "invalid JSON body")
	}

	if len(bodyMap) > 0 {
		for from, to := range bodyMap {
			if v, ok := body[from]; ok {
				setField(payload, to, v)
			}
		}
		return nil
	}

	// kein Mapping definiert -> ganzen Body übernehmen
	for k, v := range body {
		payload[k] = v
	}
	return nil
}

func applyQueryParams(ctx *gin.Context, payload map[string]any, rt helper.Routes) *httpError {
	for from, to := range rt.Mapping.QueryParams {
		if raw, ok := getQueryMaybe(ctx, from); ok {
			qt := rt.Http.Query[from] // "number" | "boolean" | "string"
			var val any = raw
			switch qt {
			case "number":
				n, err := strconv.ParseInt(raw, 10, 64)
				if err != nil {
					return newHTTPError(http.StatusBadRequest, fmt.Sprintf("query %s must be number", from))
				}
				val = n
			case "boolean":
				b, err := strconv.ParseBool(raw)
				if err != nil {
					return newHTTPError(http.StatusBadRequest, fmt.Sprintf("query %s must be boolean", from))
				}
				val = b
				// "string" oder leer -> raw beibehalten
			}
			setField(payload, to, val)
		}
	}
	return nil
}

func (c *Client) buildPayload(ctx *gin.Context, rt helper.Routes) (map[string]any, *httpError) {
	payload := make(map[string]any)

	if err := applyPathParams(ctx, payload, rt.Mapping.PathParams); err != nil {
		return nil, err
	}

	if err := applyQueryParams(ctx, payload, rt); err != nil {
		return nil, err
	}

	if ctx.Request.Method == http.MethodPost || ctx.Request.Method == http.MethodPut || ctx.Request.Method == http.MethodPatch {
		if err := applyBody(ctx, payload, rt.Mapping.Body); err != nil {
			return nil, err
		}
	}

	return payload, nil
}

func (c *Client) registerRoute(method, path string, handler gin.HandlerFunc) {
	switch method {
	case http.MethodGet:
		c.router.GET(path, handler)
	case http.MethodPost:
		c.router.POST(path, handler)
	case http.MethodPut:
		c.router.PUT(path, handler)
	case http.MethodPatch:
		c.router.PATCH(path, handler)
	case http.MethodDelete:
		c.router.DELETE(path, handler)
	default:
		panic("unsupported method: " + method)
	}
}

func (c *Client) validateTargets(routes []helper.Routes, reg registryclient.Registry) {
	for _, rt := range routes {
		key := normalizeTarget(rt.Grpc.Target)
		if _, ok := reg.Get(key); !ok {
			log.Fatalf("route %q refers to unknown grpc target %q", rt.Name, key)
		}
	}
}

func (c *Client) handleRequest(ctx *gin.Context, rt helper.Routes, reg registryclient.Registry) *httpError {
	entry, ok := reg.Get(normalizeTarget(rt.Grpc.Target))
	if !ok {
		return newHTTPError(http.StatusNotImplemented, "unmapped grpc target")
	}
	payload, errBP := c.buildPayload(ctx, rt)
	if errBP != nil {
		return errBP
	}
	req := entry.NewReq()
	marshaled, errJSON := jsonMarshal(payload)
	if errJSON != nil {
		return newHTTPError(http.StatusBadRequest, errJSON.Error())
	}
	if err := protojson.Unmarshal(marshaled, req); err != nil {
		return newHTTPError(http.StatusBadRequest, "invalid request: "+err.Error())
	}
	cc, err := c.pool.Get(ctx.Request.Context(), rt.Grpc.Service)
	if err != nil {
		return newHTTPError(http.StatusBadGateway, "downstream dial failed: "+err.Error())
	}
	to := time.Duration(rt.Grpc.Timeoutms) * time.Millisecond
	if to <= 0 {
		to = 5 * time.Second
	}
	cctx, cancel := context.WithTimeout(ctx.Request.Context(), to)
	defer cancel()
	resp, err := entry.Invoke(cctx, cc, req)
	if err != nil {
		httpCode, msg := mapGrpcError(err)
		return newHTTPError(httpCode, msg)
	}

	b, err := protojson.MarshalOptions{UseProtoNames: true}.Marshal(resp)
	if err != nil {
		return newHTTPError(http.StatusInternalServerError, "marshal response failed")
	}

	ctx.Data(http.StatusOK, "application/json", b)
	return nil
}

func NewClient(envName, port string, pool *grpcclient.Client) *Client {
	setMode(envName)
	return &Client{
		router: gin.Default(),
		port:   port,
		pool:   pool,
	}
}

func (c *Client) BuildHealthRoute(
	monitor *healthmon.Monitor,
	requiredServices []string,
	routes []helper.Routes,
	reg registryclient.Registry,
	livenessPath, readinessPath string,
) {
	c.router.GET(readinessPath, func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "ok")
	})

	c.router.GET(livenessPath, func(ctx *gin.Context) {
		for _, rt := range routes {
			key := normalizeTarget(rt.Grpc.Target)
			if _, ok := reg.Get(key); !ok {
				ctx.JSON(http.StatusServiceUnavailable, gin.H{
					"status": "not_ready",
					"reason": "unmapped_target",
					"target": key,
				})
				return
			}
		}

		snap := monitor.Snapshot()
		if !monitor.AllOK(requiredServices) {
			ctx.JSON(http.StatusServiceUnavailable, gin.H{
				"status":   "not_ready",
				"snapshot": snap,
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"status":   "ready",
			"snapshot": snap,
		})
	})
}

func (c *Client) BuildRoutes(routes []helper.Routes, reg registryclient.Registry) {
	for _, route := range routes {
		rt := route
		path := yamlToGinPath(rt.Http.Path)

		handler := func(ctx *gin.Context) {
			if err := c.handleRequest(ctx, rt, reg); err != nil {
				ctx.JSON(err.Code, gin.H{"error": err.Message})
				return
			}
		}

		c.registerRoute(rt.Http.Method, path, handler)
	}

	c.validateTargets(routes, reg)
}

func (c *Client) RunServer() error {
	if err := c.router.Run(c.port); err != nil {
		return err
	}
	return nil
}

func (c *Client) GetRouter() *gin.Engine {
	return c.router
}
