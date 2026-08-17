package log

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"mime/multipart"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/rl404/fairy/errors/stack"
)

// APIMiddlewareConfig is log config for middleware.
type APIMiddlewareConfig struct {
	// Show request header.
	RequestHeader bool
	// Show request body.
	RequestBody bool
	// Show response header.
	ResponseHeader bool
	// Show response body.
	ResponseBody bool
	// Show raw path (includes query params).
	RawPath bool
	// Show error stack.
	Error bool
}

// HTTPMiddlewareWithLog is http middleware that will log the request and response.
func HTTPMiddlewareWithLog(logger Logger, middlewareConfig ...APIMiddlewareConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return HTTPHandlerWithLog(logger, next, middlewareConfig...)
	}
}

// HTTPHandlerFuncWithLog is http handler func with log.
func HTTPHandlerFuncWithLog(logger Logger, next http.HandlerFunc, middlewareConfig ...APIMiddlewareConfig) http.HandlerFunc {
	return HTTPHandlerWithLog(logger, next, middlewareConfig...).(http.HandlerFunc)
}

// HTTPHandlerWithLog is http handler with log.
// Also includes error stack tracing feature
// if you use it.
func HTTPHandlerWithLog(logger Logger, next http.Handler, middlewareConfig ...APIMiddlewareConfig) http.Handler {
	var cfg APIMiddlewareConfig
	if len(middlewareConfig) > 0 {
		cfg = middlewareConfig[0]
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if logger == nil {
			next.ServeHTTP(w, r)
			return
		}

		// Prepare error stack tracing.
		ctx := stack.Init(r.Context())
		start := time.Now()

		// Prepare response writer for logging.
		rw := responseWriter{
			ResponseWriter: w,
		}

		// Prepare map for logging.
		ctx = initMultipartRequestBody(ctx)
		m := logHTTPRequest(r.WithContext(ctx), cfg)

		// Call next handler.
		next.ServeHTTP(&rw, r.WithContext(ctx))

		// Prepare map for logging.
		logHTTPResponse(ctx, start, &rw, cfg, m)

		logger.Log(m)
	})
}

func logHTTPRequest(r *http.Request, cfg APIMiddlewareConfig) map[string]interface{} {
	m := map[string]interface{}{
		"method": r.Method,
		"ip":     getIP(r),
		"path":   getRoutePattern(r),
	}

	if cfg.RawPath {
		m["raw_path"] = r.RequestURI
	}

	if cfg.RequestHeader {
		m["request_header"] = r.Header
	}

	if cfg.RequestBody {
		contentType := strings.ToLower(r.Header.Get("Content-Type"))
		switch {
		case strings.HasPrefix(contentType, "application/json"):
			logHTTPRequestJSONBody(r, m)
		case strings.HasPrefix(contentType, "application/x-www-form-urlencoded"):
			logHTTPRequestFormBody(r, m)
		case strings.HasPrefix(contentType, "multipart/form-data"):
			logHTTPRequestMultipartBody(r, m)
		default:
			// Ignore other content type.
		}
	}

	return m
}

func logHTTPRequestJSONBody(r *http.Request, m map[string]interface{}) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return
	}

	// Restore body so the actual handler can read it.
	r.Body = io.NopCloser(bytes.NewReader(body))

	var bodyMap interface{}
	if err := json.Unmarshal(body, &bodyMap); err != nil {
		m["request_body"] = string(body)
		return
	}

	m["request_body"] = bodyMap
}

func logHTTPRequestFormBody(r *http.Request, m map[string]interface{}) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return
	}

	// Restore body for the actual handler.
	r.Body = io.NopCloser(bytes.NewReader(body))

	values, err := url.ParseQuery(string(body))
	if err != nil {
		m["request_body"] = string(body)
		return
	}

	m["request_body"] = values
}

func logHTTPRequestMultipartBody(r *http.Request, m map[string]interface{}) {
	m["request_body"] = getMultipartFromCtx(r.Context())
}

func logHTTPResponse(ctx context.Context, start time.Time, rw *responseWriter, cfg APIMiddlewareConfig, m map[string]interface{}) {
	m["level"] = getLevelFromStatus(rw.statusCode)
	m["duration"] = time.Since(start).String()
	m["code"] = rw.statusCode

	if cfg.ResponseHeader {
		m["response_header"] = rw.Header()
	}

	if cfg.ResponseBody {
		contentType := strings.ToLower(rw.Header().Get("Content-Type"))
		switch {
		case strings.HasPrefix(contentType, "application/json"):
			logHTTPResponseJSONBody(rw, m)
		case strings.HasPrefix(contentType, "text/"):
			logHTTPResponseTextBody(rw, m)
		default:
			// Ignore other content type.
		}
	}

	// Include the error stack if you use it.
	errStack := stack.Get(ctx)
	if !cfg.Error || len(errStack) == 0 {
		return
	}

	// Copy slice to prevent reversed multiple times
	// if using multiple middleware.
	errTmp := cpSlice(errStack)

	// Reverse the stack order.
	for i, j := 0, len(errTmp)-1; i < j; i, j = i+1, j-1 {
		errTmp[i], errTmp[j] = errTmp[j], errTmp[i]
	}

	m["error"] = errTmp
}

func logHTTPResponseJSONBody(rw *responseWriter, m map[string]interface{}) {
	body := rw.body.Bytes()

	var bodyMap interface{}
	if err := json.Unmarshal(body, &bodyMap); err != nil {
		m["response_body"] = string(body)
		return
	}

	m["response_body"] = bodyMap
}

func logHTTPResponseTextBody(rw *responseWriter, m map[string]interface{}) {
	m["response_body"] = rw.body.String()
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
	body       bytes.Buffer
}

func (w *responseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *responseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func cpSlice(arr []stack.ErrStack) []stack.ErrStack {
	a := make([]stack.ErrStack, len(arr))
	copy(a, arr)
	return a
}

func getRoutePattern(r *http.Request) string {
	routePath := r.URL.Path
	if r.URL.RawPath != "" {
		routePath = r.URL.RawPath
	}

	rctx := chi.RouteContext(r.Context())
	tctx := chi.NewRouteContext()
	if rctx.Routes.Match(tctx, r.Method, routePath) {
		return tctx.RoutePattern()
	}

	return routePath
}

func getLevelFromStatus(status int) logLevel {
	switch status {
	case
		// 2xx
		http.StatusOK,
		http.StatusCreated,
		http.StatusAccepted,
		http.StatusNonAuthoritativeInfo,
		http.StatusNoContent,
		http.StatusResetContent,
		http.StatusPartialContent,
		http.StatusMultiStatus,
		http.StatusAlreadyReported,
		http.StatusIMUsed,
		// 3xx
		http.StatusMultipleChoices,
		http.StatusMovedPermanently,
		http.StatusFound,
		http.StatusSeeOther,
		http.StatusNotModified,
		http.StatusUseProxy,
		http.StatusTemporaryRedirect,
		http.StatusPermanentRedirect:
		return infoLevel
	case
		// 4xx
		http.StatusBadRequest,
		http.StatusUnauthorized,
		http.StatusPaymentRequired,
		http.StatusForbidden,
		http.StatusNotFound,
		http.StatusMethodNotAllowed,
		http.StatusNotAcceptable,
		http.StatusProxyAuthRequired,
		http.StatusRequestTimeout,
		http.StatusConflict,
		http.StatusGone,
		http.StatusLengthRequired,
		http.StatusPreconditionFailed,
		http.StatusRequestEntityTooLarge,
		http.StatusRequestURITooLong,
		http.StatusUnsupportedMediaType,
		http.StatusRequestedRangeNotSatisfiable,
		http.StatusExpectationFailed,
		http.StatusTeapot,
		http.StatusMisdirectedRequest,
		http.StatusUnprocessableEntity,
		http.StatusLocked,
		http.StatusFailedDependency,
		http.StatusTooEarly,
		http.StatusUpgradeRequired,
		http.StatusPreconditionRequired,
		http.StatusTooManyRequests,
		http.StatusRequestHeaderFieldsTooLarge,
		http.StatusUnavailableForLegalReasons:
		return warnLevel
	default:
		return errorLevel
	}
}

func getIP(r *http.Request) string {
	if host, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
		return host
	}
	return r.RemoteAddr
}

type multipartRequestBodyCtx struct{}
type multipartRequestBody map[string][]string

func initMultipartRequestBody(ctx context.Context) context.Context {
	return context.WithValue(ctx, multipartRequestBodyCtx{}, make(multipartRequestBody))
}

// AddMultipartToCtx to add multipart data to ctx for logging.
func AddMultipartToCtx(ctx context.Context, part *multipart.Part) context.Context {
	key, value := part.FormName(), part.FileName()
	if value == "" {
		var buf bytes.Buffer
		if _, err := buf.ReadFrom(part); err != nil {
			value = "<error reading value>"
		} else {
			value = buf.String()
		}
	}

	reqBody, ok := ctx.Value(multipartRequestBodyCtx{}).(multipartRequestBody)
	if !ok {
		// Create new data.
		return context.WithValue(ctx, multipartRequestBodyCtx{}, multipartRequestBody{
			key: {value},
		})
	}

	// Add to existing one.
	reqBody[key] = append(reqBody[key], value)
	return context.WithValue(ctx, multipartRequestBodyCtx{}, reqBody)
}

func getMultipartFromCtx(ctx context.Context) multipartRequestBody {
	reqBody, ok := ctx.Value(multipartRequestBodyCtx{}).(multipartRequestBody)
	if ok {
		return reqBody
	}
	return nil
}
