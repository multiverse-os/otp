package otp

import (
	"crypto/sha1"
	"hash"
	"time"
)

// otpauth://totp/Company:joe_example@gmail.com?secret=[...]&issuer=Company
// Totp is a struct holding the details for a time based hmac-sha1 otp
type TOTP struct {
	seed        string
	time        time.Time
	tokenLength int
	window      int
	windowSize  int
	base32      bool
	encoding    func() hash.Hash
}

func NewTOTP(seed string) TOTP {
	return TOTP{
		seed:        seed,
		time:        time.Now(),
		tokenLength: 6,
		window:      30,
		windowSize:  2,
		base32:      true,
		encoding:    sha1.New,
	}
}

func (self TOTP) Time(totpTime time.Time) TOTP {
	self.time = totpTime
	return self
}

func (self TOTP) TokenLength(tokenLength int) TOTP {
	self.tokenLength = tokenLength
	return self
}

func (self TOTP) Window(window int) TOTP {
	self.window = window
	return self
}
func (self TOTP) WindowSize(windowSize int) TOTP {
	self.windowSize = windowSize
	return self
}

func (self TOTP) Base32(base32 bool) TOTP {
	self.base32 = base32
	return self
}

func (self TOTP) Encoding(encoding func() hash.Hash) TOTP {
	self.encoding = encoding
	return self
}

func (totp TOTP) Generate() string {
	totpWindow := int(totp.time.Unix()) / totp.window
	hotp := HOTP{
		seed:        totp.seed,
		counter:     totpWindow,
		tokenLength: totp.tokenLength,
		window:      totp.window,
		base32:      totp.base32,
		encoding:    totp.encoding,
	}
	return hotp.Generate()
}

func (totp TOTP) Check(otp string) bool {
	otpTime := totp.time
	for i := -1 * totp.windowSize; i < totp.windowSize; i++ {
		totp.time = otpTime.Add(time.Second * time.Duration(i*totp.window))
		if totp.Generate() == otp {
			return true
		}
	}
	return false
}
