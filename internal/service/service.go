package service

import (
	"encoding/base64"
	"go-infra/internal/config"
	"go-infra/internal/container"
	"go-infra/internal/i18n"
	"go-infra/internal/infra/repository"
)

// AppService all services ep
type AppService interface {
	Config() *config.AppConfig
	// Logger() logger.AppLogger

	UserLang(code string) i18n.UserLang
	HasLang(code string) bool

	Repository() repository.AppRepository

	SmsSender() SmsSender
	EmailSender() EmailSender
}
type appService struct {
	container   container.AppContainer
	smsSender   SmsSender
	emailSender EmailSender
}

// MustNewAppServiceProd prod
func MustNewAppServiceProd() AppService {

	appContainer := container.MustNewAppContainer()
	appConfig := appContainer.Config()
	createRepository(appContainer)

	return &appService{
		smsSender:   NewSmsSender(appConfig),
		emailSender: NewEmailSender(appConfig),
		//
		container: appContainer,
	}
}

// MustNewAppServiceTesting testing
func MustNewAppServiceTesting() AppService {
	return MustNewAppServiceProd()
}

func (x *appService) Config() *config.AppConfig { return x.container.Config() }

// func (x *appService) Logger() logger.AppLogger  { return x.container.Logger() }

func (x *appService) UserLang(code string) i18n.UserLang { return x.container.UserLang(code) }
func (x *appService) HasLang(code string) bool           { return x.container.HasLang(code) }

func (x *appService) Repository() repository.AppRepository { return x.container.Repository() }

func (x *appService) SmsSender() SmsSender     { return x.smsSender }
func (x *appService) EmailSender() EmailSender { return x.emailSender }

func BasicAuth(username, password string) string {
	// Combine username and password in the format "username:password"
	auth := username + ":" + password
	// Encode the combination into base64
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
}
