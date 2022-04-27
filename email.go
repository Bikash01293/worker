package work

import (
	"fmt"
	"net/smtp"
)



func Email() {
	// Sender data.
	from := "bikashkumar01293@gmail.com"
	password := "nhamkqzzmygfvdkh"

	// Receiver email address.
	to := []string{
		"adhikaribikash821@gmail.com",
	}

	// smtp server configuration.
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	// Message.
	message := []byte("This is a test email message.")

	// Authentication.
	auth := smtp.PlainAuth("", from, password, smtpHost)

	
	// Sending email.
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	if err != nil {
		fmt.Println(err)
		return
	}


	fmt.Println("\nEmail Sent Successfully!")
}


// func sendWork() {
// 	fmt.Println("I am bikash")
// }


