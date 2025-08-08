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

func NewClient(envName, port string, pool *grpcclient.Client) *Client {
	if envName == "prod" || envName == "staging" {
		gin.SetMode(gin.ReleaseMode)
	} else if envName == "dev" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.TestMode)
	}
	return &Client{
		router: gin.Default(),
		port:   port,
		pool:   pool,
	}
}

func (c *Client) BuildRoutes(routes []helper.Routes, reg registryclient.Registry) {
	for _, route := range routes {
		rt := route
		method := rt.Http.Method
		path := yamlToGinPath(rt.Http.Path)

		handler := func(ctx *gin.Context) {
			normalizedTarget := normalizeTarget(rt.Grpc.Target)
			entry, ok := reg.Get(normalizedTarget)
			if !ok {
				ctx.JSON(http.StatusNotImplemented, gin.H{
					"error": "unmapped grpc target",
				})
				return
			}
			payload := make(map[string]any)

			for from, to := range rt.Mapping.PathParams {
				val := ctx.Param(from)
				if val == "" {
					ctx.JSON(http.StatusBadRequest, gin.H{
						"error": fmt.Sprintf("missing path param: %s", from),
					})
					return
				}
				setField(payload, to, coercePathParam(val))
			}

			for from, to := range rt.Mapping.QueryParams {
				if raw, ok := getQueryMaybe(ctx, from); ok {
					qt := rt.Http.Query[from] // "number" | "boolean" | "string"
					var val any = raw
					switch qt {
					case "number":
						n, err := strconv.ParseInt(raw, 10, 64)
						if err != nil {
							ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("query %s must be number", from)})
							return
						}
						val = n
					case "boolean":
						b, err := strconv.ParseBool(raw)
						if err != nil {
							ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("query %s must be boolean", from)})
							return
						}
						val = b
					}
					setField(payload, to, val)
				}
			}

			req := entry.NewReq()
			marshaled, err := jsonMarshal(payload)

			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"error": err.Error(),
				})
				return
			}
			if err := protojson.Unmarshal(marshaled, req); err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"error": "invalid request: " + err.Error(),
				})
				return
			}

			to := time.Duration(rt.Grpc.Timeoutms) * time.Millisecond
			if to <= 0 {
				to = 5 * time.Second
			}
			cctx, cancel := context.WithTimeout(ctx.Request.Context(), to)
			defer cancel()

			cc, err := c.pool.Get(ctx.Request.Context(), rt.Grpc.Service)
			if err != nil {
				ctx.JSON(http.StatusBadGateway, gin.H{
					"error": "downstream dial failed: " + err.Error(),
				})
				return
			}

			resp, err := entry.Invoke(cctx, cc, req)
			if err != nil {
				httpCode, msg := mapGrpcError(err)
				ctx.JSON(httpCode, gin.H{
					"error": msg,
				})
				return
			}

			b, err := protojson.MarshalOptions{
				UseProtoNames: true,
			}.Marshal(resp)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"error": "marshal response failed",
				})
				return
			}
			ctx.Data(http.StatusOK, "application/json", b)
		}

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

	for _, rt := range routes {
		key := normalizeTarget(rt.Grpc.Target)
		if _, ok := reg.Get(key); !ok {
			log.Fatalf("route %q refers to unknown grpc target %q", rt.Name, key)
		}
	}
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
