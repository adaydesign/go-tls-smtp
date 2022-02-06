package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/smtp"
)

func main() {

	// Sender data.
	from := "[email]@gmail.com" // <------------- (1) แก้ไขอีเมลที่ใช้ส่ง
	password := "[passwos]"     // <------- (2) แก้ไขรหัสผ่านของอีเมลที่ใช้ส่ง

	// Receiver email address.
	to := "[receiver's email]" // <-------------- (3) แก้ไขอีเมลของผู้รับ หากใส่หลายเมล จะไปอยู่ที่ cc

	// smtp server configuration.
	smtpHost := "smtp.gmail.com"
	smtpPort := "465" // <----------------- (4) port ใช้ tls คือ 465
	servername := smtpHost + ":" + smtpPort
	// Message.
	subj := "Hello !"                 // <----------------- (5) หัวเรื่อง
	body := "This is a test message." // <----------------- (6) เนื้อความ

	// Setup headers
	headers := make(map[string]string)
	headers["From"] = from
	headers["To"] = to
	headers["Subject"] = subj

	// Setup message
	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body

	// Authentication
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// TLS config
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         smtpHost,
	}

	// ต้องเรียกใช้ tls.Dial แทน smtp.Dial
	// สำหรับเซิฟเวอร์ smtp ที่รันบนพอร์ท 465 ต้องทำการเชื่อมต่อเป็นแบบ ssl connection

	conn, err := tls.Dial("tcp", servername, tlsconfig)
	if err != nil {
		log.Panic(err)
	}

	c, err := smtp.NewClient(conn, smtpHost)
	if err != nil {
		log.Panic(err)
	}

	// Auth
	if err = c.Auth(auth); err != nil {
		log.Panic(err)
	}

	// To && From
	if err = c.Mail(from); err != nil {
		log.Panic(err)
	}

	if err = c.Rcpt(to); err != nil {
		log.Panic(err)
	}

	// Data
	w, err := c.Data()
	if err != nil {
		log.Panic(err)
	}

	// Sending email.
	_, err = w.Write([]byte(message))
	if err != nil {
		log.Panic(err)
	}

	err = w.Close()
	if err != nil {
		log.Panic(err)
	}

	c.Quit()

	fmt.Println("Email Sent Successfully!")
}
