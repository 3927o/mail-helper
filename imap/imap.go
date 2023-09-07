package imap

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"io"
	"mail-helper/provider"
	"strings"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-message"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/proxy"
)

func ConnectImapServer(provider provider.MailProviderData, dialer proxy.Dialer) (*client.Client, error) {
	if dialer == nil {
		dialer = proxy.Direct
	}

	addr := fmt.Sprintf("%s:%d", provider.ImapHost, provider.ImapPort)
	c, err := client.DialWithDialerTLS(dialer, addr, &tls.Config{ServerName: provider.ImapHost})
	if err != nil {
		return nil, err
	}
	return c, nil
}

func GetMailMessages(c *client.Client, mailBoxes []string, getAllMail bool, setToSeen bool) ([]*imap.Message, error) {
	//选择收件箱
	var recentMails []*imap.Message
	if len(mailBoxes) == 0 {
		mailBoxes = GetAllMailBoxes(c)
	}

	for _, mailBox := range mailBoxes {
		//获取邮件
		var mails []*imap.Message
		mails, err := FetchMailFromBox(c, mailBox, getAllMail, setToSeen)
		if err != nil {
			return nil, errors.New("fetch mail from box " + mailBox + " failed: " + err.Error())
		}
		recentMails = append(recentMails, mails...)
	}

	return recentMails, nil
}

func GetMailMessagesSimple(c *client.Client) ([]*imap.Message, error) {
	mailBoxes := GetAllMailBoxes(c)
	return GetMailMessages(c, mailBoxes, false, false)
}

func getMailBoxes(c *client.Client, filter func(string) bool) (res []string) {
	// 获取账号所有的收件箱
	mailboxes := make(chan *imap.MailboxInfo)
	go c.List("", "*", mailboxes)
	for mailbox := range mailboxes {
		if filter(mailbox.Name) {
			res = append(res, mailbox.Name)
		}
	}
	return
}

func GetInMailBox(c *client.Client) (res []string) {
	return getMailBoxes(c, provider.IsInBox)
}

func GetJunkMailBox(c *client.Client) (res []string) {
	return getMailBoxes(c, provider.IsSpam)
}

func GetAllMailBoxes(c *client.Client) (res []string) {
	return getMailBoxes(c, func(s string) bool { return true })
}

func FetchMailFromBox(c *client.Client, mailBoxName string, getAllMail bool, setToSeen bool) ([]*imap.Message, error) {
	// 选择收件箱
	_, err := c.Select(mailBoxName, false)
	if err != nil {
		logrus.Errorf("select mailbox %s err: %s", mailBoxName, err.Error())
		return nil, err
	}

	// 查找未读邮件
	criteria := imap.NewSearchCriteria()
	if !getAllMail {
		criteria.WithoutFlags = []string{imap.SeenFlag}
	}
	ids, err := c.Search(criteria)
	if err != nil {
		logrus.Errorf("search unread mail failed: %s", err.Error())
		return nil, err
	}
	fmt.Println("unread mail ids: ", ids)

	if len(ids) == 0 {
		return []*imap.Message{}, nil
	}

	seqSet := new(imap.SeqSet)
	seqSet.AddNum(ids...)

	// 标记所有未读邮件为已读
	if setToSeen {
		item := imap.FormatFlagsOp(imap.AddFlags, true)
		flags := []interface{}{imap.SeenFlag}
		if err := c.Store(seqSet, item, flags, nil); err != nil {
			logrus.Errorf("store message to seen failed: %s", err.Error())
			return nil, err
		}
	}

	// 获取未读邮件
	messages := make(chan *imap.Message, len(seqSet.Set))
	done := make(chan error, 1)
	go func() {
		done <- c.Fetch(seqSet, []imap.FetchItem{imap.FetchRFC822, imap.FetchEnvelope}, messages)
	}()
	var recentMails []*imap.Message
	for msg := range messages {
		msg.Items["Mailbox"] = c.Mailbox().Name
		recentMails = append(recentMails, msg)
	}
	if err := <-done; err != nil {
		logrus.Errorf("fetch unread mail failed: %s", err.Error())
		return nil, err
	}

	return recentMails, nil
}

func parseEntity(e *message.Entity) map[string]interface{} {
	res := make(map[string]interface{})
	fields := e.Header.Fields()
	headers := make(map[string]string)
	for fields.Next() {
		headers[fields.Key()] = fields.Value()
	}
	res["header"] = headers
	res["attachmentFlag"] = 0

	mReader := e.MultipartReader()
	if mReader != nil {
		var parts []map[string]interface{}
		for {
			part, err := mReader.NextPart()
			if err != nil {
				break
			}
			parsedPart := parseEntity(part)
			if parsedPart["attachmentFlag"] == 1 {
				res["attachmentFlag"] = 1
			}
			parts = append(parts, parsedPart)
		}
		res["parts"] = parts
	} else {
		var a string
		if disp, _, _ := e.Header.ContentDisposition(); disp == "attachment" {
			res["attachmentFlag"] = 1
			a = "Please acquire attachment in /api/v1/<mid>.eml"
		} else {
			buf, _ := io.ReadAll(e.Body)
			a = string(buf)
		}

		res["content"] = a
	}
	return res
}

func ExtractMail(msg *imap.Message) (res map[string]interface{}, err error) {
	mid := uuid.New().String()
	emlString := ExtractMailEml(msg)

	entity, err := message.Read(bufio.NewReader(strings.NewReader(emlString)))
	if err != nil {
		return nil, err
	}
	res = parseEntity(entity)

	envelope := formatEnvelope(msg.Envelope)
	for k, v := range envelope {
		res[k] = v
	}
	res["ID"] = mid
	res["Mailbox"] = msg.Items["Mailbox"]
	return
}

func ExtractMailEml(msg *imap.Message) (res string) {
	for _, body := range msg.Body {
		s, _ := io.ReadAll(body)
		res += string(s)
	}
	return res
}

func formatEnvelope(envelope *imap.Envelope) map[string]interface{} {
	return map[string]interface{}{
		"From":       envelope.From[0].Address(),
		"To":         envelope.To[0].Address(),
		"Subject":    envelope.Subject,
		"Message-Id": envelope.MessageId,
		"Date":       envelope.Date,
	}
}

func ReceiveMail(cli *client.Client, mailBoxes []string, getAllMail bool, setToSeen bool) ([]map[string]interface{}, error) {
	messages, err := GetMailMessages(cli, mailBoxes, getAllMail, setToSeen)
	if err != nil {
		return nil, errors.New("get mail messages failed: " + err.Error())
	}

	var res []map[string]interface{}
	for _, msg := range messages {
		m, err := ExtractMail(msg)
		if err != nil {
			logrus.Errorf("extract mail subject=[%s] failed: %s", msg.Envelope.Subject, err.Error())
			continue
		}
		res = append(res, m)
	}
	return res, nil
}

func ReceiveMailSimple(account, password string) ([]map[string]interface{}, error) {
	provider, ok := provider.GetMailProviderByAddr(account)
	if !ok {
		return nil, errors.New("unkown mail provider: " + account)
	}

	cli, err := ConnectImapServer(provider, nil)
	if err != nil {
		return nil, errors.New("connect imap server fail: " + err.Error())
	}

	if err := cli.Login(account, password); err != nil {
		return nil, errors.New("login fail: " + err.Error())
	}

	messages, err := ReceiveMail(cli, nil, false, true)
	if err != nil {
		return nil, errors.New("receive mail fail: " + err.Error())
	}
	return messages, nil
}
