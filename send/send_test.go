package send

import (
	"net/http"
	"strings"
	"testing"

	"github.com/haiyiyun/log"
	"github.com/haiyiyun/mail/data"
)

/**
日志格式

1. sender_domain`~~`local_ip`~~`receive_domain`~~`remote_mta_ip`~~`MESSAGE-ID(ObjectId+8位批次@发送域名)`~~`form`~~`to`~~`subject`~~`code`~~`msg
*/

func TestSend(t *testing.T) {
	log.SetLevel(log.LEVEL_INFO)
	localIp := "192.168.8.8"
	senderDomain := "cloudmail.cn"
	fromDomain := "cloudmail.cc"
	from := "330222918@" + fromDomain
	receiveDomain := "mail-tester.com"
	to := "web-4xrobx@" + receiveDomain
	subject := "cloudmail -- 云=邮 -中?=国- 云?邮 -- send test"
	ip, mxErr := RandomIP(receiveDomain)
	if mxErr != nil {
		t.Error(mxErr)
		return
	}

	body := &data.Body{
		Plain: "text 云邮",
		Html: `<p style="WHITE-SPACE: normal; TEXT-TRANSFORM: none; WORD-SPACING: 0px; COLOR: rgb(0,0,0); FONT: 14px/23px Helvetica, 'Microsoft Yahei', verdana; LETTER-SPACING: normal; TEXT-INDENT: 0px; -webkit-text-stroke-width: 0px" lucida="" grande="" verdana="" margin:="" padding:=""><font style="LINE-HEIGHT: 30px" color="#ff0000" size="4" face="microsoft yahei"><span style="LINE-HEIGHT: 27px"><b>免费试发 
100封，接有诚意的卖家，本软件98%进QQ收件箱</b></span></font></p>
<p style="WHITE-SPACE: normal; TEXT-TRANSFORM: none; WORD-SPACING: 0px; COLOR: rgb(0,0,0); FONT: 14px/23px Helvetica, 'Microsoft Yahei', verdana; LETTER-SPACING: normal; TEXT-INDENT: 0px; -webkit-text-stroke-width: 0px" lucida="" grande="" verdana="" margin:="" padding:=""><font style="LINE-HEIGHT: 30px" color="#ff0000" size="4" face="microsoft yahei"><span style="LINE-HEIGHT: 27px"><b>由于最近接单的人数太多，现在每天只发20万封，先到先推广，邮件营销不再没有效果</b></span></font></p>
<p style="WHITE-SPACE: normal; TEXT-TRANSFORM: none; WORD-SPACING: 0px; COLOR: rgb(0,0,0); FONT: 14px/23px Helvetica, 'Microsoft Yahei', verdana; LETTER-SPACING: normal; TEXT-INDENT: 0px; -webkit-text-stroke-width: 0px" lucida="" grande="" verdana="" margin:="" padding:=""><font style="LINE-HEIGHT: 30px" color="#ff0000" size="4" face="microsoft yahei"><span style="LINE-HEIGHT: 27px"><b>我们软件的优势</b></span></font></p>
<p style="WHITE-SPACE: normal; TEXT-TRANSFORM: none; WORD-SPACING: 0px; COLOR: rgb(0,0,0); FONT: 14px/23px Helvetica, 'Microsoft Yahei', verdana; LETTER-SPACING: normal; TEXT-INDENT: 0px; -webkit-text-stroke-width: 0px" lucida="" grande="" verdana="" margin:="" padding:=""><font style="LINE-HEIGHT: 30px" color="#ff0000" size="4" face="microsoft yahei"><span style="LINE-HEIGHT: 27px"><b>1、</b></span></font><span style="FONT-WEIGHT: bold; LINE-HEIGHT: normal; BACKGROUND-COLOR: rgb(239,245,251)">投递技术，不用发件箱，不用变量，不用换IP，不用担心敏感词 
98%不进垃圾</span></p>
<p style="WHITE-SPACE: normal; TEXT-TRANSFORM: none; WORD-SPACING: 0px; COLOR: rgb(0,0,0); FONT: 14px/23px Helvetica, 'Microsoft Yahei', verdana; LETTER-SPACING: normal; TEXT-INDENT: 0px; -webkit-text-stroke-width: 0px" lucida="" grande="" verdana="" margin:="" padding:=""></p>
<p style="WHITE-SPACE: normal; TEXT-TRANSFORM: none; WORD-SPACING: 0px; COLOR: rgb(0,0,0); TEXT-ALIGN: left; FONT: 14px/24px Helvetica, 'Microsoft Yahei', verdana; LETTER-SPACING: normal; TEXT-INDENT: 0px; -webkit-text-stroke-width: 0px" align="center" lucida="lucida" grande="grande" verdana="verdana"><font style="LINE-HEIGHT: 25px" color="#c43cc4" size="3">2、ESP白名单确认，保证邮件进收件箱而不是垃圾箱</font></p>
<p style="WHITE-SPACE: normal; TEXT-TRANSFORM: none; WORD-SPACING: 0px; COLOR: rgb(0,0,0); TEXT-ALIGN: left; FONT: 14px/24px Helvetica, 'Microsoft Yahei', verdana; LETTER-SPACING: normal; TEXT-INDENT: 0px; -webkit-text-stroke-width: 0px" align="center" lucida="lucida" grande="grande" verdana="verdana"><font style="LINE-HEIGHT: 25px" color="#c43cc4" size="3">3、无需发件箱可以日发送13万封邮件</font></p>
<p style="WHITE-SPACE: normal; TEXT-TRANSFORM: none; WORD-SPACING: 0px; COLOR: rgb(0,0,0); TEXT-ALIGN: left; FONT: 14px/24px Helvetica, 'Microsoft Yahei', verdana; LETTER-SPACING: normal; TEXT-INDENT: 0px; -webkit-text-stroke-width: 0px" align="center" lucida="lucida" grande="grande" verdana="verdana"><font style="LINE-HEIGHT: 25px" color="#c43cc4" size="3">4、实时邮件监控，监控邮件到达率、打开率等，并以图文形式展现。</font></p>
<p style="WHITE-SPACE: normal; TEXT-TRANSFORM: none; WORD-SPACING: 0px; COLOR: rgb(0,0,0); TEXT-ALIGN: left; FONT: 14px/24px Helvetica, 'Microsoft Yahei', verdana; LETTER-SPACING: normal; TEXT-INDENT: 0px; -webkit-text-stroke-width: 0px" align="center" lucida="lucida" grande="grande" verdana="verdana"><font style="LINE-HEIGHT: 25px" color="#c43cc4" size="3">5、打造自己的企业发件箱，让各大服务商默认你，进入收件箱。</font></p>
<p style="WHITE-SPACE: normal; TEXT-TRANSFORM: none; WORD-SPACING: 0px; COLOR: rgb(0,0,0); TEXT-ALIGN: left; FONT: 14px/24px Helvetica, 'Microsoft Yahei', verdana; LETTER-SPACING: normal; TEXT-INDENT: 0px; -webkit-text-stroke-width: 0px" align="center" lucida="lucida" grande="grande" verdana="verdana"><font style="LINE-HEIGHT: 25px" color="#c43cc4" size="3">6、模拟企业发件箱，我可以让发件人变成百度的，当然发件邮箱也是百度的。（此功能支持国能所有服务商）</font></p>
<p style="WHITE-SPACE: normal; TEXT-TRANSFORM: none; WORD-SPACING: 0px; COLOR: rgb(0,0,0); TEXT-ALIGN: left; FONT: 14px/24px Helvetica, 'Microsoft Yahei', verdana; LETTER-SPACING: normal; TEXT-INDENT: 0px; -webkit-text-stroke-width: 0px" align="center" lucida="lucida" grande="grande" verdana="verdana"><font style="LINE-HEIGHT: 25px" color="#c43cc4" size="3">以前我们在拼命的找别人的大站协议去发邮件或者是用amazon 
 ses现在我们打造自己的大站邮箱群发系统</font></p>
<p style="WHITE-SPACE: normal; TEXT-TRANSFORM: none; WORD-SPACING: 0px; COLOR: rgb(0,0,0); TEXT-ALIGN: left; FONT: 14px/23px Helvetica, 'Microsoft Yahei', verdana; LETTER-SPACING: normal; TEXT-INDENT: 0px; -webkit-text-stroke-width: 0px" align="center" lucida="lucida" grande="grande" verdana="verdana"><font color="#c43cc4" size="3"><span style="LINE-HEIGHT: 25px">软件价格</span></font></p>`,
	}
	header := make(http.Header)
	mailData := data.NewDefault("fdk45gfgfdio4545kgfjgf", senderDomain, from, to, subject, body, header)
	mailHeader := mailData.Header()
	logArr := []string{}
	logArr = append(logArr, mailHeader.Get("X-Sender-Domain"))
	logArr = append(logArr, localIp)
	logArr = append(logArr, receiveDomain)
	logArr = append(logArr, ip)
	logArr = append(logArr, mailHeader.Get("Message-Id"))
	logArr = append(logArr, mailHeader.Get("From"))
	logArr = append(logArr, mailHeader.Get("To"))
	logArr = append(logArr, mailData.Subject())
	err := SendMail(ip+":25", "127.0.0.1", mailData, false, nil, "")
	if err != nil {
		code, msg := getCodeMsg(err.Error())
		logArr = append(logArr, code)
		logArr = append(logArr, msg)
	} else {
		logArr = append(logArr, "250")
		logArr = append(logArr, "")
	}

	log.Info(strings.Join(logArr, "`~~`"))

	if err != nil {
		t.Error(err)
	}
}

func getCodeMsg(line string) (string, string) {
	if len(line) < 4 {
		return "", ""
	}

	code := line[:3]
	msg := line[4:]
	return code, msg
}
