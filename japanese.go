//go:build ja_JP

package charencoder

import (
	"bytes"
	"errors"
	"io/ioutil"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"

	"github.com/saintfish/chardet"
)

// Custom error
// Encoded messages
var (
	err_mac       = errors.New("Normalized from UTF-8-mac to UTF-8")
	err_shift_jis = errors.New("Encoded from ShiftJIS to UTF-8")
	err_euc_jp    = errors.New("Encoded from EUC-JP to UTF-8")
)

type Encode_error struct {
	err error
}

func (e *Encode_error) Error() string {
	return e.err.Error()
}

// return custom error or error
// https://qiita.com/uchiko/items/1810ddacd23fd4d3c934
func Encorder(b *[]byte) error {
	res, err := chardet.NewTextDetector().DetectBest(*b)
	if err != nil {
		return err
	}
	// log.Println("Charset", res.Charset)

	switch res.Charset {
	case "UTF-8":
		// NFC (Normalization Form Composition), from UTF-8-mac to UTF-8
		nfc := norm.NFC.Bytes(*b)
		if bytes.Equal(*b, nfc) {
			nfc = nil
			return nil
		}
		*b = nfc
		nfc = nil
		return &Encode_error{err: err_mac}
	case "Shift_JIS":
		if *b, err = ioutil.ReadAll(transform.NewReader(bytes.NewReader(*b),
			japanese.ShiftJIS.NewDecoder())); err == nil {
			return &Encode_error{err: err_shift_jis}
			// return errors.New("!")
		}
	case "EUC-JP":
		if *b, err = ioutil.ReadAll(transform.NewReader(bytes.NewReader(*b),
			japanese.EUCJP.NewDecoder())); err == nil {
			return &Encode_error{err: err_euc_jp}
		}
	}
	return err
}

// golang struct interface is too free
//
//	It's like putting water in a colander
func IsEncoded(err error) bool {
	if e, ok := err.(*Encode_error); ok {
		switch e.err {
		case err_mac, err_shift_jis, err_euc_jp:
			return true
		}
	}
	return false
}
