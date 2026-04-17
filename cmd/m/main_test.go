package main

import "testing"

func TestParseCount(t *testing.T) {
	cases := []struct {
		arg    string
		wantN  int
		wantOK bool
	}{
		{"-5", 5, true},
		{"-1", 1, true},
		{"-200", 200, true},
		{"-254", 254, true},
		{"-255", 255, true},
		{"-0", 0, true},
		{"-", 0, false},
		{"", 0, false},
		{"5", 0, false},
		{"--force", 0, false},
		{"-f", 0, false},
		{"-abc", 0, false},
		{"-5abc", 0, false},
		{"-5.0", 0, false},
		{"./-5", 0, false},
	}
	for _, c := range cases {
		n, ok := parseCount(c.arg)
		if ok != c.wantOK || n != c.wantN {
			t.Errorf("parseCount(%q) = (%d, %v), want (%d, %v)", c.arg, n, ok, c.wantN, c.wantOK)
		}
	}
}

func TestValidateCount(t *testing.T) {
	cases := []struct {
		n       int
		wantErr bool
	}{
		{1, false},
		{5, false},
		{200, false},
		{254, false},
		{0, true},
		{-5, true},
		{255, true},
		{999999, true},
	}
	for _, c := range cases {
		err := validateCount(c.n)
		if (err != nil) != c.wantErr {
			t.Errorf("validateCount(%d) err=%v, wantErr=%v", c.n, err, c.wantErr)
		}
	}
}
