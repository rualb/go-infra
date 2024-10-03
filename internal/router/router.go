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
	xlog "go-infra/internal/tool/toollog"
)

func Init(e *echo.Echo, appService service.AppService) {

	initHealthController(e, appService)

	initMessengerController(e, appService)

	initConfigsController(e, appService)
}

func initConfigsController(e *echo.Echo, appService service.AppService) {

	// http://127.0.0.1:30780/configs/api/go-auth/config.development.json
	appConfig := appService.Config()

	if appConfig.Configs.Dir == "" {
		return
	}

	path, err := filepath.Abs(appConfig.Configs.Dir)

	if err != nil {
		xlog.Error("Error with configs dir: %v", appConfig.Configs.Dir)
		panic(err)
	}

	path = filepath.Clean(path)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic(fmt.Errorf("error directory does not exist: %v", path))
	}

	xlog.Info("Configs from dir: %v", path)

	e.Static(consts.PathConfigsAPI, path)

	//

}

func initHealthController(e *echo.Echo, appService service.AppService) {

	controller.SelfTest(appService)

	handler := func(c echo.Context) error {
		ctrl := controller.NewHealthController(appService, c)
		return ctrl.Check()
	}

	e.GET(consts.PathTestHealthAPI, handler)
	//
	e.GET(consts.PathTestPingAPI, func(c echo.Context) error {

		return c.String(http.StatusOK, "pong")

	})

}
func initMessengerController(e *echo.Echo, appService service.AppService) {

	factory := func(c echo.Context) *controller.MessengerController {
		return controller.NewMessengerController(appService, c)
	}

	group := e.Group(consts.PathMessengerAPI)

	group.POST("/sms-text", func(c echo.Context) error { return factory(c).SmsText() })
	group.POST("/email-html", func(c echo.Context) error { return factory(c).EmailHTML() })
	group.POST("/sms-secret-code", func(c echo.Context) error { return factory(c).SmsSecretCode() })
	group.POST("/email-secret-code", func(c echo.Context) error { return factory(c).EmailSecretCode() })

	//

}

/////////////////////////////////////////////////////
