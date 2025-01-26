package controller

// Handler web req handler
// benchmark db http://127.0.0.1:30780/sys/api/messenger?service_code=sms_passcode&to=+000123456789&passcode=123456&lang=en
// benchmark db http://127.0.0.1:30780/sys/api/messenger?service_code=email_passcode&to=test@example.com&passcode=123456&lang=en

import (
	"fmt"
	"go-infra/internal/service"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

// MessengerController controller
type MessengerController struct {
	appService service.AppService
	webCtxt    echo.Context
	Debug      bool
}

// NewMessengerController new controller
func NewMessengerController(appService service.AppService, c echo.Context) *MessengerController {

	appConfig := appService.Config()
	return &MessengerController{
		Debug:      appConfig.Debug,
		appService: appService,
		webCtxt:    c,
	}
}

type smsPasscodeData struct {
	Message  service.SmsMessage
	Passcode string
}

type emailPasscodeData struct {
	Message  service.EmailMessage
	Passcode string
}

type messageDTO struct {
	To       string `form:"to"`
	Text     string `form:"text"`
	HTML     string `form:"html"`
	Passcode string `form:"passcode"`
	Lang     string `form:"lang"`
}

func (x messageDTO) validate() any {

	if x.To == "" {
		return map[string]string{
			"status":  "empty_arg",
			"message": "argument is empty: to",
		}
	}

	if x.Text == "" && x.HTML == "" && x.Passcode == "" {
		return map[string]string{
			"status":  "empty_arg",
			"message": "argument is empty: content",
		}
	}

	return nil
}

// SmsText send sms text
func (x *MessengerController) SmsText() error {
	/*
		url /sys/api/messenger?to=123&text=foo
	*/

	c := x.webCtxt
	dto := &messageDTO{}
	err := c.Bind(dto)
	if err != nil {
		return err
	}

	if resp := dto.validate(); resp != nil {

		return c.JSONPretty(http.StatusBadRequest, resp, "")

	}

	data := smsPasscodeData{}

	data.Message.CreatedAt = time.Now()
	data.Message.To = dto.To
	data.Message.Text = dto.Text

	err = x.appService.SmsSender().Send(data.Message)
	if err != nil {
		return err
	}

	return c.String(http.StatusOK, data.Message.Text)

}

// SmsPasscode send sms secret code
func (x *MessengerController) SmsPasscode() error {

	c := x.webCtxt
	dto := &messageDTO{}
	err := c.Bind(dto)
	if err != nil {
		return err
	}

	if resp := dto.validate(); resp != nil {

		return c.JSONPretty(http.StatusBadRequest, resp, "")

	}

	data := smsPasscodeData{}
	data.Message.MaxAge = 30 // seconds
	data.Message.CreatedAt = time.Now()
	data.Message.To = dto.To
	data.Passcode = dto.Passcode
	data.Message.Lang = dto.Lang

	data.Message.Text = fmt.Sprintf("%s: %s",
		x.appService.UserLang(data.Message.Lang).Lang("Secret code"),
		data.Passcode,
	)

	err = x.appService.SmsSender().Send(data.Message)
	if err != nil {
		return err
	}

	return c.String(http.StatusOK, data.Message.Text)

}

// EmailHTML send email html
func (x *MessengerController) EmailHTML() error {

	c := x.webCtxt
	dto := &messageDTO{}
	err := c.Bind(dto)
	if err != nil {
		return err
	}

	if resp := dto.validate(); resp != nil {

		return c.JSONPretty(http.StatusBadRequest, resp, "")

	}

	data := emailPasscodeData{}
	data.Message.CreatedAt = time.Now()
	data.Message.From = ""
	data.Message.To = dto.To
	data.Message.HTML = dto.HTML

	err = x.appService.EmailSender().Send(data.Message)
	if err != nil {
		return err
	}

	return c.HTML(http.StatusOK, data.Message.HTML)

}

// EmailPasscode send email secret code
func (x *MessengerController) EmailPasscode() error {

	appConfig := x.appService.Config()

	c := x.webCtxt
	dto := &messageDTO{}
	err := c.Bind(dto)
	if err != nil {
		return err
	}

	if resp := dto.validate(); resp != nil {

		return c.JSONPretty(http.StatusBadRequest, resp, "")

	}

	data := emailPasscodeData{}
	data.Message.CreatedAt = time.Now()
	data.Message.From = ""
	data.Message.To = dto.To
	data.Passcode = dto.Passcode
	data.Message.Lang = dto.Lang

	userLang := x.appService.UserLang(data.Message.Lang)
	labelPasscode := userLang.Lang("Secret code")
	data.Message.Subject = fmt.Sprintf("%v - %v", labelPasscode, appConfig.Title)

	bu := strings.Builder{}
	err = TemplateEmailPasscode().Execute(&bu, templateEmailPasscodeData{
		LangCode:      userLang.LangCode(),
		AppTitle:      appConfig.Title,
		LabelPasscode: labelPasscode,
		Passcode:      data.Passcode,
		Subject:       data.Message.Subject,
	})

	if err != nil {
		return err
	}

	data.Message.HTML = bu.String()

	err = x.appService.EmailSender().Send(data.Message)
	if err != nil {
		return err
	}

	return c.HTML(http.StatusOK, data.Message.HTML)

}
