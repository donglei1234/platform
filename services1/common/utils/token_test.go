package utils

import (
	"fmt"
	"testing"
)

func TestDecodeToken(t *testing.T) {

	testcases := []string{
		"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MDAzMTE4MjUsInVpZCI6IjM5Mzc4NjU5MTYifQ._a6ygCNvgWm2FmrWoOPATqTzJn5yh4fYr5eDe6mm2RQ",
		"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MDAzMTQwMjksInVpZCI6IjE5NzIyMzcwNzQifQ.LwscTb25dSingPSJxnI05qsc0Uk8fMpb85Kjf56OXkI",
		"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MDAzMTQwNTYsInVpZCI6IjE0NzE1MDk4MiJ9.w1-jWqDc65zHQItBN3OsrLjV5AXU7A32QxRu1os-cd0",
		"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MDAzMTQwNzgsInVpZCI6IjUxNDQ5OTI2MyJ9.ALZTOCDzDcycrg32LOxESZy0T54rGbJft9ZyUbQJDIA",
		"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MDAzMTQxMDMsInVpZCI6IjM2NjA0OTc5NzkifQ.dc3x96m62SpeHDlFK23e9tqUOEyvmFJYJa0w8lp7UKA",
		"",
		"12345678-12345678-12345678",
		"ksajdklasjdklasjkd",
		"&^^%&&&@!&*#",
		"12345678",
	}
	for i, val := range testcases {
		uid, err := DecodeToken(val)
		if err != nil {
			t.Errorf("Test %d: decode token err:%v", i, err)
		}
		fmt.Println(uid)
	}
}
