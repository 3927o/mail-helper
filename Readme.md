# Go-Mail-Helper
a helper for go to send/receive mail by protocol

## Usage
send mail
```go
func TestSendMessageSimple(t *testing.T) {
	account := "xxxxx@qq.com"
	pwd := "xxxxx"

	if err := SendMessageSimple(account, pwd, account, "test", "testaaaaaaa"); err != nil {
		t.Error("send message fail: " + err.Error())
	}
}

func TestSendMail(t *testing.T) {
	account := "xxxxx@qq.com"
	pwd := "xxxxx"

	provider, _ := provider.GetMailProviderByAddr(account)
	cli, _ := ConnectSimple(provider)
s
	auth := NewLoginAuth(account, pwd)
	err := cli.Auth(auth)
	if err != nil {
		t.Error("auth fail: " + err.Error())
	}

	if err := SendMessage(cli, message.Message{From: account, To: []string{account}, Subject: "test", Content: []byte("testaaaaaaa")}); err != nil {
		t.Error("send message fail: " + err.Error())
	}
}
```

receive mail
```go
func TestReceiveMailSimple(t *testing.T) {
	account := "xxxxx@qq.com"
	pwd := "xxxxx"

	if messages, err := ReceiveMailSimple(account, pwd); err != nil {
		t.Error("receive mail fail: " + err.Error())
	} else {
		fmt.Println(messages)
	}
}
```