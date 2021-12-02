package views

import (
	"log"

	"github.com/mpanelo/gocookit/models"
)

const (
	AlertLevelSuccess = "success"
	AlertLevelDanger  = "danger"
	AlertLevelWarning = "warning"

	AlertGenericMsg = "Something went wrong. Please try again. If the problem persists, contact support@gocookit.io"
)

type Data struct {
	Alert *Alert
	User  *models.User
	Yield interface{}
}

func (d *Data) SetSuccess(msg string) {
	d.Alert = &Alert{
		Level: AlertLevelSuccess,
		Msg:   msg,
	}
}

func (d *Data) SetAlertDanger(err error) {
	d.Alert = &Alert{
		Level: AlertLevelDanger,
	}

	if alerter, ok := err.(Alerter); ok {
		d.Alert.Msg = alerter.Alert()
	} else {
		log.Println(err)
		d.Alert.Msg = AlertGenericMsg
	}
}

type Alert struct {
	Level string
	Msg   string
}

type Alerter interface {
	Alert() string
}
