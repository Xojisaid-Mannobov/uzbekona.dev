// Package slug o'zbekcha (lotin/kirill) matndan URL uchun xavfsiz slug yasaydi.
package slug

import (
	"regexp"
	"strings"
)

var cyrillic = map[rune]string{
	'а': "a", 'б': "b", 'в': "v", 'г': "g", 'д': "d", 'е': "e", 'ё': "yo", 'ж': "j", 'з': "z",
	'и': "i", 'й': "y", 'к': "k", 'л': "l", 'м': "m", 'н': "n", 'о': "o", 'п': "p", 'р': "r",
	'с': "s", 'т': "t", 'у': "u", 'ф': "f", 'х': "x", 'ц': "ts", 'ч': "ch", 'ш': "sh", 'щ': "sh",
	'ъ': "", 'ы': "i", 'ь': "", 'э': "e", 'ю': "yu", 'я': "ya", 'ў': "o", 'қ': "q", 'ғ': "g", 'ҳ': "h",
}

var nonAlnum = regexp.MustCompile(`[^a-z0-9]+`)

// Valid — slug formati: kichik harf, raqam va chiziqcha.
var Valid = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)

// Make matnni slugga aylantiradi: "G‘oyalar va O‘sish" → "goyalar-va-osish".
func Make(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))

	var b strings.Builder
	for _, r := range s {
		if v, ok := cyrillic[r]; ok {
			b.WriteString(v)
			continue
		}
		switch r {
		// o‘, g‘ dagi belgilar va tutuq belgisi tashlab yuboriladi
		case '‘', '’', '\'', '`', 'ʻ', 'ʼ':
			continue
		}
		b.WriteRune(r)
	}

	out := nonAlnum.ReplaceAllString(b.String(), "-")
	out = strings.Trim(out, "-")
	if len(out) > 120 {
		out = strings.Trim(out[:120], "-")
	}
	return out
}
