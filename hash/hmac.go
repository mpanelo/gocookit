package hash

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"hash"
)

type Hmac struct {
	hmac hash.Hash
}

func NewHmac(key string) *Hmac {
	h := hmac.New(sha256.New, []byte(key))
	return &Hmac{h}
}

func (h *Hmac) Hash(token string) string {
	h.hmac.Reset()
	h.hmac.Write([]byte(token))
	b := h.hmac.Sum(nil)
	return base64.URLEncoding.EncodeToString(b)
}
