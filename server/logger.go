package server

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/blendle/zapdriver"
	"github.com/go-chi/chi/middleware"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

// Returns a middleware function for logging requests
func loggerHTTPMiddlewareStackdriver() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			var requestID string
			if reqID := r.Context().Value(middleware.RequestIDKey); reqID != nil {
				requestID = reqID.(string)
			}
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			// Parse the request
			next.ServeHTTP(ww, r)
			// Don't log the version endpoint, it's too noisy
			if r.RequestURI == "/version" {
				return
			}
			// If the remote IP is being proxied, use the real IP
			remoteIP := r.Header.Get("X-Forwarded-For")
			if remoteIP == "" {
				remoteIP = r.RemoteAddr
			}
			zap.L().Info("HTTP Request", []zapcore.Field{
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
				zap.String("request-id", requestID),
			}...)
		})
	}
}

func loggerHTTPMiddlewareDefault() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			var requestID string
			if reqID := r.Context().Value(middleware.RequestIDKey); reqID != nil {
				requestID = reqID.(string)
			}
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			next.ServeHTTP(ww, r)

			// Don't log the version endpoint, it's too noisy
			if r.RequestURI == "/version" {
				return
			}

			fields := []zapcore.Field{
				zap.Int("status", ww.Status()),
				zap.Duration("duration", time.Since(start)),
				zap.String("request", r.RequestURI),
				zap.String("method", r.Method),
				zap.String("package", "server.request"),
			}
			if requestID != "" {
				fields = append(fields, zap.String("request-id", requestID))
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

func loggerGRPCUnaryDefault() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// Don't log the version endpoint
		if info.FullMethod == "/rpc.VersionRPC/Version" {
			return handler(ctx, req)
		}
		start := time.Now()
		resp, err := handler(ctx, req)
		fields := []zapcore.Field{
			zap.String("request", info.FullMethod),
			zap.String("method", "GRPC"),
			zap.String("package", "server.grpc"),
			zap.Duration("duration", time.Since(start)),
			zap.String("status", status.Code(err).String()),
		}
		if err != nil {
			fields = append(fields, zap.String("error", err.Error()))
		}
		zap.L().Info("GRPC Request", fields...)
		return resp, err
	}
}

func loggerGRPCStreamDefault() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		start := time.Now()
		fields := []zapcore.Field{
			zap.String("request", info.FullMethod),
			zap.String("method", "GRPC Stream"),
			zap.String("package", "server.grpc.stream"),
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

func loggerGRPCUnaryStackdriver() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// Don't log the version endpoint
		if info.FullMethod == "/rpc.VersionRPC/Version" {
			return handler(ctx, req)
		}

		start := time.Now()

		// Call the handler
		resp, err := handler(ctx, req)

		// Build the Stackdriver HTTPPayload
		fields := []zapcore.Field{
			zapdriver.HTTP(&zapdriver.HTTPPayload{
				RequestMethod: "GRPC Request",
				RequestURL:    info.FullMethod,
				Protocol:      "GRPC",
				Latency:       fmt.Sprintf("%fs", time.Since(start).Seconds()),
			}),
			zap.String("package", "server.grpc"),
			zap.String("status", status.Code(err).String()),
		}
		if err != nil {
			fields = append(fields, zap.String("error", err.Error()))
		}

		zap.L().Info("GRPC Request", fields...)
		return resp, err
	}
}

func loggerGRPCStreamStackdriver() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		start := time.Now()

		// Build the Stackdriver HTTPPayload
		fields := []zapcore.Field{
			zapdriver.HTTP(&zapdriver.HTTPPayload{
				RequestMethod: "GRPC Stream Stop",
				RequestURL:    info.FullMethod,
				Protocol:      "GRPC",
			}),
			zap.String("package", "server.grpc.stream"),
		}
		zap.L().Info("GRPC Stream Start", fields...)

		// Call the stream
		err := handler(srv, ss)

		fields = []zapcore.Field{
			zapdriver.HTTP(&zapdriver.HTTPPayload{
				RequestMethod: "GRPC Stream Start",
				RequestURL:    info.FullMethod,
				Protocol:      "GRPC",
				Latency:       fmt.Sprintf("%fs", time.Since(start).Seconds()),
			}),
			zap.String("package", "server.grpc.stream"),
		}
		fields = append(fields,
			zap.String("status", status.Code(err).String()),
		)
		if err != nil {
			fields = append(fields, zap.String("error", err.Error()))
		}
		zap.L().Info("GRPC Stream Complete", fields...)
		return err
	}
}
