// Package container app main container
package container

import (
	config "go-infra/internal/config"
	"net/http"
	"time"

	xlog "go-infra/internal/tool/toollog"

	i18n "go-infra/internal/i18n" // session "go-infra/session"

	repository "go-infra/internal/infra/repository"
)

// AppContainer running app container
type AppContainer interface {
	Repository() repository.AppRepository

	Config() *config.AppConfig
	// Logger() logger.AppLogger

	UserLang(code string) i18n.UserLang
	HasLang(code string) bool
}

type appContainer struct {
	// logger       logger.AppLogger
	configSource *config.AppConfigSource
	repository   repository.AppRepository

	lang i18n.AppLang
}

func (c *appContainer) Repository() repository.AppRepository {
	return c.repository
}

func (c *appContainer) Config() *config.AppConfig {
	return c.configSource.Config()
}

// func (c *appContainer) Logger() logger.AppLogger {
// 	return c.logger
// }

func (c *appContainer) UserLang(code string) i18n.UserLang {
	return c.lang.UserLang(code)
}
func (c *appContainer) HasLang(code string) bool {
	return c.lang.HasLang(code)
}
func initRuntime(appConfig *config.AppConfig) {
	t, ok := http.DefaultTransport.(*http.Transport)

	if ok {
		c := appConfig.HTTPTransport

		if c.MaxIdleConns > 0 {
			t.MaxIdleConns = c.MaxIdleConns
		}
		if c.IdleConnTimeout > 0 {
			t.IdleConnTimeout = time.Duration(c.IdleConnTimeout) * time.Second
		}
		if c.MaxConnsPerHost > 0 {
			t.MaxConnsPerHost = c.MaxConnsPerHost
		}

		if c.MaxIdleConnsPerHost > 0 {
			t.MaxIdleConnsPerHost = c.MaxIdleConnsPerHost
		}

	} else {
		xlog.Info("[ERROR] Cannot init http.Transport")
	}
}
func MustNewAppContainer() (cont AppContainer) {

	configSource := config.MustNewAppConfigSource()

	appConfig := configSource.Config() // first call, init

	//
	// appLogger := logger.InitLogger(appConfig)
	//

	appLang := i18n.MustNewAppLang(appConfig)

	repo := repository.MustNewRepository(appConfig) // , appLogger)

	cont = &appContainer{
		configSource: configSource,
		// logger:       appLogger,
		repository: repo,
		lang:       appLang,
	}
	{
		initRuntime(appConfig)

	}
	return
}
