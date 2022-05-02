package logger

import "testing"

func init() {
	Setup(&Config{
		Output:       "std",
		CipherKey:    "0123456789abcdef",
		CipherFields: []string{"secret", "mobile", "phone", "password"},
	})
}

func TestMatchReplace(t *testing.T) {
	cases := []struct {
		input string
		want  string
	}{
		{
			input: `{"more":"{\"phone\":\"13800138000\"}","password":"123456"}`,
			want:  `{"more":"{\"phone\":\"8fIrGTYAQzFR+tN2Wt0yDQ==\"}","password":"D89D/OvFrItmxI9ct8JbAg=="}`,
		}, {
			input: `{"more":"{\"cellPhone\":\"13800138000\"}","user_password":"123456"}`,
			want:  `{"more":"{\"cellPhone\":\"8fIrGTYAQzFR+tN2Wt0yDQ==\"}","user_password":"D89D/OvFrItmxI9ct8JbAg=="}`,
		},
		{
			input: `{"more":"{\"secret\":\"\"}","password":"","mobile_phone":"13800138000"}`,
			want:  `{"more":"{\"secret\":\"\"}","password":"","mobile_phone":"8fIrGTYAQzFR+tN2Wt0yDQ=="}`,
		},
		{
			input: `{"more":"{\"secret\":\"secret\"}","password":"password"}`,
			want:  `{"more":"{\"secret\":\"iNo+4RE7nEnpk338CYGhcw==\"}","password":"R4lDIBO/32oRLZSjtsPrGQ=="}`,
		},
	}
	for i, c := range cases {
		if got := matchReplace([]byte(c.input)); string(got) != c.want {
			t.Errorf("(%v) want: %s\ngot:%s\n", i, c.want, got)
		}
	}
}
