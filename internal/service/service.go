package service

import (
	"encoding/base64"
	"go-infra/internal/config"
	"go-infra/internal/i18n"
	"go-infra/internal/repository"
	"os"
	"time"

	xlog "go-infra/internal/util/utillog"
	"net/http"
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
type defaultAppService struct {
	smsSender   SmsSender
	emailSender EmailSender

	configSource *config.AppConfigSource
	repository   repository.AppRepository

	lang i18n.AppLang
}

func (x *defaultAppService) mustConfig() {

	d, _ := os.Getwd()

	xlog.Info("Current work dir: %v", d)

	x.configSource = config.MustNewAppConfigSource()

	appConfig := x.Config() // first call, init

	mustConfigRuntime(appConfig)

}

func (x *defaultAppService) mustBuild() {

	appConfig := x.Config()

	x.lang = i18n.MustNewAppLang(appConfig)

	x.repository = repository.MustNewRepository(appConfig) // , appLogger)

	x.smsSender = NewSmsSender(appConfig)
	x.emailSender = NewEmailSender(appConfig)

	mustCreateRepository(x)
}

func mustConfigRuntime(appConfig *config.AppConfig) {
	t, ok := http.DefaultTransport.(*http.Transport)

	if ok {
		x := appConfig.HTTPTransport

		if x.MaxIdleConns > 0 {
			xlog.Info("Http.Transport.MaxIdleConns=%v", x.MaxIdleConns)
			t.MaxIdleConns = x.MaxIdleConns
		}
		if x.IdleConnTimeout > 0 {
			xlog.Info("Http.Transport.IdleConnTimeout=%v", x.IdleConnTimeout)
			t.IdleConnTimeout = time.Duration(x.IdleConnTimeout) * time.Second
		}
		if x.MaxConnsPerHost > 0 {
			xlog.Info("Http.Transport.MaxConnsPerHost=%v", x.MaxConnsPerHost)
			t.MaxConnsPerHost = x.MaxConnsPerHost
		}

		if x.MaxIdleConnsPerHost > 0 {
			xlog.Info("Http.Transport.MaxIdleConnsPerHost=%v", x.MaxIdleConnsPerHost)
			t.MaxIdleConnsPerHost = x.MaxIdleConnsPerHost
		}

	} else {
		xlog.Error("Cannot init http.Transport")
	}
}

// MustNewAppServiceProd prod
func MustNewAppServiceProd() AppService {

	appService := &defaultAppService{}

	appService.mustConfig()
	appService.mustBuild()

	return appService

}

// MustNewAppServiceTesting testing
func MustNewAppServiceTesting() AppService {

	return MustNewAppServiceProd()
}

func (x *defaultAppService) Config() *config.AppConfig { return x.configSource.Config() }

// func (x *appService) Logger() logger.AppLogger  { return x.container.Logger() }

func (x *defaultAppService) UserLang(code string) i18n.UserLang { return x.lang.UserLang(code) }
func (x *defaultAppService) HasLang(code string) bool           { return x.lang.HasLang(code) }

func (x *defaultAppService) Repository() repository.AppRepository { return x.repository }

func (x *defaultAppService) SmsSender() SmsSender     { return x.smsSender }
func (x *defaultAppService) EmailSender() EmailSender { return x.emailSender }

func BasicAuth(username, password string) string {
	// Combine username and password in the format "username:password"
	auth := username + ":" + password
	// Encode the combination into base64
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
}
