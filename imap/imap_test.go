package imap

import (
	"fmt"
	"testing"
)

func TestReceiveMailSimple(t *testing.T) {
	account := "lin3927@outlook.com"
	pwd := "lovqdtzwhquojdeg"

	if messages, err := ReceiveMailSimple(account, pwd); err != nil {
		t.Error("receive mail fail: " + err.Error())
	} else {
		fmt.Println(messages)
	}
}
