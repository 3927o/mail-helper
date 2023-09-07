package message

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"mime/multipart"
	"mime/quotedprintable"
	"net/textproto"
	"strings"
)

// Message
// 发件信息
type Message struct {
	To          []string          `json:"to"`
	Nickname    string            `json:"nickname"`
	From        string            `json:"from"`
	Subject     string            `json:"subject"`
	Content     []byte            `json:"content"`
	Headers     map[string]string `json:"headers"`
	Attachments map[string][]byte `json:"attachments"`
}

func BuildMessageSimple(from, subject, content string, to []string) Message {
	message := Message{
		From:    from,
		To:      to,
		Subject: subject,
		Content: []byte(content),
	}
	return message
}

func BuildMessageWithNickname(from, subject, content, nickname string, to []string) Message {
	message := Message{
		From:     from,
		To:       to,
		Subject:  subject,
		Content:  []byte(content),
		Nickname: nickname,
	}
	return message
}

func (m Message) Bytes() []byte {
	buf := bytes.NewBuffer(nil)

	_, _ = fmt.Fprintf(buf, "Subject: =?UTF-8?B?%s?=\r\n", base64.StdEncoding.EncodeToString([]byte(m.Subject)))
	if len(m.To) == 0 {
		_, _ = fmt.Fprintf(buf, "To: %s\r\n", strings.Join(m.To, ","))
	}
	if m.Nickname != "" && m.From != "" {
		_, _ = fmt.Fprintf(buf, "From: =?UTF-8?B?%s?= <%s>\r\n", base64.StdEncoding.EncodeToString([]byte(m.Nickname)), m.From)
	} else if m.From != "" {
		_, _ = fmt.Fprintf(buf, "From: %s\r\n", m.From)
	}

	for k, v := range m.Headers {
		_, _ = fmt.Fprintf(buf, "%s: %s\r\n", k, v)
	}

	buf.WriteString("MIME-Version: 1.0\r\n")

	writer := multipart.NewWriter(buf)
	boundary := writer.Boundary()

	_, _ = fmt.Fprintf(buf, "Content-Type: multipart/mixed; boundary=%s\r\n\r\n", boundary)
	//_ , _ = fmt.Fprintf(buf, "--%s\r\n\r\n", boundary)
	htmlPart, _ := writer.CreatePart(textproto.MIMEHeader{"Content-Type": {"text/html;charset=\"UTF-8\""}, "Content-Transfer-Encoding": {"quoted-printable"}})

	qp := quotedprintable.NewWriter(htmlPart)
	_, _ = qp.Write(m.Content)
	_ = qp.Close()

	buf.WriteString("\r\n\r\n")
	if len(m.Attachments) > 0 {
		for k, v := range m.Attachments {
			filename := base64.StdEncoding.EncodeToString([]byte(k))

			content := base64.StdEncoding.EncodeToString(v)
			filePart, _ := writer.CreatePart(textproto.MIMEHeader{
				"Content-Type":              {"application/octet-stream"},
				"Content-Transfer-Encoding": {"base64"},
				"Content-Disposition":       {"attachment; filename=\"=?UTF-8?B?" + filename + "?=\""}},
			)
			// 一行超过76好像也没事？
			// _, _ = filePart.Write([]byte(content))

			qp := quotedprintable.NewWriter(filePart)
			_, _ = qp.Write([]byte(content))
			_ = qp.Close()
			buf.WriteString("\r\n\r\n")
		}
	}

	_ = writer.Close()

	return buf.Bytes()
}
