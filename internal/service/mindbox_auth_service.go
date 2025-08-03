package service

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sso/internal/config"
	"sso/internal/logger"
	"sso/internal/models"
	"sso/internal/repository"
	"time"
)

type AuthMindboxService interface {
	RegisterUser(user *models.User, platform string, websiteID string, deviceUUID string, userAgent string) error
	LoginUser(user *models.User, platform string, websiteID string, mindboxID string, deviceUUID string, userAgent string) error
}

type authMindboxService struct {
	userRepo        repository.UserRepository
	userMindBoxRepo repository.UserMindBoxRepository
	log             *logger.Logger
	cfg             *config.Config
	client          *http.Client
}

func createHTTPClient(cfg *config.Config) *http.Client {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: false,
		},
	}

	return &http.Client{
		Transport: transport,
		Timeout:   cfg.TLS.Timeout,
	}
}

func NewAuthMindboxService(log *logger.Logger, userRepo repository.UserRepository, userMindBoxRepo repository.UserMindBoxRepository, cfg *config.Config) AuthMindboxService {
	client := createHTTPClient(cfg)

	log.Debug("Mindbox HTTP client configured",
		"tls_skip_verify", false,
		"timeout", cfg.TLS.Timeout,
		"environment", cfg.AppEnv)

	return &authMindboxService{
		userRepo:        userRepo,
		userMindBoxRepo: userMindBoxRepo,
		log:             log,
		cfg:             cfg,
		client:          client,
	}
}

func (s *authMindboxService) RegisterUser(user *models.User, platform string, websiteID string, deviceUUID string, userAgent string) error {
	userMindBox, err := s.userMindBoxRepo.FindByUserID(user.ID)
	if err != nil {
		s.log.Warn("Failed to check user in user_mind_box", "user_id", user.ID, "error", err)
	} else if userMindBox == nil {
		var createErr error
		_, createErr = s.userMindBoxRepo.Create(user.ID)
		if createErr != nil {
			s.log.Warn("Failed to create user in user_mind_box", "user_id", user.ID, "error", createErr)
		} else {
			s.log.Debug("Created user record in user_mind_box", "user_id", user.ID)
		}
	}

	customerData := map[string]any{
		"mobilePhone": user.Phone.String,
		"ids": map[string]string{
			"websiteID": websiteID,
		},
		"subscriptions": []map[string]any{
			{"pointOfContact": "Email", "isSubscribed": 1},
			{"pointOfContact": "SMS", "isSubscribed": 1},
			{"pointOfContact": "Webpush", "isSubscribed": 1},
		},
	}

	if user.FirstName.Valid && user.FirstName.String != "" {
		customerData["firstName"] = user.FirstName.String
	}
	if user.LastName.Valid && user.LastName.String != "" {
		customerData["lastName"] = user.LastName.String
	}
	if user.Email.Valid && user.Email.String != "" {
		customerData["email"] = user.Email.String
	}

	requestBody := map[string]any{
		"customer":             customerData,
		"executionDateTimeUtc": time.Now().UTC().Format("2006-01-02 15:04:05.000"),
	}

	err = s.Send(platform, "RegisterCustomer", requestBody, deviceUUID, userAgent)
	if err != nil {
		s.log.Error("failed to send RegisterUser to mindbox", "error", err.Error())
		return fmt.Errorf("failed to register user in mindbox: %w", err)
	}

	return nil
}

func (s *authMindboxService) LoginUser(user *models.User, platform string, websiteID string, mindboxID string, deviceUUID string, userAgent string) error {
	userMindBox, err := s.userMindBoxRepo.FindByUserID(user.ID)
	if err != nil {
		s.log.Warn("Failed to check user in user_mind_box", "user_id", user.ID, "error", err)
	} else if userMindBox == nil {
		var createErr error
		_, createErr = s.userMindBoxRepo.Create(user.ID)
		if createErr != nil {
			s.log.Warn("Failed to create user in user_mind_box", "user_id", user.ID, "error", createErr)
		} else {
			s.log.Debug("Created user record in user_mind_box", "user_id", user.ID)
		}
	}

	customerData := map[string]any{
		"mobilePhone": user.Phone.String,
		"ids": map[string]string{
			"websiteID": websiteID,
		},
		"subscriptions": []map[string]any{
			{"pointOfContact": "Email", "isSubscribed": 1},
			{"pointOfContact": "SMS", "isSubscribed": 1},
			{"pointOfContact": "Webpush", "isSubscribed": 1},
		},
	}

	if mindboxID != "" {
		customerData["ids"].(map[string]string)["mindboxId"] = mindboxID
	}

	if user.FirstName.Valid && user.FirstName.String != "" {
		customerData["firstName"] = user.FirstName.String
	}
	if user.LastName.Valid && user.LastName.String != "" {
		customerData["lastName"] = user.LastName.String
	}
	if user.Email.Valid && user.Email.String != "" {
		customerData["email"] = user.Email.String
	}

	requestBody := map[string]any{
		"customer":             customerData,
		"executionDateTimeUtc": time.Now().UTC().Format("2006-01-02 15:04:05.000"),
	}

	err = s.Send(platform, "AuthorizeCustomer", requestBody, deviceUUID, userAgent)
	if err != nil {
		s.log.Error("failed to send LoginUser to mindbox", "error", err.Error())
		return fmt.Errorf("failed to login user in mindbox: %w", err)
	}

	return nil
}

func (s *authMindboxService) Send(platform, operation string, data any, deviceUUID string, userAgent string) error {
	var auth, endpointId string

	switch platform {
	case "android":
		auth = s.cfg.Mindbox.Android.Auth
		endpointId = s.cfg.Mindbox.Android.EndpointID
	case "ios":
		auth = s.cfg.Mindbox.IOS.Auth
		endpointId = s.cfg.Mindbox.IOS.EndpointID
	case "web":
		auth = s.cfg.Mindbox.Web.Auth
		endpointId = s.cfg.Mindbox.Web.EndpointID
	default:
		return errors.New("unknown platform")
	}

	url := fmt.Sprintf("%s?endpointId=%s&operation=%s.%s", s.cfg.Mindbox.Url, endpointId, s.cfg.Mindbox.OperationPrefix, operation)
	if deviceUUID != "" {
		url = fmt.Sprintf("%s&deviceUUID=%s", url, deviceUUID)
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal request data: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Authorization", auth)
	if userAgent != "" {
		req.Header.Set("User-Agent", userAgent)
	}
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		s.log.Error("Mindbox API error response",
			"status", resp.Status,
			"body", string(body))
	} else {
		s.log.Debug("Mindbox API response", "status", resp.Status)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	return nil
}
