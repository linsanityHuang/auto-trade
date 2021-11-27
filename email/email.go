package main

import (
	"auto-trade/global"
	"auto-trade/pkg/setting"
	"fmt"
	"log"
	"net/smtp"
)

func init() {
	err := setupSetting()
	if err != nil {
		log.Fatalf("init.setupSetting err: %v", err)
	}
}

func setupSetting() error {
	setting, err := setting.NewSetting()
	if err != nil {
		return err
	}
	err = setting.ReadSection("Email", &global.EmailSetting)
	if err != nil {
		return err
	}

	return nil
}

func SendEmail(receiver string) {
	auth := smtp.PlainAuth("", global.EmailSetting.SMTPUsername, global.EmailSetting.SMTPPassword, global.EmailSetting.SMTPHost)
	msg := []byte("Subject: 这里是标题内容\r\n\r\n" + "这里是正文内容\r\n")
	addr := fmt.Sprintf("%s:%d", global.EmailSetting.SMTPHost, global.EmailSetting.SMTPPort)
	err := smtp.SendMail(addr, auth, global.EmailSetting.SMTPUsername, []string{receiver}, msg)

	if err != nil {
		log.Fatalf("failed to send email: %v\n", err)
	}
}

func main() {
	SendEmail(global.EmailSetting.Receiver)
}
