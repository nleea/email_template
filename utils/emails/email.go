package emails

import (
	"gopkg.in/gomail.v2"
	CO "sequency/config"
)

func SendEmail(from string, to string, subject string, template string) {

	envs := CO.ConfigEnv()

	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", "Hello <b>Bob</b> and <i>Cora</i>!")

	d := gomail.NewDialer("smtp.gmail.com", 587, envs["USER_EMAIL"], envs["PASSWORD_EMAIL"])

	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}

}
