package data

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/haiyiyun/log"
	"github.com/haiyiyun/mail/quote_printable"
)

var (
	senderDomain = "go.haiyiyun.cn"
	from         = "ldf@haiyiyu.cn"
	to           = "fook@haiyiyun.cn"
	subject      = "cloudmail -- 云=邮 -中?=国- 云?邮 -- send test"
	plant        = "text 云邮"
	html         = "html 云邮"
	subjectTests = [][2]string{
		[2]string{
			"云邮", "=?utf-8?Q?=E4=BA=91=E9=82=AE?=",
		},
		[2]string{
			"中?国", "=?utf-8?Q?=E4=B8=AD=3F=E5=9B=BD?=",
		},
		[2]string{
			"cloudmail -- 云=邮 -中?=国- 云?邮 -- send test", "=?utf-8?Q?cloudmail -- =E4=BA=91=3D=E9=82=AE -=E4=B8=AD=3F=3D=E5=9B=BD- =E4=BA=91=3F?=\r\n =?utf-8?Q?=E9=82=AE -- send test?=",
		},
	}
	dataTests = [][]byte{
		[]byte("To: " + to),
		[]byte("From: " + from),
		[]byte("Sender: " + Sender(from).Encode("", senderDomain)),
		[]byte("X-Priority: 3"),
		[]byte("Subject: " + QP(subject).EncodeQP()),
		quote_printable.Encode([]byte(plant)),
		quote_printable.Encode([]byte(html)),
	}
)

func TestSubject(t *testing.T) {
	log.SetLevel(log.LEVEL_DISABLE)
	for _, test := range subjectTests {
		if enSub := QP(test[0]).EncodeQP(); enSub != test[1] {
			t.Error("invalid Subject Encode:", enSub)
		}

		if deSub := QP(test[1]).DecodeQP(); deSub != test[0] {
			t.Error("invalid Subject Decode:", deSub)
		}
	}
}

func testDataer(d Dataer) []byte {
	return d.Data()
}

func TestData(t *testing.T) {
	log.SetLevel(log.LEVEL_DISABLE)
	headers := make(http.Header)
	headers.Add("X-Priority", "3")

	md := NewDefault("fdk45gfgfdio4545kgfjgf", senderDomain, from, to, subject, &Body{
		Plain: plant,
		Html:  html,
	}, headers)
	bt := testDataer(md)
	for _, test := range dataTests {
		if !bytes.Contains(bt, test) {
			t.Error("Test:", string(test))
			t.Error("Data:", string(md.Data()))
			return
		}
	}
}
