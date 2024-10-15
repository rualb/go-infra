package controller

// Handler web req handler
// benchmark db http://127.0.0.1:30780/messenger/api?service_code=sms_secret_code&to=+000123456789&secret_code=123456&lang=en
// benchmark db http://127.0.0.1:30780/messenger/api?service_code=email_secret_code&to=test@example.com&secret_code=123456&lang=en

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

type smsSecretCodeData struct {
	Message    service.SmsMessage
	SecretCode string
}

type emailSecretCodeData struct {
	Message    service.EmailMessage
	SecretCode string
}

type messageDto struct {
	To         string `form:"to"`
	Text       string `form:"text"`
	HTML       string `form:"html"`
	SecretCode string `form:"secret_code"`
	Lang       string `form:"lang"`
}

func (x messageDto) validate() any {

	if x.To == "" {
		return map[string]string{
			"status":  "empty_arg",
			"message": "argument is empty: to",
		}
	}

	if x.Text == "" && x.HTML == "" && x.SecretCode == "" {
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
		url /messenger/api?to=123&text=foo
	*/

	c := x.webCtxt
	dto := &messageDto{}
	err := c.Bind(dto)
	if err != nil {
		return err
	}

	if resp := dto.validate(); resp != nil {

		return c.JSONPretty(http.StatusBadRequest, resp, "")

	}

	data := smsSecretCodeData{}

	data.Message.CreatedAt = time.Now()
	data.Message.To = dto.To
	data.Message.Text = dto.Text

	err = x.appService.SmsSender().Send(data.Message)
	if err != nil {
		return err
	}

	return c.String(http.StatusOK, data.Message.Text)

}

// SmsSecretCode send sms secret code
func (x *MessengerController) SmsSecretCode() error {

	c := x.webCtxt
	dto := &messageDto{}
	err := c.Bind(dto)
	if err != nil {
		return err
	}

	if resp := dto.validate(); resp != nil {

		return c.JSONPretty(http.StatusBadRequest, resp, "")

	}

	data := smsSecretCodeData{}
	data.Message.MaxAge = 30 // seconds
	data.Message.CreatedAt = time.Now()
	data.Message.To = dto.To
	data.SecretCode = dto.SecretCode
	data.Message.Lang = dto.Lang

	data.Message.Text = fmt.Sprintf("%v: %v", x.appService.UserLang(data.Message.Lang).Lang("Secret code"), data.SecretCode)

	err = x.appService.SmsSender().Send(data.Message)
	if err != nil {
		return err
	}

	return c.String(http.StatusOK, data.Message.Text)

}

// EmailHTML send email html
func (x *MessengerController) EmailHTML() error {

	c := x.webCtxt
	dto := &messageDto{}
	err := c.Bind(dto)
	if err != nil {
		return err
	}

	if resp := dto.validate(); resp != nil {

		return c.JSONPretty(http.StatusBadRequest, resp, "")

	}

	data := emailSecretCodeData{}
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

// EmailSecretCode send email secret code
func (x *MessengerController) EmailSecretCode() error {

	appConfig := x.appService.Config()

	c := x.webCtxt
	dto := &messageDto{}
	err := c.Bind(dto)
	if err != nil {
		return err
	}

	if resp := dto.validate(); resp != nil {

		return c.JSONPretty(http.StatusBadRequest, resp, "")

	}

	data := emailSecretCodeData{}
	data.Message.CreatedAt = time.Now()
	data.Message.From = ""
	data.Message.To = dto.To
	data.SecretCode = dto.SecretCode
	data.Message.Lang = dto.Lang

	userLang := x.appService.UserLang(data.Message.Lang)
	labelSecretCode := userLang.Lang("Secret code")
	data.Message.Subject = fmt.Sprintf("%v - %v", labelSecretCode, appConfig.Title)

	bu := strings.Builder{}
	err = TemplateEmailSecretCode().Execute(&bu, templateEmailSecretCodeData{
		LangCode:        userLang.LangCode(),
		AppTitle:        appConfig.Title,
		LabelSecretCode: labelSecretCode,
		SecretCode:      data.SecretCode,
		Subject:         data.Message.Subject,
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
