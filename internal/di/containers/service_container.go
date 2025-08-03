package containers

import (
	"sso/internal/config"
	"sso/internal/logger"
	"sso/internal/repository"
	"sso/internal/service"
)

type ServiceContainer struct {
	ssoService service.SSOService
	jwtService service.JWTService
	logger     *logger.Logger
}

func NewSSOService(
	testRepo repository.TestAccountRepository,
	userRepo repository.UserRepository,
	tokenRepo repository.TokenRepository,
	codeCache service.CacheService,
	jwtService service.JWTService,
	smsService service.SMSCService,
	mindboxService service.AuthMindboxService,
	logger *logger.Logger,
) service.SSOService {
	return &service.SSOAuthService{
		TestAccountRepo: testRepo,
		UserRepo:        userRepo,
		TokenRepo:       tokenRepo,
		CodeCache:       codeCache,
		JWTService:      jwtService,
		SMSService:      smsService,
		MindboxService:  mindboxService,
		Logger:          logger,
	}
}

func NewServiceContainer(repoContainer *RepositoryContainer, cacheContainer *CacheContainer, cfg *config.Config, logger *logger.Logger) (*ServiceContainer, error) {
	container := &ServiceContainer{
		logger: logger,
	}

	jwtService := service.NewJWTService(cfg.JWT.SecretKey, cacheContainer.GetCodeCache())
	container.jwtService = jwtService

	ssoService := NewSSOService(
		repoContainer.TestAccountRepo,
		repoContainer.UserRepo,
		repoContainer.TokenRepo,
		cacheContainer.GetCodeCache(),
		jwtService,
		service.NewSMSCService(cfg.SMSCLogin, cfg.SMSCPassword, logger),
		service.NewAuthMindboxService(logger, repoContainer.UserRepo, repoContainer.UserMindBoxRepo, cfg),
		logger,
	)
	container.ssoService = ssoService

	logger.Debug("All services initialized successfully")
	return container, nil
}

func (c *ServiceContainer) GetSSOService() service.SSOService {
	return c.ssoService
}

func (c *ServiceContainer) GetJWTService() service.JWTService {
	return c.jwtService
}
