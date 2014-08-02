package wraps

import (
	"net/http"
	"regexp"
	"testing"
)

type matchQ struct {
	Matcher
	*http.Request
}

func TestMatch(t *testing.T) {
	reqA, _ := http.NewRequest("GET", "/a", nil)
	reqB, _ := http.NewRequest("POST", "http://localhost/b?name=me#frag", nil)
	reqB.Header.Set("X-my", "val")
	//reqB.URL.Host = "localhost"
	//reqB.URL.Scheme = "http"

	tests := map[matchQ]bool{
		matchQ{MatchPath("/b"), reqA}: false,
		matchQ{MatchPath("/b"), reqB}: true,

		matchQ{MatchQuery("name", "me"), reqA}: false,
		matchQ{MatchQuery("name", "me"), reqB}: true,

		matchQ{MatchHeader("X-my", "val"), reqA}: false,
		matchQ{MatchHeader("X-my", "val"), reqB}: true,

		matchQ{MatchHost("localhost"), reqA}: false,
		matchQ{MatchHost("localhost"), reqB}: true,

		matchQ{MatchMethod("POST"), reqA}: false,
		matchQ{MatchMethod("POST"), reqB}: true,

		matchQ{MatchScheme("http"), reqA}: false,
		matchQ{MatchScheme("http"), reqB}: true,

		matchQ{MatchHostRegex(regexp.MustCompile("local")), reqA}: false,
		matchQ{MatchHostRegex(regexp.MustCompile("local")), reqB}: true,

		matchQ{MatchPathRegex(regexp.MustCompile("b")), reqA}: false,
		matchQ{MatchPathRegex(regexp.MustCompile("b")), reqB}: true,
	}

	for mQ, doesMatch := range tests {
		res := mQ.Matcher.Match(mQ.Request)
		if res != doesMatch {
			t.Errorf("error in matching %s %s with matcher %T, expected %v but got %v",
				mQ.Request.Method, mQ.Request.URL.Path, mQ.Matcher, doesMatch, res,
			)
		}
	}
}

func TestMatchAnd(t *testing.T) {
	req, _ := http.NewRequest("GET", "/a", nil)

	m := And(MatchMethod("GET"), MatchPath("/a"))

	if !m.Match(req) {
		t.Error("should match but does not")
	}
}

func TestMatchNotAnd(t *testing.T) {
	req, _ := http.NewRequest("GET", "/a", nil)

	m := And(MatchMethod("GET"), MatchPath("/b"))

	if m.Match(req) {
		t.Error("should not match but does")
	}
}

func TestMatchOr(t *testing.T) {
	req, _ := http.NewRequest("GET", "/a", nil)

	m := Or(MatchMethod("GET"), MatchPath("/b"))

	if !m.Match(req) {
		t.Error("should match but does not")
	}
}

func TestMatchNotOr(t *testing.T) {
	req, _ := http.NewRequest("GET", "/a", nil)

	m := Or(MatchMethod("POST"), MatchPath("/b"))

	if m.Match(req) {
		t.Error("should not match but does")
	}
}
