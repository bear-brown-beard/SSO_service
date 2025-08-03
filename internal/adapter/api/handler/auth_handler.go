package api

import (
	"encoding/json"
	"net/http"
	dto "sso/internal/adapter/api/handler/dto"
	"sso/internal/logger"
	"sso/internal/service"
	response "sso/pkg/response"
	"strings"
)

type VerificationHandler struct {
	SSOService service.SSOService
	Logger     *logger.Logger
}

func NewVerificationHandler(s service.SSOService, logger *logger.Logger) *VerificationHandler {
	return &VerificationHandler{
		SSOService: s,
		Logger:     logger,
	}
}
func (h *VerificationHandler) Verification(w http.ResponseWriter, r *http.Request) {
	var req dto.VerificationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Return(w, http.StatusOK, false, "Invalid request format", nil)
		return
	}

	if req.Platform == "" {
		req.Platform = r.Header.Get("platform")
		if req.Platform == "" {
			req.Platform = "web"
		}
	}

	err := h.SSOService.Verification(req.Phone, req.Signature, req.Platform)
	if err != nil {
		if strings.Contains(err.Error(), "invalid phone number") {
			response.Return(w, http.StatusOK, false, err.Error(), nil)
			return
		}
		if strings.Contains(err.Error(), "failed to save verification code") {
			response.Return(w, http.StatusInternalServerError, false, "Internal server error", nil)
			return
		}
		response.Return(w, http.StatusOK, false, "Failed to send verification code", nil)
		return
	}

	response.Return(w, http.StatusOK, true, "Verification code sent successfully", nil)
}

func (h *VerificationHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.Logger.Error("Error decoding request body", "error", err)
		response.Return(w, http.StatusOK, false, "Invalid request format", nil)
		return
	}

	deviceUUID := r.Header.Get("X-DeviceUUID")
	platform := r.Header.Get("platform")
	if platform == "" {
		platform = "web"
	}

	ip := r.RemoteAddr
	if forwardedFor := r.Header.Get("X-Forwarded-For"); forwardedFor != "" {
		ip = strings.Split(forwardedFor, ",")[0]
	}

	agent := r.UserAgent()
	token, err := h.SSOService.Login(req.Phone, req.Code, platform, deviceUUID, agent, ip, req.WebsiteID)
	if err != nil {
		h.Logger.Error("Error from SSOService.Login", "error", err)
		if strings.Contains(err.Error(), "invalid phone number") ||
			strings.Contains(err.Error(), "invalid verification code") ||
			strings.Contains(err.Error(), "verification code expired") {
			response.Return(w, http.StatusOK, false, err.Error(), nil)
			return
		}
		if strings.Contains(err.Error(), "failed to save token") {
			response.Return(w, http.StatusInternalServerError, false, "Internal server error", nil)
			return
		}
		response.Return(w, http.StatusOK, false, "Login failed", nil)
		return
	}

	response.Return(w, http.StatusOK, true, "Login successful", map[string]string{
		"token": token,
	})
}

func (h *VerificationHandler) Logout(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		h.Logger.Error("Logout failed: Authorization header is missing")
		response.Return(w, http.StatusUnauthorized, false, "Authorization header is required", nil)
		return
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		h.Logger.Error("Logout failed: Invalid authorization header format")
		response.Return(w, http.StatusUnauthorized, false, "Invalid authorization header format", nil)
		return
	}

	token := parts[1]

	err := h.SSOService.Logout(token)
	if err != nil {
		h.Logger.Error("Logout error", "error", err)
		response.Return(w, http.StatusUnauthorized, false, err.Error(), nil)
		return
	}

	response.Return(w, http.StatusOK, true, "Successfully logged out", nil)
}
