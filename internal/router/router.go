// Package router main rounter
package router

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/labstack/echo/v4"

	controller "go-infra/internal/controller"

	"go-infra/internal/config/consts"
	"go-infra/internal/service"

	xlog "go-infra/internal/util/utillog"

	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4/middleware"
)

func Init(e *echo.Echo, appService service.AppService) {

	initDebugController(e, appService)

	initMessengerController(e, appService)

	initConfigsController(e, appService)

	initSys(e, appService)
}
func initSys(e *echo.Echo, appService service.AppService) {

	// !!! DANGER for private(non-public) services only
	// or use non-public port via echo.New()

	appConfig := appService.Config()

	listen := appConfig.HTTPServer.Listen
	listenSys := appConfig.HTTPServer.ListenSys
	sysMetrics := appConfig.HTTPServer.SysMetrics
	hasAnyService := sysMetrics
	sysAPIKey := appConfig.HTTPServer.SysAPIKey
	hasAPIKey := sysAPIKey != ""
	hasListenSys := listenSys != ""
	startNewListener := listenSys != listen

	if !hasListenSys {
		return
	}

	if !hasAnyService {
		return
	}

	if !hasAPIKey {
		xlog.Panic("sys api key is empty")
		return
	}

	if startNewListener {

		e = echo.New() // overwrite override

		e.Use(middleware.Recover())
		// e.Use(middleware.Logger())
	} else {
		xlog.Warn("sys api serve in main listener: %v", listen)
	}

	sysAPIAccessAuthMW := middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		KeyLookup: "query:api-key,header:Authorization",
		Validator: func(key string, c echo.Context) (bool, error) {
			return key == sysAPIKey, nil
		},
	})

	if sysMetrics {
		// may be eSys := echo.New() // this Echo will run on separate port
		e.GET(
			consts.PathSysMetricsAPI,
			echoprometheus.NewHandler(),
			sysAPIAccessAuthMW,
		) // adds route to serve gathered metrics

	}

	if startNewListener {

		// start as async task
		go func() {
			xlog.Info("sys api serve on: %v main: %v", listenSys, listen)

			if err := e.Start(listenSys); err != nil {
				if err != http.ErrServerClosed {
					xlog.Error("%v", err)
				} else {
					xlog.Info("shutting down the server")
				}
			}
		}()

	} else {
		xlog.Info("sys api server serve on main listener: %v", listen)
	}

}
func initConfigsController(e *echo.Echo, appService service.AppService) {

	// http://127.0.0.1:30780/sys/api/configs/go-auth/config.development.json
	appConfig := appService.Config()

	if appConfig.Configs.Dir == "" {
		return
	}

	path, err := filepath.Abs(appConfig.Configs.Dir)

	if err != nil {
		xlog.Error("error with configs dir: %v", appConfig.Configs.Dir)
		panic(err)
	}

	path = filepath.Clean(path)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic(fmt.Errorf("error directory does not exist: %v", path))
	}

	xlog.Info("configs from dir: %v", path)

	e.Static(consts.PathSysConfigsAPI, path)

	//

}

func initDebugController(e *echo.Echo, _ service.AppService) {
	e.GET(consts.PathInfraPingDebugAPI, func(c echo.Context) error { return c.String(http.StatusOK, "pong") })
	// publicly-available-no-sensitive-data
	e.GET("/health", func(c echo.Context) error { return c.JSON(http.StatusOK, struct{}{}) })

}
func initMessengerController(e *echo.Echo, appService service.AppService) {

	factory := func(c echo.Context) *controller.MessengerController {
		return controller.NewMessengerController(appService, c)
	}

	group := e.Group(consts.PathSysMessengerAPI)

	group.POST("/sms-text", func(c echo.Context) error { return factory(c).SmsText() })
	group.POST("/email-html", func(c echo.Context) error { return factory(c).EmailHTML() })
	group.POST("/sms-passcode", func(c echo.Context) error { return factory(c).SmsPasscode() })
	group.POST("/email-passcode", func(c echo.Context) error { return factory(c).EmailPasscode() })

	//

}

/////////////////////////////////////////////////////
