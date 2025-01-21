package common4all

import (
	"bytes"
	"fmt"
	"github.com/strongo/strongoapp"
)

type deeplink struct {
}

func (deeplink) AppHashPathToReceipt(receiptID string) string {
	return fmt.Sprintf("receipt=%s", receiptID)
}

var Deeplink = deeplink{}

type Linker struct {
	userID string
	locale string
	issuer string
	host   string
}

func NewLinker(environment string, userID string, locale, issuer string) Linker {
	return Linker{
		userID: userID,
		locale: locale,
		issuer: issuer,
		host:   host(environment),
	}
}

func host(environment string) string {
	switch environment {
	case "prod":
		return "debtus.app"
	case strongoapp.LocalHostEnv:
		return "local.debtus.app"
	case "dev":
		return "dev1.debtus.app"
	}
	panic(fmt.Sprintf("Unknown environment: %v", environment))
}

func (l Linker) UrlToContact(contactID string) string {
	return l.url("/contact", fmt.Sprintf("?id=%s", contactID), "")
}

func (l Linker) url(path, query, hash string) string {
	var buffer bytes.Buffer
	buffer.WriteString("https://" + l.host + path + query)
	if hash != "" {
		buffer.WriteString(hash)
	}
	if query != "" || hash != "" {
		buffer.WriteString("&")
	}
	//isAdmin := false // TODO: How to get isAdmin?
	//token, _ := token4auth.IssueAuthToken(ctx, l.userID, l.issuer)
	buffer.WriteString("lang=" + l.locale)
	buffer.WriteString("&secret=TODO")
	return buffer.String()
}

func (l Linker) ToMainScreen() string {
	return l.url("/app/", "", "#")
}
