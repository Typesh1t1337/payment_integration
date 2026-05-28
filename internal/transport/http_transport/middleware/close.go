package middleware

import "net/http"

func CloseBodyMiddleware() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r != nil && r.Body != nil {
				defer func() {
					_ = r.Body.Close()
				}()
			}
			next.ServeHTTP(w, r)
		})
	}
}
