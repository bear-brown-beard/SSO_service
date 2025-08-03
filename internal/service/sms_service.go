package service

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sso/internal/logger"
	"sso/internal/models"
)

type SMSCService interface {
	SendVerificationCode(phone, code, signature, platform string) error
}

type smscService struct {
	login    string
	password string
	apiURL   string
	logger   *logger.Logger
}

func NewSMSCService(login, password string, logger *logger.Logger) SMSCService {
	return &smscService{
		login:    login,
		password: password,
		apiURL:   "https://smsc.kz/sys/send.php",
		logger:   logger,
	}
}

func (s *smscService) SendVerificationCode(phone, code, signature, platform string) error {
	messageText := fmt.Sprintf("%s код доступа для авторизации", code)

	if platform == "android" && signature != "" {
		messageText = fmt.Sprintf("%s\n%s", messageText, signature)
	}

	params := url.Values{}
	params.Set("login", s.login)
	params.Set("psw", s.password)
	params.Set("phones", phone)
	params.Set("mes", messageText)
	params.Set("fmt", "3")

	requestURL := fmt.Sprintf("%s?%s", s.apiURL, params.Encode())

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	resp, err := client.Get(requestURL)
	if err != nil {
		s.logger.Error("Failed to send HTTP request", "error", err)
		return fmt.Errorf("failed to send SMS request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	var smscResp models.SMSCResponse
	if err := json.Unmarshal(body, &smscResp); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	if smscResp.Status != 0 {
		return fmt.Errorf("SMS service error: %s", smscResp.Error)
	}

	return nil
}
