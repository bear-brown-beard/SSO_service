package api

import (
	"net/http"
	api "sso/internal/adapter/api/handler"
	"sso/internal/config"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func NewRouter(verificationHandler *api.VerificationHandler, cfg *config.Config) http.Handler {
	r := chi.NewRouter()

	if cfg.AppEnv == "local" {
		r.Use(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				if strings.HasPrefix(req.URL.Path, "/sso") {
					req.URL.Path = strings.TrimPrefix(req.URL.Path, "/sso")
					if req.URL.Path == "" {
						req.URL.Path = "/"
					}
				}
				next.ServeHTTP(w, req)
			})
		})
	}
	corsConfig := config.GetCORSConfig(cfg)
	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   corsConfig.AllowedOrigins,
		AllowedMethods:   corsConfig.AllowedMethods,
		AllowedHeaders:   corsConfig.AllowedHeaders,
		ExposedHeaders:   corsConfig.ExposedHeaders,
		AllowCredentials: corsConfig.AllowCredentials,
		MaxAge:           corsConfig.MaxAge,
	})

	r.Use(corsMiddleware.Handler)

	r.Post("/verification", verificationHandler.Verification)
	r.Post("/login", verificationHandler.Login)
	r.Post("/logout", verificationHandler.Logout)

	return r
}
