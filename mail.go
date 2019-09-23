package main

import (
	"fmt"
	"log"
	"strconv"

	"gopkg.in/gomail.v2"
)

// SendMail SendMail
func SendMail(mailTo []string, subject string, body string, config map[string]string) error {
	port, _ := strconv.Atoi(config["port"])
	m := gomail.NewMessage()

	m.SetHeader("From", "<"+config["username"]+">")
	m.SetHeader("To", mailTo...)    //发送给多个用户
	m.SetHeader("Subject", subject) //设置邮件主题
	m.SetBody("text/html", body)    //设置邮件正文

	dialer := gomail.NewDialer(config["host"], port, config["username"], config["password"])
	err := dialer.DialAndSend(m)
	if err != nil {
		log.Println(fmt.Sprintf("[-] 发送邮件 [%s] 失败", subject), err)
	} else {
		log.Println(fmt.Sprintf("[+] 发送邮件 [%s] 成功", subject))
	}
	return err
}
