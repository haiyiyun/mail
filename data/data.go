package data

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"net/textproto"
	"strings"

	"github.com/haiyiyun/mail/quote_printable"
)

type Sender string

func (s Sender) Encode(prefix, senderDomain string) string {
	ss := string(s)
	sss := strings.Split(ss, "@")
	replace := true
	if len(sss) > 1 {
		if sss[1] == senderDomain {
			replace = false
		}
	}

	if replace {
		ss = strings.Replace(ss, "@", "=", -1) + "@" + senderDomain
		if prefix != "" {
			ss = prefix + "-" + ss
		}
	}

	return ss
}

func (s Sender) Decode(senderDomain string) (string, string) {
	var prefix string
	ss := string(s)
	sss := strings.Split(ss, "@")

	replace := true
	if len(sss) > 1 {
		if sss[1] == senderDomain {
			replace = false
		}
	}

	if replace {
		if pos := strings.LastIndex(ss, "-"); pos > 0 {
			prefix = ss[:pos-1]
			ss = ss[pos+1:]
		}
		if strings.Contains(ss, "=") {
			ss = strings.Replace(sss[0], "=", "@", -1)
		}
	}

	return prefix, ss
}

var charset = "utf-8"

type QP string

func (s QP) Q() string {
	return "=?" + charset + "?Q?"
}

func (s QP) encodeQP(src []byte) string {
	ss := []string{}
	encStr := quote_printable.EncodeToString(src)
	if !strings.Contains(encStr, "=\r\n") {
		ss = append(ss, s.Q()+encStr+"?=")
	} else {
		for _, str := range strings.Split(encStr, "=\r\n") {
			ss = append(ss, s.Q()+str+"?=")
		}
	}

	return strings.Join(ss, "\r\n ")
}

func (s QP) EncodeQP() string {
	return s.encodeQP([]byte(s))

}

func (s QP) decodeQP(src []byte) string {
	ssrc := string(src)
	return quote_printable.DecodeToString([]byte(ssrc[len(s.Q()) : len(ssrc)-2]))
}

func (s QP) DecodeQP() string {
	if !strings.Contains(string(s), "\r\n ") {
		return s.decodeQP([]byte(string(s)))
	}

	ss := ""
	for _, str := range strings.Split(string(s), "\r\n ") {
		ss += s.decodeQP([]byte(str))
	}

	return ss
}

type Body struct {
	Plain string
	Html  string
	Ext   map[string]string
}

func (b *Body) QP() []byte {
	var crlf = "\r\n"
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	boundary := w.Boundary()
	fmt.Fprint(&buf, "Mime-Version: 1.0", crlf)
	fmt.Fprint(&buf, `Content-Type: multipart/alternative;`, crlf, ` boundary="`+boundary+`";`, crlf, ` charset=`+charset, crlf, crlf)
	if b.Plain != "" {
		mh := make(textproto.MIMEHeader)
		mh.Add("Content-Type", "text/plain;"+crlf+" charset="+charset)
		mh.Add("Content-Transfer-Encoding", quote_printable.TransferEncodingType)
		w.CreatePart(mh)
		buf.Write(quote_printable.Encode([]byte(b.Plain)))
		buf.WriteString(crlf + crlf)
	}

	if b.Html != "" {
		mh := make(textproto.MIMEHeader)
		mh.Add("Content-Type", "text/html;"+crlf+" charset="+charset)
		mh.Add("Content-Transfer-Encoding", quote_printable.TransferEncodingType)
		w.CreatePart(mh)
		buf.Write(quote_printable.Encode([]byte(b.Html)))
		buf.WriteString(crlf + crlf)
	}

	if b.Ext != nil || len(b.Ext) > 0 {
		for t, bd := range b.Ext {
			mh := make(textproto.MIMEHeader)
			mh.Add("Content-Type", "text/"+t+";"+crlf+" charset="+charset)
			mh.Add("Content-Transfer-Encoding", quote_printable.TransferEncodingType)
			w.CreatePart(mh)
			buf.Write(quote_printable.Encode([]byte(bd)))
			buf.WriteString(crlf + crlf)
		}
	}

	w.Close()

	return buf.Bytes()
}
