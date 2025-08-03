package dto

type LoginRequest struct {
	Phone     string `json:"phone"`
	Code      string `json:"code"`
	WebsiteID string `json:"websiteID,omitempty"`
}

type LoginResponse struct {
	Message string `json:"message"`
	Token   string `json:"token"`
}

type LogoutResponse struct {
	Message string `json:"message"`
}

type VerificationRequest struct {
	Phone     string `json:"phone"`
	Signature string `json:"signature"`
	Platform  string `json:"platform,omitempty"`
	WebsiteID string `json:"websiteID,omitempty"`
}

type VerificationResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
