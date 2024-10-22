package application

import (
	"fmt"
	"net/http"
)

func (a *App) addHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set(
			"Content-Security-Policy",
			"default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com",
		)
		writer.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		writer.Header().Set("X-Content-Type-Options", "nosniff")
		writer.Header().Set("X-Frame-Options", "deny")
		writer.Header().Set("X-XSS-Protection", "0")

		writer.Header().Set("Server", "Go")

		next.ServeHTTP(writer, request)
	})
}

func (a *App) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		a.logger.Info("request received",
			"ip", request.RemoteAddr,
			"proto", request.Proto,
			"method", request.Method,
			"uri", request.RequestURI,
		)
		next.ServeHTTP(writer, request)
	})
}

func (a *App) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				writer.Header().Set("Connection", "close")
				a.internalServerError(writer, request, fmt.Errorf("panic recovered: %v", r))
			}
		}()
		next.ServeHTTP(writer, request)
	})
}
