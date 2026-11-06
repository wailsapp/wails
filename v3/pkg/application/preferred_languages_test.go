package application

import (
	"reflect"
	"testing"
)

func TestPosixPreferredLanguages(t *testing.T) {
	tests := []struct {
		name string
		env  map[string]string
		want []string
	}{
		{"LANG with codeset", map[string]string{"LANG": "nb_NO.UTF-8"}, []string{"nb-NO"}},
		{"script modifier", map[string]string{"LANG": "sr_RS@latin"}, []string{"sr-Latn-RS"}},
		{"script modifier after codeset", map[string]string{"LANG": "sr_RS.UTF-8@latin"}, []string{"sr-Latn-RS"}},
		{"other modifier dropped", map[string]string{"LANG": "de_DE@euro"}, []string{"de-DE"}},
		{"nothing set", map[string]string{}, []string{}},
		{"LANG=C", map[string]string{"LANG": "C"}, []string{}},
		{"LANG=POSIX", map[string]string{"LANG": "POSIX"}, []string{}},
		{"LANG=C.UTF-8", map[string]string{"LANG": "C.UTF-8"}, []string{}},
		{"LANGUAGE before the locale", map[string]string{"LANGUAGE": "nb:en", "LANG": "en_US.UTF-8"}, []string{"nb", "en", "en-US"}},
		{"LANGUAGE in POSIX form", map[string]string{"LANGUAGE": "nb_NO:nn_NO", "LANG": "nb_NO.UTF-8"}, []string{"nb-NO", "nn-NO"}},
		{"LANGUAGE with empty entries", map[string]string{"LANGUAGE": ":nb::", "LANG": "en_GB"}, []string{"nb", "en-GB"}},
		{"LC_ALL=C ignores LANGUAGE", map[string]string{"LC_ALL": "C", "LANGUAGE": "nb"}, []string{}},
		{"LANG=C.UTF-8 ignores LANGUAGE", map[string]string{"LANG": "C.UTF-8", "LANGUAGE": "nb"}, []string{}},
		{"LANGUAGE alone is ignored", map[string]string{"LANGUAGE": "nb"}, []string{}},
		{"LC_ALL wins", map[string]string{"LC_ALL": "en_GB", "LC_MESSAGES": "nb_NO", "LANG": "de_DE"}, []string{"en-GB"}},
		{"LC_MESSAGES before LANG", map[string]string{"LC_MESSAGES": "nb_NO", "LANG": "de_DE"}, []string{"nb-NO"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			getenv := func(key string) string { return tt.env[key] }
			got := normaliseLanguageTags(posixPreferredLanguages(getenv))
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestNormaliseLanguageTags(t *testing.T) {
	tests := []struct {
		name string
		raw  []string
		want []string
	}{
		{"canonical case, duplicates and junk", []string{"zh-hans-cn", "nb-NO", "nb-NO", "!!"}, []string{"zh-Hans-CN", "nb-NO"}},
		{"duplicates in different spellings", []string{"nb_NO", "nb-no", "en"}, []string{"nb-NO", "en"}},
		{"undetermined dropped", []string{"und", "en-US"}, []string{"en-US"}},
		{"blank entries dropped", []string{"", "  ", "fr-CA"}, []string{"fr-CA"}},
		{"nil", nil, []string{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := normaliseLanguageTags(tt.raw)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}
