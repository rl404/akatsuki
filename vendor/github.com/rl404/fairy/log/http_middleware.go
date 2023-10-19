package log

import (
	"bytes"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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

		var bw bodyWriter
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		ww.Tee(&bw)

		// Get request body.
		var body []byte
		if r.Body != nil {
			b, err := io.ReadAll(r.Body)
			if err != nil {
				_ = stack.Wrap(ctx, err)
			}
			body, r.Body = b, io.NopCloser(bytes.NewBuffer(b))
		}

		// Call next handler.
		next.ServeHTTP(ww, r.WithContext(ctx))

		// Prepare map for logging.
		m := map[string]interface{}{
			"level":    getLevelFromStatus(ww.Status()),
			"duration": time.Since(start).String(),
			"method":   r.Method,
			"path":     r.URL.Path,
			"code":     ww.Status(),
			"ip":       getIP(r),
		}

		if path, ok := getRoutePattern(r); ok {
			m["path"] = path
		}

		if cfg.RawPath {
			m["raw_path"] = r.RequestURI
		}

		if cfg.RequestHeader {
			m["request_header"] = r.Header
		}

		if cfg.RequestBody {
			m["request_body"] = string(body)
		}

		if cfg.ResponseHeader {
			m["response_header"] = ww.Header()
		}

		if cfg.ResponseBody {
			m["response_body"] = string(bw.Body)
		}

		// Include the error stack if you use it.
		errStack := stack.Get(ctx)
		if cfg.Error && len(errStack) > 0 {
			// Copy slice to prevent reversed multiple times
			// if using multiple middleware.
			errTmp := cpSlice(errStack)

			// Reverse the stack order.
			for i, j := 0, len(errTmp)-1; i < j; i, j = i+1, j-1 {
				errTmp[i], errTmp[j] = errTmp[j], errTmp[i]
			}

			m["error"] = errTmp
		}

		logger.Log(m)
	})
}

func cpSlice(arr []stack.ErrStack) []stack.ErrStack {
	a := make([]stack.ErrStack, len(arr))
	copy(a, arr)
	return a
}

func getRoutePattern(r *http.Request) (string, bool) {
	routePath := r.URL.Path
	if r.URL.RawPath != "" {
		routePath = r.URL.RawPath
	}

	rctx := chi.RouteContext(r.Context())
	tctx := chi.NewRouteContext()
	if rctx.Routes.Match(tctx, r.Method, routePath) {
		return tctx.RoutePattern(), true
	}

	return "", false
}

func getLevelFromStatus(status int) logLevel {
	switch status {
	case
		http.StatusOK,
		http.StatusCreated,
		http.StatusAccepted,
		http.StatusMultipleChoices,
		http.StatusMovedPermanently,
		http.StatusFound,
		http.StatusSeeOther,
		http.StatusNotModified,
		http.StatusTemporaryRedirect,
		http.StatusPermanentRedirect:
		return infoLevel
	case
		http.StatusBadRequest,
		http.StatusUnauthorized,
		http.StatusForbidden,
		http.StatusNotFound,
		http.StatusMethodNotAllowed,
		http.StatusNotAcceptable,
		http.StatusRequestTimeout,
		http.StatusConflict,
		http.StatusGone,
		http.StatusPreconditionFailed,
		http.StatusExpectationFailed,
		http.StatusMisdirectedRequest,
		http.StatusUnprocessableEntity,
		http.StatusFailedDependency,
		http.StatusTooManyRequests:
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

type bodyWriter struct {
	Body []byte
}

func (w *bodyWriter) Write(b []byte) (int, error) {
	w.Body = b
	return len(b), nil
}
