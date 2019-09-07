package server

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httputil"
	"strconv"
	"time"

	"github.com/blendle/zapdriver"
	"github.com/go-chi/chi/middleware"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

const grpcGatewayIdentifier = "grpcgateway" // Used to identify reuqests from the grpc gateway

func loggerHTTPMiddlewareDefault(logRequestBody bool, disabledEndpoints []string) func(http.Handler) http.Handler {
	// Make a map lookup for disabled endpoints
	disabled := make(map[string]struct{})
	for _, d := range disabledEndpoints {
		disabled[d] = struct{}{}
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// If Disabled
			if _, ok := disabled[r.RequestURI]; ok {
				next.ServeHTTP(w, r)
				return
			}

			start := time.Now()

			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			var response *bytes.Buffer
			if logRequestBody {
				response = new(bytes.Buffer)
				ww.Tee(response)
			}

			next.ServeHTTP(ww, r)

			fields := []zapcore.Field{
				zap.Int("status", ww.Status()),
				zap.Duration("duration", time.Since(start)),
				zap.String("path", r.RequestURI),
				zap.String("method", r.Method),
				zap.String("package", "server.http"),
			}

			if reqID := r.Context().Value(middleware.RequestIDKey); reqID != nil {
				fields = append(fields, zap.String("request-id", reqID.(string)))
			}

			if logRequestBody {
				if req, err := httputil.DumpRequest(r, true); err == nil {
					fields = append(fields, zap.ByteString("request", req))
				}
				fields = append(fields, zap.ByteString("response", response.Bytes()))
			}

			// If we have an x-Forwarded-For header, use that for the remote
			if forwardedFor := r.Header.Get("X-Forwarded-For"); forwardedFor != "" {
				fields = append(fields, zap.String("remote", forwardedFor))
			} else {
				fields = append(fields, zap.String("remote", r.RemoteAddr))
			}
			zap.L().Info("HTTP Request", fields...)
		})
	}
}

// Returns a middleware function for logging requests
func loggerHTTPMiddlewareStackdriver(logRequestBody bool, disabledEndpoints []string) func(http.Handler) http.Handler {
	// Make a map lookup for disabled endpoints
	disabled := make(map[string]struct{})
	for _, d := range disabledEndpoints {
		disabled[d] = struct{}{}
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// If Disabled
			if _, ok := disabled[r.RequestURI]; ok {
				next.ServeHTTP(w, r)
				return
			}

			start := time.Now()

			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			var response *bytes.Buffer
			if logRequestBody {
				response = new(bytes.Buffer)
				ww.Tee(response)
			}

			next.ServeHTTP(ww, r)

			// If the remote IP is being proxied, use the real IP
			remoteIP := r.Header.Get("X-Forwarded-For")
			if remoteIP == "" {
				remoteIP = r.RemoteAddr
			}

			fields := []zapcore.Field{
				zapdriver.HTTP(&zapdriver.HTTPPayload{
					RequestMethod: r.Method,
					RequestURL:    r.RequestURI,
					RequestSize:   strconv.FormatInt(r.ContentLength, 10),
					Status:        ww.Status(),
					ResponseSize:  strconv.Itoa(ww.BytesWritten()),
					UserAgent:     r.UserAgent(),
					RemoteIP:      remoteIP,
					Referer:       r.Referer(),
					Latency:       fmt.Sprintf("%fs", time.Since(start).Seconds()),
					Protocol:      r.Proto,
				}),
				zap.String("package", "server.http"),
			}

			if reqID := r.Context().Value(middleware.RequestIDKey); reqID != nil {
				fields = append(fields, zap.String("request-id", reqID.(string)))
			}

			if logRequestBody {
				if req, err := httputil.DumpRequest(r, true); err == nil {
					fields = append(fields, zap.ByteString("request", req))
				}
				fields = append(fields, zap.ByteString("response", response.Bytes()))
			}

			zap.L().Info("HTTP Request", fields...)
		})
	}
}

func loggerGRPCUnaryDefault(logRequestBody bool, disabledEndpoints []string) grpc.UnaryServerInterceptor {
	// Make a map lookup for disabled endpoints
	disabled := make(map[string]struct{})
	for _, d := range disabledEndpoints {
		disabled[d] = struct{}{}
	}
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// If disabled or from grpc gateway, don't log
		if _, ok := disabled[info.FullMethod]; ok || fromGRPCGateway(ctx) {
			return handler(ctx, req)
		}

		start := time.Now()
		resp, err := handler(ctx, req)
		fields := []zapcore.Field{
			zap.String("path", info.FullMethod),
			zap.String("method", "GRPC"),
			zap.String("package", "server.grpc"),
			zap.Duration("duration", time.Since(start)),
			zap.String("status", status.Code(err).String()),
		}

		if remote := grpcMetadataGetFirst(ctx, "x-forwarded-for"); remote != "" {
			fields = append(fields, zap.String("remote", remote))
		} else if p, ok := peer.FromContext(ctx); ok {
			fields = append(fields, zap.String("remote", p.Addr.String()))
		}

		if err != nil {
			fields = append(fields, zap.String("error", err.Error()))
		}

		if logRequestBody {
			fields = append(fields, zap.Any("request", req), zap.Any("response", resp))
		}

		zap.L().Info("GRPC Request", fields...)
		return resp, err
	}
}

func loggerGRPCStreamDefault(disabledEndpoints []string) grpc.StreamServerInterceptor {
	// Make a map lookup for disabled endpoints
	disabled := make(map[string]struct{})
	for _, d := range disabledEndpoints {
		disabled[d] = struct{}{}
	}
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		// If disabled or from grpc gateway, don't log
		if _, ok := disabled[info.FullMethod]; ok || fromGRPCGateway(ss.Context()) {
			return handler(srv, ss)
		}

		start := time.Now()
		fields := []zapcore.Field{
			zap.String("path", info.FullMethod),
			zap.String("method", "GRPC Stream"),
			zap.String("package", "server.grpc.stream"),
		}

		if remote := grpcMetadataGetFirst(ss.Context(), "x-forwarded-for"); remote != "" {
			fields = append(fields, zap.String("remote", remote))
		} else if p, ok := peer.FromContext(ss.Context()); ok {
			fields = append(fields, zap.String("remote", p.Addr.String()))
		}

		zap.L().Info("GRPC Stream Start", fields...)

		// Call the stream
		err := handler(srv, ss)

		fields = append(fields,
			zap.Duration("duration", time.Since(start)),
			zap.String("status", status.Code(err).String()),
		)

		if err != nil {
			fields = append(fields, zap.String("error", err.Error()))
		}

		zap.L().Info("GRPC Stream Complete", fields...)
		return err
	}
}

func loggerGRPCUnaryStackdriver(logRequestBody bool, disabledEndpoints []string) grpc.UnaryServerInterceptor {
	// Make a map lookup for disabled endpoints
	disabled := make(map[string]struct{})
	for _, d := range disabledEndpoints {
		disabled[d] = struct{}{}
	}
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// If disabled or from grpc gateway, don't log
		if _, ok := disabled[info.FullMethod]; ok || fromGRPCGateway(ctx) {
			return handler(ctx, req)
		}

		start := time.Now()

		// Call the handler
		resp, err := handler(ctx, req)

		httpPayload := &zapdriver.HTTPPayload{
			RequestMethod: "GRPC Request",
			RequestURL:    info.FullMethod,
			Protocol:      "GRPC",
			Latency:       fmt.Sprintf("%fs", time.Since(start).Seconds()),
		}

		if remote := grpcMetadataGetFirst(ctx, "x-forwarded-for"); remote != "" {
			httpPayload.RemoteIP = remote
		} else if p, ok := peer.FromContext(ctx); ok {
			httpPayload.RemoteIP = p.Addr.String()
		}

		// Build the Stackdriver HTTPPayload
		fields := []zapcore.Field{
			zapdriver.HTTP(httpPayload),
			zap.String("package", "server.grpc"),
			zap.String("status", status.Code(err).String()),
		}

		if err != nil {
			fields = append(fields, zap.String("error", err.Error()))
		}

		if logRequestBody {
			fields = append(fields, zap.Any("request", req), zap.Any("response", resp))
		}

		zap.L().Info("GRPC Request", fields...)
		return resp, err
	}
}

func loggerGRPCStreamStackdriver(disabledEndpoints []string) grpc.StreamServerInterceptor {
	// Make a map lookup for disabled endpoints
	disabled := make(map[string]struct{})
	for _, d := range disabledEndpoints {
		disabled[d] = struct{}{}
	}
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		// If disabled or from grpc gateway, don't log
		if _, ok := disabled[info.FullMethod]; ok || fromGRPCGateway(ss.Context()) {
			return handler(srv, ss)
		}

		start := time.Now()

		httpPayload := &zapdriver.HTTPPayload{
			RequestMethod: "GRPC Stream Start",
			RequestURL:    info.FullMethod,
			Protocol:      "GRPC",
		}
		if remote := grpcMetadataGetFirst(ss.Context(), "x-forwarded-for"); remote != "" {
			httpPayload.RemoteIP = remote
		} else if p, ok := peer.FromContext(ss.Context()); ok {
			httpPayload.RemoteIP = p.Addr.String()
		}

		// Build the Stackdriver HTTPPayload
		fields := []zapcore.Field{
			zapdriver.HTTP(httpPayload),
			zap.String("package", "server.grpc.stream"),
		}

		zap.L().Info("GRPC Stream Start", fields...)

		// Call the stream
		err := handler(srv, ss)

		httpPayload.RequestMethod = "GRPC Stream Complete"
		httpPayload.Latency = fmt.Sprintf("%fs", time.Since(start).Seconds())

		fields = []zapcore.Field{
			zapdriver.HTTP(httpPayload),
			zap.String("package", "server.grpc.stream"),
			zap.String("status", status.Code(err).String()),
		}
		if err != nil {
			fields = append(fields, zap.String("error", err.Error()))
		}
		zap.L().Info("GRPC Stream Complete", fields...)
		return err
	}
}

func grpcMetadataGetFirst(ctx context.Context, key string) string {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		headers := md.Get(key)
		if len(headers) > 0 {
			return headers[0]
		}
	}
	return ""
}

func fromGRPCGateway(ctx context.Context) bool {
	if grpcMetadataGetFirst(ctx, grpcGatewayIdentifier) != "" {
		return true
	}
	return false
}
