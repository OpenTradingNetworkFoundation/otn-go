package httpserver

import (
	"net/http"
	"time"

	"github.com/gorilla/handlers"
)

type CORSConfig struct {
	Origins []string `json:"origins"`
	Methods []string `json:"methods"`
	Headers []string `json:"headers"`
}

type Config struct {
	Addr string      `json:"addr"`
	CORS *CORSConfig `json:"cors"`
}

func configureCORS(cfg *CORSConfig, next http.Handler) http.Handler {
	var opts []handlers.CORSOption
	if cfg.Origins != nil && len(cfg.Origins) > 0 {
		opts = append(opts, handlers.AllowedOrigins(cfg.Origins))
	}

	if cfg.Methods != nil && len(cfg.Methods) > 0 {
		opts = append(opts, handlers.AllowedMethods(cfg.Methods))
	}

	if cfg.Headers != nil && len(cfg.Headers) > 0 {
		opts = append(opts, handlers.AllowedHeaders(cfg.Headers))
	}

	if opts != nil {
		return handlers.CORS(opts...)(next)
	}

	return next
}

func NewHTTPServer(cfg *Config, router http.Handler) *http.Server {
	if cfg.CORS != nil {
		router = configureCORS(cfg.CORS, router)
	}

	s := &http.Server{
		Addr:         cfg.Addr,
		Handler:      router,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
	}

	return s
}
