package application

import (
	"strings"

	"golang.org/x/text/language"
)

// PreferredLanguages returns the languages the user has chosen for
// application user interfaces, most preferred first, as canonical BCP-47 tags
// such as "nb-NO", "en-US" or "zh-Hans-CN". It returns an empty slice when the
// platform reports none; the app then applies its own default.
//
// The list comes from the operating system's language settings. It is not the
// regional format for dates, numbers and currency, which every platform keeps
// as a separate setting: someone running English Windows in Norway gets an
// English tag first, not "nb-NO".
//
// Each call asks the platform; nothing is cached. It needs no application
// instance, so it can be called before [New], except on Android, where the
// list comes through the WailsBridge and is empty until the activity exists.
func PreferredLanguages() []string {
	return normaliseLanguageTags(platformPreferredLanguages())
}

// PreferredLanguage returns the first entry of [PreferredLanguages], or ""
// when the list is empty.
func PreferredLanguage() string {
	if langs := PreferredLanguages(); len(langs) > 0 {
		return langs[0]
	}
	return ""
}

// posixPreferredLanguages applies GNU gettext's rules to the POSIX locale
// environment read through getenv:
//
//   - the message locale is LC_ALL, else LC_MESSAGES, else LANG;
//   - if that is unset, "C" or "POSIX" (with any codeset, so "C.UTF-8" too),
//     the user asked for no localisation, and the list is empty: gettext
//     ignores LANGUAGE in that case as well;
//   - otherwise the entries of LANGUAGE, a colon-separated priority list such
//     as "nb:no:en", come first, followed by the message locale.
//
// The entries are returned raw; normaliseLanguageTags turns them into tags.
// getenv is a parameter so tests need not touch the process environment.
func posixPreferredLanguages(getenv func(string) string) []string {
	locale := getenv("LC_ALL")
	if locale == "" {
		locale = getenv("LC_MESSAGES")
	}
	if locale == "" {
		locale = getenv("LANG")
	}
	if isPOSIXCLocale(locale) {
		return nil
	}

	var out []string
	for _, entry := range strings.Split(getenv("LANGUAGE"), ":") {
		if entry != "" {
			out = append(out, entry)
		}
	}
	return append(out, locale)
}

// isPOSIXCLocale reports whether a POSIX locale name asks for no
// localisation: unset, "C" or "POSIX", with or without a codeset or modifier.
func isPOSIXCLocale(locale string) bool {
	name, _, _ := strings.Cut(locale, "@")
	name, _, _ = strings.Cut(name, ".")
	return name == "" || name == "C" || name == "POSIX"
}

// posixScriptModifiers maps the glibc locale modifiers that name a script to
// its ISO 15924 code. Other modifiers, such as "@euro", say nothing about the
// language and are dropped.
var posixScriptModifiers = map[string]string{
	"latin":      "Latn",
	"cyrillic":   "Cyrl",
	"devanagari": "Deva",
}

// normaliseLanguageTags turns the raw strings a platform reports, in POSIX
// ("nb_NO.UTF-8", "sr_RS@latin") or BCP-47 ("zh-hans-cn") form, into
// canonical BCP-47 tags. Entries that don't parse are dropped, and a
// duplicate keeps its first position. The result is never nil.
func normaliseLanguageTags(raw []string) []string {
	out := make([]string, 0, len(raw))
	seen := make(map[string]bool, len(raw))
	for _, s := range raw {
		tag, ok := languageTag(s)
		if !ok || seen[tag] {
			continue
		}
		seen[tag] = true
		out = append(out, tag)
	}
	return out
}

// languageTag converts one platform string to a canonical BCP-47 tag.
func languageTag(raw string) (string, bool) {
	// POSIX names are language[_territory][.codeset][@modifier].
	s, modifier, _ := strings.Cut(strings.TrimSpace(raw), "@")
	s, _, _ = strings.Cut(s, ".")
	if isPOSIXCLocale(s) {
		return "", false
	}

	tag, err := language.Parse(strings.ReplaceAll(s, "_", "-"))
	if err != nil || tag == language.Und {
		return "", false
	}
	if code, ok := posixScriptModifiers[strings.ToLower(modifier)]; ok {
		script := language.MustParseScript(code)
		if withScript, err := language.Compose(tag, script); err == nil {
			tag = withScript
		}
	}
	return tag.String(), true
}
