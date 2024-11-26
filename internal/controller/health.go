package controller

import (
	"fmt"
	"go-infra/internal/service"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

// HealthController controller
type HealthController struct {
	appService service.AppService
	webCtxt    echo.Context
}

// NewHealthController new controller
func NewHealthController(appService service.AppService, c echo.Context) *HealthController {
	return &HealthController{
		appService: appService,
		webCtxt:    c,
	}
}

// Check health
// benchmark db ?cmd=db
// benchmark db ?cmd=db_raw
// get string of length ?cmd=length&length=10
// get stats db ?cmd=stats
func (x *HealthController) Check() error {

	c := x.webCtxt
	res := ""

	switch c.QueryParam("cmd") {

	case "db":

		//
		if val, err := checkRunDB(x.appService); err != nil {
			return err
		} else {
			res = val
		}

	case "db_raw":

		if val, err := checkRunDBRaw(x.appService); err != nil {
			return err
		} else {
			res = val
		}

	case "stats":

		if val, err := checkStats(x.appService); err != nil {
			return err
		} else {
			res = val
		}

	default:
		return c.String(http.StatusBadRequest, "cmd not defined")
	}

	return c.String(http.StatusOK, res)
}

// SelfTest run benchmark for db
func SelfTest(_ service.AppService) {

	// count := 1000

	// {
	// 	_, _ = checkRunDB(appService)
	// 	df := time.Now()
	// 	for i := 0; i < count; i++ {
	// 		_, _ = checkRunDB(appService)
	// 	}
	// 	xlog.Info("Self test RunDb Msec:%v Count:%v", time.Since(df).Milliseconds(), count)
	// }

	// {
	// 	_, _ = checkRunDBRaw(appService)
	// 	df := time.Now()
	// 	for i := 0; i < count; i++ {
	// 		_, _ = checkRunDBRaw(appService)
	// 	}
	// 	xlog.Info("Self test RunDbRaw Msec:%v Count:%v", time.Since(df).Milliseconds(), count)
	// }

}
func checkRunDB(appService service.AppService) (string, error) {
	type Result struct {
		Value int
		Now   time.Time
	}
	//
	data := Result{}
	repo := appService.Repository()
	err := repo.Raw("select cast(? as int)+cast(? as int) as value, now() as now", 2, 3).Scan(&data).Error
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(" %v %v", data.Value, data.Now), nil

}
func checkStats(appService service.AppService) (string, error) {

	bu := strings.Builder{}

	{
		// stats from db

		db, err := appService.Repository().Driver().DB()

		if err != nil {
			return "", err
		}
		stats := db.Stats()

		fmt.Fprintf(&bu, "db_max_open_connections %d\n", stats.MaxOpenConnections)
		fmt.Fprintf(&bu, "db_open_connections %d\n", stats.OpenConnections)
		fmt.Fprintf(&bu, "db_in_use_connections %d\n", stats.InUse)
		fmt.Fprintf(&bu, "db_idle_connections %d\n", stats.Idle)
		fmt.Fprintf(&bu, "db_wait_count %d\n", stats.WaitCount)
		fmt.Fprintf(&bu, "db_wait_duration_msec %v\n", stats.WaitDuration.Milliseconds())
		fmt.Fprintf(&bu, "db_max_idle_closed %d\n", stats.MaxIdleClosed)
		fmt.Fprintf(&bu, "db_max_lifetime_closed %d\n", stats.MaxLifetimeClosed)

	}

	return bu.String(), nil
}
func stringOfLength(_ service.AppService, length int) string {

	length = min(max(length, 0), 32000)

	return strings.Repeat("A", length)
}

func checkRunDBRaw(appService service.AppService) (string, error) {

	type Data struct {
		Value int
		Now   time.Time
	}
	data := Data{}
	repo := appService.Repository()
	db, _ := repo.Driver().DB()

	// postgres driver $1, mysql,sql ?
	buf := db.QueryRow(`select cast($1 as int)+cast($2 as int) as value, now() as now`, 2, 3)

	if buf.Err() != nil {
		return "", buf.Err()
	}

	_ = buf.Scan(&data.Value, &data.Now)

	return fmt.Sprintf(" %v %v", data.Value, data.Now), nil

}
