package service

import (
	"fmt"
	"go-infra/internal/config"
	xlog "go-infra/internal/tool/toollog"
	"go-infra/internal/tool/tooltaskqueue"
	"time"
)

type EmailMessage struct {
	From      string
	To        string
	Lang      string
	Subject   string
	HTML      string
	CreatedAt time.Time
	MaxAge    int16 // expires after createdAt+MaxAge if MaxAge>0
}

type EmailSender interface {
	Send(message EmailMessage) error
}

type emailSender struct {
	Debug     bool
	taskQueue *tooltaskqueue.TaskQueue[EmailMessage]
}

func (x *emailSender) Send(message EmailMessage) error {

	return x.taskQueue.Enqueue(&message)

}

type emailTaskQueue struct {
	Debug   bool
	gateway config.AppConfigMessageGateway
}

func (message *EmailMessage) exctractValueForEmail(name string) (string, error) {

	//	message.From = fmt.Sprintf("%s <$s>", message.From, x.gateway.From)
	switch name {
	case "from": // important field for gaiteways `title <mail>`
		return message.From, nil
	case "to":
		return message.To, nil
	case "subject":
		return message.Subject, nil
	case "html":
		return message.HTML, nil
	}

	return "", fmt.Errorf("prop not exists: %s", name)
}
func (x emailTaskQueue) handlerEmail(emailMessage *EmailMessage) error {

	gw := x.gateway

	emailMessage.From = gw.From

	if x.Debug || gw.Stdout {
		xlog.Info("To: `%v` Subject: `%v` Message: `%v`", emailMessage.To, emailMessage.Subject, emailMessage.HTML)
	}

	if gw.HTTP {

		sd := newDataSender()

		err := sd.fillQuery(gw, emailMessage.exctractValueForEmail)

		if err != nil {
			return err
		}
		err = sd.fillBody(gw, emailMessage.exctractValueForEmail)

		if err != nil {
			return err
		}

		return sd.sendData(gw)

	}
	return nil
}

func NewEmailSender(appConfig *config.AppConfig) EmailSender {

	tq := emailTaskQueue{

		Debug:   appConfig.Debug,
		gateway: appConfig.EmailGateway,
	}

	return &emailSender{
		Debug:     appConfig.Debug,
		taskQueue: tooltaskqueue.NewTaskQueue("email sender", tq.handlerEmail, 1),
	}
}
