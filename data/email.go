package data

import (
	"fmt"
	"net/smtp"
)

type Email struct {
	From    string
	To      string
	Subject string
	Body    string
}

func SendEmail(useremail string) {
	// TODO : AJOUTER LES INFOS DANS LE FICHIER D'ENVIRONNEMENT
	//env, err := LoadConfig(".")
	//if err != nil {
	//	fmt.Print("cannot load config file")
	//}

	var email Email

	email.From = "exemple@negosud.com"
	email.To = useremail
	email.Subject = "negosud - vos identifiants"
	email.Body = "vos identifiants blablabla"

	auth := smtp.PlainAuth(
		"",
		"	",
		"",
		"smtp.gmail.com",
	)

	to := []string{email.To}
	msg := []byte("To: " + email.To + "\r\n" +
		"Subject: " + email.Subject + "\r\n" +
		"\r\n" +
		email.Body + "\r\n")
	err := smtp.SendMail("smtp.gmail.com:587", auth, email.From, to, msg)
	if err != nil {
		fmt.Print("Unable to send email to new user")
	}
}
