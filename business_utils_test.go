package kit

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_IsIpV4Valid(t *testing.T) {
	tests := []struct {
		in    string
		valid bool
	}{
		{
			in:    "",
			valid: false,
		}, {
			in:    "invalid",
			valid: false,
		}, {
			in:    "0.0.0.0",
			valid: true,
		}, {
			in:    "10.10.20.4",
			valid: true,
		}, {
			in:    "0",
			valid: false,
		}, {
			in:    "0.0",
			valid: false,
		}, {
			in:    "0.0.0..0",
			valid: false,
		}, {
			in:    "1.1.1.1",
			valid: true,
		}, {
			in:    "0.0.0.0 ",
			valid: false,
		}, {
			in:    "255.255.255.255",
			valid: true,
		}, {
			in:    " 255.255.255.255",
			valid: false,
		}, {
			in:    "256.255.255.255",
			valid: false,
		}, {
			in:    "255.256.255.255",
			valid: false,
		}, {
			in:    "255.255.256.255",
			valid: false,
		}, {
			in:    "255.255.255.256",
			valid: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			assert.Equal(t, tt.valid, IsIpV4Valid(tt.in))
		})
	}
}

func Test_IsIpV6Valid(t *testing.T) {
	tests := []struct {
		in    string
		valid bool
	}{
		{
			in:    "",
			valid: false,
		}, {
			in:    "invalid",
			valid: false,
		}, {
			in:    "1:2:3:4:5:6:7:8",
			valid: true,
		}, {
			in:    "1::",
			valid: true,
		}, {
			in:    "1::8",
			valid: true,
		}, {
			in:    "1::7:8",
			valid: true,
		}, {
			in:    "1::6:7:8",
			valid: true,
		}, {
			in:    "1::5:6:7:8",
			valid: true,
		}, {
			in:    "1::4:5:6:7:8",
			valid: true,
		}, {
			in:    "1::3:4:5:6:7:8",
			valid: true,
		}, {
			in:    "::2:3:4:5:6:7:8",
			valid: true,
		}, {
			in:    "::8",
			valid: true,
		}, {
			in:    "1:2:3:4::6:7:8",
			valid: true,
		}, {
			in:    "::255.255.255.255",
			valid: true,
		}, {
			in:    "::ffff:255.255.255.255",
			valid: true,
		}, {
			in:    "::ffff:0:255.255.255.255",
			valid: true,
		}, {
			in:    "2001:db8:3:4::192.0.2.33",
			valid: true,
		}, {
			in:    "64:ff9b::192.0.2.33",
			valid: true,
		}, {
			in:    "2001:0db8:11a3:09d7:1f34:8a2e:07a0:765d",
			valid: true,
		}, {
			in:    " 64:ff9b::192.0.2.33",
			valid: false,
		}, {
			in:    "64:ff9b::192.0.2.33 ",
			valid: false,
		}, {
			in:    "::iiii:255.255.255.255",
			valid: false,
		}, {
			in:    "2001:0db8:0000:::ff00:0042:8329",
			valid: false,
		}, {
			in:    "20018:0db8:0000::ff00:0042:8329",
			valid: false,
		}, {
			in:    "2001:0db8:0000::ff00:80042:8329",
			valid: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			assert.Equal(t, tt.valid, IsIpV6Valid(tt.in))
		})
	}
}

func Test_IsUrlValid(t *testing.T) {
	tests := []struct {
		in  string
		out bool
	}{
		{out: true, in: "http://www.foo.com"},
		{out: true, in: "http://www.foo.group"},
		{out: true, in: "http://www.foo.group.com"},
		{out: true, in: "http://www.foo.group.local.com"},
		{out: true, in: "https://www.foo.com"},
		{out: true, in: "https://www.foo.group"},
		{out: true, in: "https://www.foo.group.com"},
		{out: true, in: "https://www.foo.group.local.com"},
		{out: true, in: "www.foo.com"},
		{out: true, in: "www.foo.group"},
		{out: true, in: "www.foo.group.com"},
		{out: true, in: "www.foo.group.local.com"},
		{out: true, in: "foo.com"},
		{out: true, in: "foo.group"},
		{out: true, in: "foo.group.com"},
		{out: true, in: "group.local.com"},
		{out: false, in: "httpd://www.foo.com"},
		{out: false, in: "httpd://www.foo.group"},
		{out: false, in: "httpd://www.foo.group.com"},
		{out: false, in: "httpd://www.foo.group.local.com"},
		{out: true, in: "http://www.foo.com/local"},
		{out: true, in: "http://www.foo.group/local/group/data"},
		{out: true, in: "http://www.foo.group.com?page=local"},
		{out: true, in: "http://www.foo.group.local.com/page.js"},
		{out: false, in: "://www.foo.com/local"},
		{out: false, in: "http//www.foo.group/local/group/data"},
		{out: false, in: "http:/www.foo.group.com?page=local"},
		{out: false, in: "www../local"},
		{out: false, in: ""},
		{out: false, in: "local"},
		{out: false, in: "."},
		{out: false, in: ".com"},
		{out: false, in: "y."},
		{out: false, in: "http://y."},
		{out: false, in: "http://.com"},
		{out: true, in: "y.com"},
		{out: true, in: "99.com"},
		{out: true, in: "http://localhost:9999/page.js"},
		{out: true, in: "http://localhost:9999"},
	}
	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			assert.Equal(t, tt.out, IsUrlValid(tt.in))
		})
	}
}

func Test_IsEmailValid(t *testing.T) {
	tests := []struct {
		name string
		in   string
		out  bool
	}{
		{name: "valid", in: "test@test.com", out: true},
		{name: "with dot in domain", in: "test@test.test.com", out: true},
		{name: "with dot", in: "test.test@test.com", out: true},
		{name: "with plus and minus signs", in: "test.test+test-test@test.com", out: true},
		{name: "with quotes around username", in: "\"test\"@test.com", out: true},
		{name: "with quotes and @ inside username", in: "\"test.@.test\"@test.com", out: true},
		{name: "with allowed special symbols", in: "#!$%&'*+-/=?^_`{}|~@test.com", out: true},
		{name: "cyrillic", in: "тест@тест.рф", out: true},
		{name: "cyrillic username with latin domain", in: "тест@test.com", out: true},
		{name: "latin username with cyrillic domain", in: "test@тест.рф", out: true},
		{name: "too long username", in: "abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklm@test.com", out: true},
		{name: "not valid tld", in: "3@test.abcdefghij", out: true},
		{name: "without domain", in: "test@test", out: true},

		{name: "double @", in: "te@st@test.com", out: false},
		{name: "with open brace", in: "tes(t@test.com", out: false},
		{name: "with close brace", in: "tes)t@test.com", out: false},
		{name: "with <", in: "tes<t@test.com", out: false},
		{name: "with >", in: "tes>t@test.com", out: false},
		{name: "with comma", in: "tes,t@test.com", out: false},
		{name: "with colon", in: "tes:t@test.com", out: false},
		{name: "with semicolon", in: "tes;t@test.com", out: false},
		{name: "empty email", in: "", out: false},
		{name: "dot at the end", in: "test@test.com.", out: false},
		{name: "dot at the end of username", in: "test.@test.com", out: false},
		{name: "dot at the beginning of username", in: ".test@test.com", out: false},
		{name: "too long hostname", in: "1@abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijfg.com", out: false},
		{name: "too long domain", in: "2@test.abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijfghabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijfghabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijfghabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzab", out: false},
		{name: "without @", in: "testtest.com", out: false},
		{name: "without username", in: "@test.com", out: false},
		{name: "with multiple @", in: "A@b@c@test.com", out: false},
		{name: "with quotes inside username", in: "just\"not\"right@test.com", out: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.out, IsEmailValid(tt.in))
		})
	}
}

func Test_IsTelegramValid(t *testing.T) {
	tests := []struct {
		name string
		in   string
		out  bool
	}{
		{name: "valid with https", in: "https://t.me/something", out: true},
		{name: "valid with www", in: "www.t.me/something", out: true},
		{name: "valid without protocol", in: "telegram.me/23648724something", out: true},
		{name: "not valid with symbols", in: "t.me/#$%", out: false},
		{name: "not valid with site", in: "q.me/123", out: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.out, IsTelegramChannelValid(tt.in))
		})
	}
}

func Test_ExtractExtension(t *testing.T) {
	tests := []struct {
		url      string
		expected string
		err      bool
	}{
		{
			url:      "Smt smt",
			expected: "",
			err:      false,
		}, {
			url:      "https://example.com/secure/Dashboard.jspa",
			expected: "jspa",
			err:      false,
		}, {
			url:      "https://example.com/secure/Dashboard.jspa?tt=12&ss=image.jpg",
			expected: "jspa",
			err:      false,
		}, {
			url:      "Dashboard.jspa?tt=12&ss=jpg",
			expected: "jspa",
			err:      false,
		}, {
			url:      "jspa?tt=12&ss=jpg",
			expected: "",
			err:      false,
		}, {
			url:      "tt=12&ss=image.jpg",
			expected: "",
			err:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.url, func(t *testing.T) {
			ext, err := ExtractUrlExtension(tt.url)
			assert.Equal(t, tt.expected, ext)
			if tt.err {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
