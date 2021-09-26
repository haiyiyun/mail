package data

import (
	"bytes"
	"net/http"
	"strings"
	"time"

	"github.com/haiyiyun/log"
	"github.com/haiyiyun/validator"
)

type Default struct {
	H http.Header
	S string
	B *Body
}

func NewDefault(receivedId, senderDomain, from, to, subject string, body *Body, hs http.Header) *Default {
	//防止hs里面会相关header的值，所以使用set来操作，使有的话，覆盖掉
	hs.Set("Date", time.Now().Format(time.RFC1123Z))
	hs.Set("Message-Id", "<"+receivedId+"@"+senderDomain+">")
	fromAddr := validator.EmailRegexp.FindString(from)
	if pos := strings.Index(from, " <"); pos != -1 {
		hs.Set("From", QP(from[:pos]).EncodeQP()+from[pos:])
	} else {
		hs.Set("From", fromAddr)
	}

	toAddr := validator.EmailRegexp.FindString(to)
	if pos := strings.Index(to, " <"); pos != -1 {
		hs.Set("To", QP(to[:pos]).EncodeQP()+to[pos:])
	} else {
		hs.Set("To", toAddr)
	}

	hs.Set("Sender", Sender(fromAddr).Encode("", senderDomain))
	hs.Set("Return-Path", "<"+strings.Replace(toAddr, "@", "=", -1)+"@"+senderDomain+">")
	return &Default{
		H: hs,
		S: subject,
		B: body,
	}
}

func (d *Default) Date() string {
	return d.H.Get("Date")
}

func (d *Default) MessageId() string {
	return d.H.Get("Message-Id")
}

func (d *Default) Sender() string {
	return d.H.Get("Sender")
}

func (d *Default) Subject() string {
	return d.S
}

func (d *Default) Header() http.Header {
	return d.H
}

func (d *Default) HeaderByte() []byte {
	var buf bytes.Buffer
	if err := d.H.Write(&buf); err != nil {
		log.Error("<Header> Error:", err)
	}

	return buf.Bytes()
}

func (d *Default) Body() []byte {
	return d.B.QP()
}

func (d *Default) Data() []byte {
	var buf bytes.Buffer
	header := d.HeaderByte()
	body := d.Body()
	buf.Write(header)
	buf.WriteString("Subject: " + QP(d.Subject()).EncodeQP() + "\r\n")
	buf.Write(body)
	return buf.Bytes()
}
