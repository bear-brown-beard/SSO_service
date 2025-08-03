package service

import (
	"context"
	"fmt"
	"math/rand/v2"
	"regexp"
	"time"

	"sso/internal/logger"
	"sso/internal/repository"
)

type SSOService interface {
	Verification(phone, signature, platform string) error
	Login(phone, code, platform, deviceUUID, agent, ip, websiteID string) (string, error)
	Logout(token string) error
}

type SSOAuthService struct {
	TestAccountRepo repository.TestAccountRepository
	UserRepo        repository.UserRepository
	TokenRepo       repository.TokenRepository
	CodeCache       CacheService
	JWTService      JWTService
	SMSService      SMSCService
	MindboxService  AuthMindboxService
	Logger          *logger.Logger
}

func validatePhone(phone string) (string, error) {
	phone = regexp.MustCompile(`[^\d]`).ReplaceAllString(phone, "")
	if len(phone) != 11 {
		return phone, fmt.Errorf("invalid phone number: phone number must be 11 digits")
	}
	if phone[0] != '7' && phone[0] != '8' {
		return phone, fmt.Errorf("invalid phone number: phone number must start with 7 or 8")
	}
	if phone[0] == '8' {
		phone = "7" + phone[1:]
	}
	return phone, nil
}

func generateCode() string {
	return fmt.Sprintf("%04d", rand.IntN(9000)+1000)
}

func (s *SSOAuthService) Verification(phone, signature, platform string) error {
	normalizedPhone, err := validatePhone(phone)
	if err != nil {
		return err
	}

	testAccount, err := s.TestAccountRepo.FindByPhone(normalizedPhone)
	if err != nil {
		return fmt.Errorf("error checking test account: %w", err)
	}

	var code string
	if testAccount != nil {
		if !testAccount.Code.Valid {
			return fmt.Errorf("test code is not set for test account")
		}
		code = testAccount.Code.String
		s.Logger.Info("Test verification code generated", "phone", normalizedPhone, "code", code)
	} else {
		code = generateCode()

		go func() {
			err := s.SMSService.SendVerificationCode(normalizedPhone, code, signature, platform)
			if err != nil {
				s.Logger.Error("Async SMS sending failed",
					"phone", normalizedPhone,
					"code", code,
					"error", err)
			} else {
				s.Logger.Info("Async SMS sent successfully",
					"phone", normalizedPhone,
					"code", code)
			}
		}()

		s.Logger.Info("Verification code generated and sent", "phone", normalizedPhone, "code", code)
	}

	ttl := 5 * time.Minute
	err = s.CodeCache.SaveCode(context.Background(), normalizedPhone, code, ttl)
	if err != nil {
		return fmt.Errorf("error saving verification code to cache: %w", err)
	}

	return nil
}

func (s *SSOAuthService) Login(phone, code, platform, deviceUUID, agent, ip, websiteID string) (string, error) {
	normalizedPhone, err := validatePhone(phone)
	if err != nil {
		return "", err
	}

	storedCode, err := s.CodeCache.GetCode(context.Background(), normalizedPhone)
	if err != nil {
		return "", fmt.Errorf("error retrieving verification code from cache: %w", err)
	}
	if storedCode == "" {
		return "", fmt.Errorf("verification code expired or not found")
	}
	if storedCode != code {
		return "", fmt.Errorf("invalid verification code")
	}

	err = s.CodeCache.DeleteCode(context.Background(), normalizedPhone)
	if err != nil {
		s.Logger.Warn("Failed to delete verification code", "error", err)
	}

	user, err := s.UserRepo.FindByPhone(normalizedPhone)
	if err != nil {
		return "", fmt.Errorf("error finding user in repository: %w", err)
	}

	if user == nil {
		user, err = s.UserRepo.Create(normalizedPhone)
		if err != nil {
			return "", fmt.Errorf("error creating user in repository: %w", err)
		}
		mindboxWebsiteID := websiteID
		if mindboxWebsiteID == "" {
			mindboxWebsiteID = fmt.Sprintf("%d", user.ID)
		}

		go func() {
			err := s.MindboxService.RegisterUser(
				user,
				platform,
				mindboxWebsiteID,
				deviceUUID,
				agent,
			)
			if err != nil {
				s.Logger.Error("Async Mindbox registration failed",
					"user_id", user.ID,
					"phone", user.Phone.String,
					"error", err)
			} else {
				s.Logger.Info("Async Mindbox registration successful",
					"user_id", user.ID,
					"phone", user.Phone.String)
			}
		}()
	} else {
		mindboxWebsiteID := websiteID
		if mindboxWebsiteID == "" {
			mindboxWebsiteID = fmt.Sprintf("%d", user.ID)
		}

		go func() {
			err := s.MindboxService.LoginUser(
				user,
				platform,
				mindboxWebsiteID,
				"",
				deviceUUID,
				agent,
			)
			if err != nil {
				s.Logger.Error("Async Mindbox login failed",
					"user_id", user.ID,
					"phone", user.Phone.String,
					"error", err)
			} else {
				s.Logger.Info("Async Mindbox login successful",
					"user_id", user.ID,
					"phone", user.Phone.String)
			}
		}()
	}

	token, err := s.JWTService.GenerateToken(user.ID, normalizedPhone)
	if err != nil {
		return "", fmt.Errorf("error generating JWT token: %w", err)
	}

	return token, nil
}

func (s *SSOAuthService) Logout(token string) error {
	_, err := s.JWTService.ValidateToken(token)
	if err != nil {
		return fmt.Errorf("invalid token: %w", err)
	}

	err = s.CodeCache.AddToBlacklist(context.Background(), token, 24*time.Hour)
	if err != nil {
		return fmt.Errorf("failed to add token to blacklist: %w", err)
	}

	err = s.TokenRepo.Deactivate(token)
	if err != nil {
		s.Logger.Warn("Failed to deactivate token in database", "error", err)
	}

	return nil
}
