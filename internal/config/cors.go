package config

import (
	"strings"
)

func ParseCORSOrigins(origins string) []string {
	if origins == "" {
		return []string{}
	}
	return strings.Split(origins, ",")
}

func ParseCORSMethods(methods string) []string {
	if methods == "" {
		return []string{}
	}
	return strings.Split(methods, ",")
}

func ParseCORSHeaders(headers string) []string {
	if headers == "" {
		return []string{}
	}
	return strings.Split(headers, ",")
}

func GetCORSConfig(cfg *Config) CORSConfig {
	switch cfg.AppEnv {
	case "local", "dev":
		return CORSConfig{
			AllowedOrigins:   []string{"http://localhost:3000", "http://localhost:8080", "http://localhost:5173"},
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-DeviceUUID", "X-Platform"},
			ExposedHeaders:   []string{"Link"},
			AllowCredentials: true,
			MaxAge:           300,
		}
	case "stage":
		return CORSConfig{
			AllowedOrigins:   []string{"https://stage.yourdomain.com"},
			AllowedMethods:   []string{"POST", "OPTIONS"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-DeviceUUID", "X-Platform"},
			ExposedHeaders:   []string{"Link"},
			AllowCredentials: true,
			MaxAge:           300,
		}
	case "prod":
		return CORSConfig{
			AllowedOrigins:   []string{"https://yourdomain.com"},
			AllowedMethods:   []string{"POST", "OPTIONS"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-DeviceUUID", "X-Platform"},
			ExposedHeaders:   []string{"Link"},
			AllowCredentials: true,
			MaxAge:           300,
		}
	default:
		return CORSConfig{
			AllowedOrigins:   []string{"*"},
			AllowedMethods:   []string{"*"},
			AllowedHeaders:   []string{"*"},
			ExposedHeaders:   []string{},
			AllowCredentials: false,
			MaxAge:           300,
		}
	}
}
