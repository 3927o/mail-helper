package smtp

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"mail-helper/message"
	"mail-helper/provider"
	"net"
	"net/smtp"

	"golang.org/x/net/proxy"
)

type loginAuth struct {
	username, password string
}

func (a *loginAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	return "LOGIN", []byte{}, nil
}

func (a *loginAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if more {
		switch string(fromServer) {
		case "Username:":
			return []byte(a.username), nil
		case "Password:":
			return []byte(a.password), nil
		default:
			return nil, errors.New("unkown fromServer")
		}
	}
	return nil, nil
}

func NewLoginAuth(username, password string) smtp.Auth {
	return &loginAuth{username, password}
}

func newSmtpClientWithConn(conn net.Conn, provider provider.MailProviderData) (*smtp.Client, error) {
	client, err := smtp.NewClient(conn, provider.SmtpHost)
	if err != nil {
		return nil, fmt.Errorf("new smtp client of host %s fail: %s", provider.SmtpHost, err)
	}

	if ok, _ := client.Extension("STARTTLS"); ok {
		if err := client.StartTLS(&tls.Config{ServerName: provider.SmtpHost}); err != nil {
			_ = client.Close()
			//log.Tracef("use ssl fail:%s", err)
			return nil, fmt.Errorf("use ssl fail:%s", err)
		}
	}

	return client, nil
}

func ConnectWithContext(ctx context.Context, provider provider.MailProviderData, dialer proxy.ContextDialer) (*smtp.Client, error) {
	if dialer == nil {
		dialer = proxy.Direct
	}

	conn, err := dialer.DialContext(ctx, "tcp", fmt.Sprintf("%s:%d", provider.SmtpHost, provider.SmtpPort))
	if err != nil {
		return nil, fmt.Errorf("dial tcp %s:%d fail: %s", provider.SmtpHost, provider.SmtpPort, err)
	}

	client, err := newSmtpClientWithConn(conn, provider)
	return client, err
}

func Connect(provider provider.MailProviderData, dialer proxy.Dialer) (*smtp.Client, error) {
	if dialer == nil {
		dialer = proxy.Direct
	}

	conn, err := dialer.Dial("tcp", fmt.Sprintf("%s:%d", provider.SmtpHost, provider.SmtpPort))
	if err != nil {
		return nil, fmt.Errorf("dial tcp %s:%d fail: %s", provider.SmtpHost, provider.SmtpPort, err)
	}

	client, err := newSmtpClientWithConn(conn, provider)
	return client, err
}

func ConnectSimple(provider provider.MailProviderData) (*smtp.Client, error) {
	return Connect(provider, nil)
}

func ConnectSimpleWithContext(ctx context.Context, provider provider.MailProviderData) (*smtp.Client, error) {
	return ConnectWithContext(ctx, provider, nil)
}

func SendMessage(client *smtp.Client, message message.Message) error {
	if err := client.Mail(message.From); err != nil {
		return fmt.Errorf("set mail from fail: %s", err)
	}

	for _, addr := range message.To {
		if err := client.Rcpt(addr); err != nil {
			return fmt.Errorf("add rcpt %s fail: %s", addr, err)
		}
	}

	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("get data writer fail: %s", err)
	}

	_, err = w.Write(message.Bytes())
	if err != nil {
		return fmt.Errorf("write data fail: %s", err)
	}

	err = w.Close()
	if err != nil {
		return fmt.Errorf("close data writer fail: %s", err)
	}

	return nil
}

func SendMessageSimple(account, password, to, subject, content string) error {
	provider, ok := provider.GetMailProviderByAddr(account)
	if !ok {
		return errors.New("unkown mail provider: " + account)
	}

	client, err := ConnectSimple(provider)
	if err != nil {
		return errors.New("new smtp client fail: " + err.Error())
	}

	auth := NewLoginAuth(account, password)
	if err := client.Auth(auth); err != nil {
		return errors.New("auth fail: " + err.Error())
	}

	msg := message.BuildMessageSimple(account, subject, content, []string{to})
	if err := SendMessage(client, msg); err != nil {
		return errors.New("send message fail: " + err.Error())
	}
	return nil
}
