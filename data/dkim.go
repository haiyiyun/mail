package data

import (
	"net/http"

	"github.com/haiyiyun/dkim"
	"github.com/haiyiyun/log"
)

type Dkim struct {
	*Default
	D string
	S string
	K string
}

func NewDkim(selector, dkimPrivateKey, receivedId, senderDomain, from, to, subject string, body *Body, hs http.Header) *Dkim {
	return &Dkim{
		Default: NewDefault(receivedId, senderDomain, from, to, subject, body, hs),
		D:       senderDomain,
		S:       selector,
		K:       dkimPrivateKey,
	}
}

func (d *Dkim) Data() []byte {
	mailData := d.Default.Data()
	conf, err := dkim.NewConf(d.D, d.S)
	if err != nil {
		log.Error("Dkim NewConf error:", err)
		return mailData
	}

	dkimNew, err := dkim.New(conf, []byte(d.K))
	if err != nil {
		log.Error("Dkim New error:", err)
		return mailData
	}

	signByte, err := dkimNew.Sign(mailData)
	if err != nil {
		log.Error("dkimNew Sign error:", err)
		return mailData
	}

	return signByte
}
