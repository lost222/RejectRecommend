package Cro

import (
	"crypto/sha256"
	"fmt"
	"io"
)

func SaltCro(password string) string{
	h := sha256.New()
	io.WriteString(h, password)

	//pwmd5等于e10adc3949ba59abbe56e057f20f883e
	pw :=fmt.Sprintf("%x", h.Sum(nil))

	//指定两个 salt： salt1 = @#$%   salt2 = ^&*()
	salt1 := "@#$%"
	salt2 := "^&*()"

	//salt1+salt2+MD5拼接
	io.WriteString(h, salt1)
	io.WriteString(h, salt2)
	io.WriteString(h, pw)

	last :=fmt.Sprintf("%x", h.Sum(nil))
	return last
}