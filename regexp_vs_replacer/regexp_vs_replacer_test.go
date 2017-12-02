package regexp_vs_replacer

import (
	"io/ioutil"
	"regexp"
	"strings"
	"testing"
)

var (
	r        *strings.Replacer
	re       *regexp.Regexp
	contents string
	keywords = []string{"One", "morning", "when", "Gregor", "Samsa", "woke", "from", "troubled", "dreams"}
	pairList = make([]string, len(keywords)*2)
)

func init() {
	_contents, err := ioutil.ReadFile("kafka.txt")
	if err != nil {
		panic(err)
	}
	contents = string(_contents)

	re = regexp.MustCompile(`(` + strings.Join(keywords, "|") + `)`)

	for i, keyword := range keywords {
		pairList[i*2], pairList[i*2+1] = keyword, strings.ToUpper(keyword)
	}
	r = strings.NewReplacer(pairList...)
}

func BenchmarkDoRegexp(b *testing.B) {
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		DoRegexp(re, contents)
	}
}

func BenchmarkDoReplacer(b *testing.B) {
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		DoReplacer(r, contents)
	}
}

func TestDoRegexp_Equal_To_DoReplacer(t *testing.T) {
	if DoRegexp(re, contents) != DoReplacer(r, contents) {
		t.Fatal("DoRegexp is not equal to DoReplacer")
	}
}
