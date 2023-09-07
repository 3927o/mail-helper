package provider

import (
	"strings"
)

type MailProviderData struct {
	SmtpHost string
	SmtpPort int
	ImapHost string
	ImapPort int
}

func GetMailProviderData(providerName string) (MailProviderData, bool) {
	data, ok := mailProvider[providerName]
	return data, ok
}

func GetMailProviderByAddr(addr string) (MailProviderData, bool) {
	providerName := strings.Split(addr, "@")[1]
	data, ok := GetMailProviderData(providerName)
	return data, ok
}

var mailProvider = make(map[string]MailProviderData)

var mailJunkBoxs = []string{"Junk", "Bulk", "Spam", "[Gmail]/Spam", "垃圾邮件", "垃圾郵件", "[Gmail]/垃圾邮件", "[Gmail]/垃圾郵件"}
var mailInBoxBoxs = []string{"INBOX", "Inbox", "收件箱"}

var mailJunkBoxData map[string]struct{}
var mailInBoxData map[string]struct{}

// IsSpam judge if mailbox is junk box
func IsSpam(mailBoxName string) bool {
	_, ok := mailJunkBoxData[mailBoxName]
	return ok
}

// IsInBox judge if mailbox is inbox
func IsInBox(mailBoxName string) bool {
	_, ok := mailInBoxData[mailBoxName]
	return ok
}

func RegisterMailJunkBox(mailBoxName string) {
	mailJunkBoxData[mailBoxName] = struct{}{}
}

func RegisterMailInBox(mailBoxName string) {
	mailInBoxData[mailBoxName] = struct{}{}
}

// RegisterMailProvider register mail provider, if your mail is `123@123.jp`, you are recomended to register providerName as `123.jp`, not `123`, so you can get provider quickly by addr.
func RegisterMailProvider(providerName string, data MailProviderData) {
	mailProvider[providerName] = data
}

func init() {
	mailInBoxData = make(map[string]struct{})
	mailJunkBoxData = make(map[string]struct{})
	for _, n := range mailInBoxBoxs {
		mailInBoxData[n] = struct{}{}
	}

	for _, n := range mailJunkBoxs {
		mailJunkBoxData[n] = struct{}{}
	}

	mailProvider["gmail.com"] = MailProviderData{"smtp.gmail.com", 587, "imap.gmail.com", 993}
	mailProvider["me.com"] = MailProviderData{"smtp.mail.me.com", 587, "imap.mail.me.com", 993}
	mailProvider["office365.com"] = MailProviderData{"smtp.office365.com", 587, "outlook.office365.com", 993}
	mailProvider["21cn.com"] = MailProviderData{"smtp.21cn.com", 587, "imap.21cn.com", 993}
	mailProvider["aliyun.com"] = MailProviderData{"smtp.aliyun.com", 587, "imap.aliyun.com", 993}
	mailProvider["qq.com"] = MailProviderData{"smtp.qq.com", 587, "imap.qq.com", 993}
	mailProvider["163.com"] = MailProviderData{"smtp.163.com", 25, "imap.163.com", 993}
	mailProvider["zoho.com"] = MailProviderData{"smtp.zoho.com", 587, "imappro.zoho.com", 993}
	mailProvider["aol.com"] = MailProviderData{"smtp.aol.com", 587, "imap.aol.com", 993}
	mailProvider["yahoo.com"] = MailProviderData{"smtp.mail.yahoo.com", 587, "imap.mail.yahoo.com", 993}
	// gmx
	mailProvider["gmx.com"] = MailProviderData{"mail.gmx.com", 587, "imap.gmx.com", 993}
	mailProvider["gmx.de"] = MailProviderData{"smtp.gmx.de", 587, "imap.gmx.de", 993}
	mailProvider["gmx.at"] = MailProviderData{"smtp.gmx.at", 587, "imap.gmx.at", 993}
	mailProvider["gmx.ch"] = MailProviderData{"smtp.gmx.ch", 587, "imap.gmx.ch", 993}
	mailProvider["gmx.net"] = MailProviderData{"smtp.gmx.net", 587, "imap.gmx.net", 993}

	// zoho
	mailProvider["zohomail.eu"] = MailProviderData{"smtp.zoho.eu", 587, "imap.zoho.eu", 993}
	mailProvider["zohomail.com"] = MailProviderData{"smtp.zoho.com", 587, "imap.zoho.com", 993}
	mailProvider["zohomail.jp"] = MailProviderData{"smtp.zoho.jp", 587, "imap.zoho.jp", 993}
	mailProvider["zohomail.cn"] = MailProviderData{"smtp.zoho.com.cn", 587, "imap.zoho.com.cn", 993}

	mailProvider["yandex.ru"] = MailProviderData{"smtp.Yandex.ru", 587, "imap.Yandex.ru", 993}
	mailProvider["web.de"] = MailProviderData{"smtp.web.de", 587, "imap.web.de", 993}
	mailProvider["vfemail.net"] = MailProviderData{"smtp.vfemail.net", 587, "imap.vfemail.net", 993}
	mailProvider["onet.pl"] = MailProviderData{"smtp.poczta.onet.pl", 587, "imap.poczta.onet.pl", 993}
	// outlook
	mailProvider["outlook.com"] = MailProviderData{"smtp-mail.outlook.com", 587, "imap-mail.outlook.com", 993}
	mailProvider["hotmail.com"] = MailProviderData{"smtp-mail.outlook.com", 587, "imap-mail.outlook.com", 993}

}
