// quotedPrintable implements an encoder for the quoted-printable
// wire format defined in RFC 2045.
package quote_printable

import (
	"bytes"
	"io"
	"strings"

	"github.com/haiyiyun/log"
)

var TransferEncodingType = "quoted-printable"

func Encode(src []byte) []byte {
	var dstBuilder []byte
	if len(src) == 0 {
		return dstBuilder
	}

	lineCounter := 0
	lineMax := 76
	for _, b := range src {
		encoded := quotedPrintableEncodeByte(b)
		lineCounter += len(encoded)
		if lineCounter > lineMax {
			dstBuilder = append(dstBuilder, []byte("=\r\n")...)
			lineCounter = len(encoded)
		}
		dstBuilder = append(dstBuilder, encoded...)
	}

	return dstBuilder
}

func EncodeToString(src []byte) string {
	return string(Encode(src))
}

func Decode(src []byte) []byte {
	var buf bytes.Buffer
	_, err := io.Copy(&buf, newQuotedPrintableReader(strings.NewReader(string(src))))
	if err != nil {
		log.Debug("<Decode> io.Copy Error:", err)
	}

	log.Debug("<Decode> String:", buf.String())
	return buf.Bytes()
}

func DecodeToString(src []byte) string {
	return string(Decode(src))
}
