package send

import (
	"crypto/tls"
	"net"
	"net/smtp"
	"strings"

	"github.com/haiyiyun/mail/data"
	"github.com/haiyiyun/mail/predefined"
)

func SendMail(addr, mailIp string, m data.Mailer, skipVerify bool, a smtp.Auth, from string) error {
	h := m.Header()
	md := m.Data()
	tos := predefined.EmailRegexp.FindAllString(h.Get("To"), -1)
	mailFrom := strings.Trim(h.Get("Return-Path"), "<>")
	mailDomain := strings.Split(mailFrom, "@")[1]
	if from != "" {
		mailFrom = from
	}

	return sendMail(skipVerify, addr, mailIp, mailDomain, a, mailFrom, tos, md)
}

func dial(addr, ip string, skipVerify bool) (*smtp.Client, error) {
	var client *smtp.Client
	var nErr error
	host, port, _ := net.SplitHostPort(addr)
	var d net.Dialer

	if ip != "" && ip != "127.0.0.1" {
		d.LocalAddr = &net.TCPAddr{
			IP: net.ParseIP(ip),
		}
	}

	conn, err := d.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}

	if port != "25" {
		var config *tls.Config
		if skipVerify {
			config = &tls.Config{
				InsecureSkipVerify: true,
			}
		}

		tlsConn := tls.Client(conn, config)
		client, nErr = smtp.NewClient(tlsConn, host)
	} else {
		client, nErr = smtp.NewClient(conn, host)
	}

	return client, nErr
}

func sendMail(skipVerify bool, addr, mailIp, mailDomain string, a smtp.Auth, from string, to []string, msg []byte) error {
	c, dialErr := dial(addr, mailIp, skipVerify)
	if dialErr != nil {
		return dialErr
	}

	defer c.Close()
	if err := c.Hello(mailDomain); err != nil {
		return err
	}

	if ok, _ := c.Extension("STARTTLS"); ok {
		if skipVerify {
			config := &tls.Config{
				InsecureSkipVerify: true,
			}

			if err := c.StartTLS(config); err != nil {
				return err
			}
		} else {
			if err := c.StartTLS(nil); err != nil {
				return err
			}
		}
	}

	if a != nil {
		if ok, _ := c.Extension("AUTH"); ok {
			if err := c.Auth(a); err != nil {
				return err
			}
		}
	}

	if err := c.Mail(from); err != nil {
		return err
	}

	for _, addr := range to {
		if err := c.Rcpt(addr); err != nil {
			return err
		}
	}

	w, err := c.Data()
	if err != nil {
		return err
	}

	_, err = w.Write(msg)
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}

	return c.Quit()
}
