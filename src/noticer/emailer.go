package noticer

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"model"
	"net/smtp"
	"strings"
	"time"
)

// define email interface, and implemented auth and send method
type Mail interface {
	Auth()
	Send(message Message) error
}

type SendMail struct {
	user     string
	password string
	host     string
	port     string
	auth     smtp.Auth
}

type Attachment struct {
	name        string
	contentType string
	withFile    bool
}

type Message struct {
	from        string
	to          []string
	cc          []string
	bcc         []string
	subject     string
	body        string
	contentType string
	attachment  Attachment
}

func (mail *SendMail) Auth() {
	mail.auth = smtp.PlainAuth("", mail.user, mail.password, mail.host)
}

func (mail SendMail) Send(message Message) error {
	mail.Auth()
	buffer := bytes.NewBuffer(nil)
	boundary := "GoBoundary"
	Header := make(map[string]string)
	Header["From"] = message.from
	Header["To"] = strings.Join(message.to, ";")
	Header["Cc"] = strings.Join(message.cc, ";")
	Header["Bcc"] = strings.Join(message.bcc, ";")
	Header["Subject"] = message.subject
	Header["Content-Type"] = "multipart/mixed;boundary=" + boundary
	Header["Mime-Version"] = "1.0"
	Header["Date"] = time.Now().String()
	mail.writeHeader(buffer, Header)

	body := "\r\n--" + boundary + "\r\n"
	body += "Content-Type:" + message.contentType + "\r\n"
	body += "\r\n" + message.body + "\r\n"
	buffer.WriteString(body)

	if message.attachment.withFile {
		attachment := "\r\n--" + boundary + "\r\n"
		attachment += "Content-Transfer-Encoding:base64\r\n"
		attachment += "Content-Disposition:attachment\r\n"
		attachment += "Content-Type:" + message.attachment.contentType + ";name=\"" + message.attachment.name + "\"\r\n"
		buffer.WriteString(attachment)
		defer func() {
			if err := recover(); err != nil {
				log.Fatalln(err)
			}
		}()
		mail.writeFile(buffer, message.attachment.name)
	}

	buffer.WriteString("\r\n--" + boundary + "--")
	smtp.SendMail(mail.host+":"+mail.port, mail.auth, message.from, message.to, buffer.Bytes())
	return nil
}

func (mail SendMail) writeHeader(buffer *bytes.Buffer, Header map[string]string) string {
	header := ""
	for key, value := range Header {
		header += key + ":" + value + "\r\n"
	}
	header += "\r\n"
	buffer.WriteString(header)
	return header
}

// read and write the file to buffer
func (mail SendMail) writeFile(buffer *bytes.Buffer, fileName string) {
	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		panic(err.Error())
	}
	payload := make([]byte, base64.StdEncoding.EncodedLen(len(file)))
	base64.StdEncoding.Encode(payload, file)
	buffer.WriteString("\r\n")
	for index, line := 0, len(payload); index < line; index++ {
		buffer.WriteByte(payload[index])
		if (index+1)%76 == 0 {
			buffer.WriteString("\r\n")
		}
	}
}

func getSender() model.EmailAccount{
	data, err := ioutil.ReadFile("account.json")
	if err != nil{
		fmt.Println("配置文件不存在，请检查")
	}
	account := model.EmailAccount{}
	err = json.Unmarshal(data, &account)
	return account
}

func SendEmail() {
	sender := getSender()
	var mail Mail
	mail = &SendMail{user: sender.Sender, password: sender.Pwd, host: "smtp.qq.com", port: "25"}
	message := Message{from: sender.Sender,
		to: []string{sender.Sender},
		cc: []string{},
		bcc: []string{},
		subject: "库存信息",
		body: "",
		contentType: "text/plain;charset=utf-8",
		attachment: Attachment{
			name:        "output.xlsx",
			contentType: "application/octet-stream",
			withFile:    true,
		},
	}
	err := mail.Send(message)
	if err != nil{
		fmt.Println("邮件发送失败")
		fmt.Println(err)
		return
	}
	fmt.Println("邮件发送成功")
}
