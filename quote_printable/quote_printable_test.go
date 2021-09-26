package quote_printable

import (
	"testing"

	"github.com/haiyiyun/log"
)

func TestQP(t *testing.T) {
	log.SetLevel(log.LEVEL_DISABLE)
	srcQP := "=E4=BA=91=3F=E9=82=AE"
	src := "云?邮"
	qpd := Encode([]byte(src))
	if srcQP != string(qpd) {
		t.Error("invalid encode in:", string(qpd), " want:", srcQP)
	}
	qoenStr := DecodeToString(qpd)
	if src != qoenStr {
		t.Error("invalid decode in:", qoenStr, " want:", src)
	}
}
