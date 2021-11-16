package views

const (
	AlertLevelSuccess = "success"
	AlertLevelDanger  = "danger"
	AlertLevelWarning = "warning"

	AlertGenericMsg = "Something went wrong. Please try again. If the problem persists, contact support@gocookit.io"
)

type Data struct {
	Alert *Alert
	Yield interface{}
}

func (d *Data) SetAlertDanger(err error) {
	d.Alert = &Alert{
		Level: AlertLevelDanger,
		Msg:   AlertGenericMsg,
	}

	if alerter, ok := err.(Alerter); ok {
		d.Alert.Msg = alerter.Alert()
	}
}

type Alert struct {
	Level string
	Msg   string
}

type Alerter interface {
	Alert() string
}
