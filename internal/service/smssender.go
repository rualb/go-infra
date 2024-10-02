package service

import (
	"fmt"
	"go-infra/internal/config"
	xlog "go-infra/internal/tool/toollog"
	"go-infra/internal/tool/tooltaskqueue"
	"time"
)

type SmsMessage struct {
	From      string
	To        string
	Lang      string
	Text      string
	CreatedAt time.Time
	MaxAge    int16 // seconds, expires after createdAt+MaxAge if MaxAge>0
}

type SmsSender interface {
	Send(message SmsMessage) error
}

type smsSender struct {
	Debug     bool
	taskQueue *tooltaskqueue.TaskQueue[SmsMessage]
}

func (x *smsSender) Send(message SmsMessage) error {

	return x.taskQueue.Enqueue(&message)

}

type smsTaskQueue struct {
	Debug   bool
	gateway config.AppConfigMessageGateway
}

func (message *SmsMessage) exctractValueForSms(name string) (string, error) {

	switch name {
	case "to":
		return message.To, nil
	case "text":
		return message.Text, nil
	}

	return "", fmt.Errorf("prop not exists: %s", name)
}

func (x smsTaskQueue) handlerSms(smsMessage *SmsMessage) error {

	gw := x.gateway

	smsMessage.From = gw.From

	if x.Debug || gw.Stdout {
		xlog.Info("To: `%v` Message: `%v`", smsMessage.To, smsMessage.Text)
	}

	if gw.HTTP {
		sd := newDataSender()

		err := sd.fillQuery(gw, smsMessage.exctractValueForSms)

		if err != nil {
			return err
		}
		err = sd.fillBody(gw, smsMessage.exctractValueForSms)

		if err != nil {
			return err
		}

		return sd.sendData(gw)

	}
	return nil
}

func NewSmsSender(appConfig *config.AppConfig) SmsSender {

	tq := smsTaskQueue{

		Debug:   appConfig.Debug,
		gateway: appConfig.SmsGateway,
	}

	return &smsSender{
		Debug:     appConfig.Debug,
		taskQueue: tooltaskqueue.NewTaskQueue("sms sender", tq.handlerSms, 1),
	}

}
