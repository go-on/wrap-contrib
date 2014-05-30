package wraps

import (
	"net/http"
	"regexp"
)

type Matcher interface {
	Match(*http.Request) bool
}

type MatchFunc func(*http.Request) bool

func (mf MatchFunc) Match(req *http.Request) bool {
	return mf(req)
}

func And(ms ...Matcher) Matcher {
	return MatchFunc(
		func(req *http.Request) bool {
			for _, m := range ms {
				if !m.Match(req) {
					return false
				}
			}
			return true
		},
	)
}

func Or(ms ...Matcher) Matcher {
	return MatchFunc(
		func(req *http.Request) (doesMatch bool) {
			for _, m := range ms {
				if m.Match(req) {
					doesMatch = true
					break
				}
			}
			return
		},
	)
}

type MatchMethod string

func (mh MatchMethod) Match(r *http.Request) bool {
	return r.Method == string(mh)
}

type MatchPath string

func (mh MatchPath) Match(r *http.Request) bool {
	return r.URL.Path == string(mh)
}

type matchPathRegex regexp.Regexp

func MatchPathRegex(r *regexp.Regexp) Matcher {
	m := matchPathRegex(*r)
	return &m
}

func (mhr *matchPathRegex) Match(r *http.Request) bool {
	rg := regexp.Regexp(*mhr)
	return (&rg).MatchString(r.URL.Path)
}

type MatchHost string

func (mh MatchHost) Match(r *http.Request) bool {
	return r.URL.Host == string(mh)
}

func MatchHostRegex(r *regexp.Regexp) Matcher {
	m := matchHostRegex(*r)
	return &m
}

type matchHostRegex regexp.Regexp

func (mhr *matchHostRegex) Match(r *http.Request) bool {
	rg := regexp.Regexp(*mhr)
	return (&rg).MatchString(r.URL.Host)
}

type MatchScheme string

func (m MatchScheme) Match(r *http.Request) bool {
	return r.URL.Scheme == string(m)
}

type matchQuery struct {
	Key, Val string
}

func MatchQuery(key, val string) Matcher {
	return &matchQuery{key, val}
}

func (m *matchQuery) Match(r *http.Request) bool {
	return r.URL.Query().Get(m.Key) == m.Val
}

type MatchFragment string

func (m MatchFragment) Match(r *http.Request) bool {
	return r.URL.Fragment == string(m)
}

type matchHeader struct {
	Key, Val string
}

func MatchHeader(key, value string) Matcher {
	return &matchHeader{key, value}
}

func (m *matchHeader) Match(r *http.Request) bool {
	return r.Header.Get(m.Key) == m.Val
}
