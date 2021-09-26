package send

import (
	"testing"

	"github.com/haiyiyun/log"
)

func TestIp(t *testing.T) {
	log.SetLevel(log.LEVEL_DISABLE)
	if _, err := RandomIP("qq.com"); err != nil {
		t.Error(err)
	}
}
