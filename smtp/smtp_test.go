package smtp

import (
	"mail-helper/message"
	"mail-helper/provider"
	"os"
	"testing"
)

func TestSendMail(t *testing.T) {
	account := "kmakunitwaire@gmx.com"
	pwd := "hPNVUiuwXV"
	to := "mikumikumi233@gmail.com"

	provider, _ := provider.GetMailProviderByAddr(account)
	cli, _ := ConnectSimple(provider)
	auth := NewLoginAuth(account, pwd)
	err := cli.Auth(auth)
	if err != nil {
		t.Error("auth fail: " + err.Error())
		return
	}

	message := message.Message{From: account, To: []string{to}, Subject: "test", Content: []byte("testaaaaaaa")}
	fileByte, _ := os.ReadFile(`C:\Users\123\Downloads\123123123.pdf`)
	message.Attachments = map[string][]byte{"123123123.pdf": fileByte}

	if err := SendMessage(cli, message); err != nil {
		t.Error("send message fail: " + err.Error())
		return
	}
}

func TestSendMessageSimple(t *testing.T) {
	account := "1624497311@qq.com"
	pwd := "xsrusnnifldoddaa"

	if err := SendMessageSimple(account, pwd, account, "test", "testaaaaaaa"); err != nil {
		t.Error("send message fail: " + err.Error())
	}
}

func TestSimpleSend(t *testing.T) {
	account := "1624497311@qq.com"
	pwd := "xsrusnnifldoddaa"

	provider, ok := provider.GetMailProviderByAddr(account)
	if !ok {
		t.Error("unkown mail provider: " + account)
	}
	client, err := ConnectSimple(provider)
	if err != nil {
		t.Error("new smtp client fail: " + err.Error())
	}

	if err = client.Auth(NewLoginAuth(account, pwd)); err != nil {
		t.Error("auth fail: " + err.Error())
	}
	if err = client.Mail(account); err != nil {
		t.Error("set mail from fail: " + err.Error())
	}
	if err = client.Rcpt(account); err != nil {
		t.Error("set mail to fail: " + err.Error())
	}
	w, err := client.Data()
	if err != nil {
		t.Error("get data writer fail: " + err.Error())
	}
	_, err = w.Write(message.Message{From: "1624497311@qq.com"}.Bytes())
	if err != nil {
		t.Error("write data fail: " + err.Error())
	}
	err = w.Close()
	if err != nil {
		t.Error("close data writer fail: " + err.Error())
	}
}
