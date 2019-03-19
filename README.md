<img src="https://avatars2.githubusercontent.com/u/24763891?s=400&u=c1150e7da5667f47159d433d8e49dad99a364f5f&v=4"  width="256px" height="256px" align="right" alt="Multiverse OS Logo">

## Scramble ID: OTP (HOTP, and TOTP) Library
**URL** [multiverse-os.org](https://multiverse-os.org)

A library to satisfy the requirements of Scramble ID, it provides by default
HOTP settings that satisfy the google authenticator protocol, while being
customizable enough to satisfy the needs of the more esoteric parts of scramble
ID and the software that uses it.


### Usage
Usage is done by chaining any customization onto the intialization function:

The simplest version, that will work with google authenicator is:

```
  hotp := otp.NewHOTP("seed-value")
```

And customization is done with chained functions:

```
  hotp := otp.NewHOTP("seed-value").Encoding(sha256.New)
```
