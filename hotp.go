package otp

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"fmt"
	"hash"
	"math"
)

// otpauth://totp/Company:joe_example@gmail.com?secret=[...]&issuer=Company
type HOTP struct {
	seed        string
	window      int
	counter     int
	tokenLength int
	base32      bool
	encoding    func() hash.Hash
}

func NewHOTP(seed string) HOTP {
	return HOTP{
		seed:        seed,
		encoding:    sha1.New,
		window:      5,
		tokenLength: 6,
		base32:      true,
	}
}

func (self HOTP) Base32(base32 bool) HOTP {
	self.base32 = base32
	return self
}

func (self HOTP) Counter(counter int) HOTP {
	self.counter = counter
	return self
}

func (self HOTP) TokenLength(tokenLength int) HOTP {
	self.tokenLength = tokenLength
	return self
}

func (self HOTP) Window(window int) HOTP {
	self.window = window
	return self
}

func (self HOTP) Encoding(encoding func() hash.Hash) HOTP {
	self.encoding = encoding
	return self
}

func (self HOTP) Seed() []byte {
	if self.base32 {
		encodedSeed, _ := base32.StdEncoding.DecodeString(self.seed)
		return encodedSeed
	} else {
		return []byte(self.seed)
	}
}

func (self HOTP) HMAC() []byte {
	hash := hmac.New(self.encoding, self.Seed())
	hash.Write([]byte(counterToBytes(self.counter)))
	return hash.Sum(nil)
}

func (self HOTP) Generate() string {
	otp := truncate(self.HMAC()) % int(math.Pow10(self.tokenLength))
	return fmt.Sprintf(fmt.Sprintf("%%0%dd", self.tokenLength), otp)
}

func (self HOTP) Check(otp string) (bool, int) {
	for i := 0; i < self.window; i++ {
		o := self.Generate()
		if o == otp {
			return true, int(self.counter)
		}
		self.counter++
	}
	return false, 0
}

func (self HOTP) Sync(otp1 string, otp2 string) (bool, int) {
	self.window = 100
	v, i := self.Check(otp1)
	if !v {
		return false, 0
	}
	self.counter = self.counter + i + 1
	self.window = 1
	v2, i2 := self.Check(otp2)
	if v2 {
		return true, i2 + 1
	}
	return false, 0
}

func truncate(hash []byte) int {
	offset := int(hash[len(hash)-1] & 0xf)
	return ((int(hash[offset]) & 0x7f) << 24) |
		((int(hash[offset+1] & 0xff)) << 16) |
		((int(hash[offset+2] & 0xff)) << 8) |
		(int(hash[offset+3]) & 0xff)
}

func counterToBytes(counter int) (text []byte) {
	text = make([]byte, 8)
	for i := (len(text) - 1); i >= 0; i-- {
		text[i] = byte(counter & 0xff)
		counter = counter >> 8
	}
	return text
}
