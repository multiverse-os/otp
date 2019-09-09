package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	otp "github.com/multiverse-os/otp"
)

func PrintToken(otpType, otpCode string) {
	switch otpType {
	case "HOTP":
		fmt.Println(otp.NewHOTP(otpCode).Generate())
	case "TOTP":
		fmt.Println(otp.NewTOTP(otpCode).Generate())
	}
}

// Primarily an example program until a password manager will be aseembled
// in addition this will be used to as an example for the README, because it is
// currently missing the TOTP example use case.
func main() {
	args := os.Args
	if len(args) < 2 {
		fmt.Println("[error] failed to provide arguments, try again using the following:\n")
		fmt.Println("    otp [hotp|totp] [seed_value]\n")
	} else {
		otpType := strings.ToUpper(args[1])
		otpCode := args[2]
		if otpType == "HOTP" || otpType == "TOTP" {
			fmt.Println("Scramble Suit: OTP")
			fmt.Println("==========================")
			fmt.Println("Generating OTP type token:")

			PrintToken(otpType, otpCode)

			tick := time.Tick(20 * time.Second)
			for {
				select {
				case <-tick:
					PrintToken(otpType, otpCode)
				}
			}
		} else {
			fmt.Println("[error] failed to provide OTP type, try again using the following:\n")
			fmt.Println("    otp [hotp|totp] [seed_value]\n")
		}
	}
}
