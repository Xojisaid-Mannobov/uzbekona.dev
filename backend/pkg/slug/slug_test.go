package slug

import "testing"

func TestMake(t *testing.T) {
	cases := map[string]string{
		"G‘oyalarni ishlaydigan mahsulotlarga": "goyalarni-ishlaydigan-mahsulotlarga",
		"  Telegram tizimlari  ":               "telegram-tizimlari",
		"AI & LLM integratsiyalar":             "ai-llm-integratsiyalar",
		"Ўзбекона дев":                         "ozbekona-dev",
		"Launch’dan keyin":                     "launchdan-keyin",
		"---":                                  "",
	}
	for in, want := range cases {
		if got := Make(in); got != want {
			t.Errorf("Make(%q) = %q, kutilgan %q", in, got, want)
		}
	}
}

func TestValid(t *testing.T) {
	for _, s := range []string{"kuaf", "web-platformalar", "a1-b2"} {
		if !Valid.MatchString(s) {
			t.Errorf("%q valid bo‘lishi kerak", s)
		}
	}
	for _, s := range []string{"Kuaf", "web--x", "-a", "a b", ""} {
		if Valid.MatchString(s) {
			t.Errorf("%q valid bo‘lmasligi kerak", s)
		}
	}
}
