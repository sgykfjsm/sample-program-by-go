package regexp_vs_replacer

import (
	"regexp"
	"strings"
)

func DoRegexp(re *regexp.Regexp, contents string) string {
	return re.ReplaceAllStringFunc(contents, func(s string) string {
		return strings.ToUpper(s)
	})
}

func DoReplacer(r *strings.Replacer, contents string) string {
	return r.Replace(contents)
}
