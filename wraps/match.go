package wraps

import (
	"net/http"
	"regexp"
)

// Matcher is an interface for types that can match against a http.Request
type Matcher interface {
	// Match returns if the request matches
	Match(*http.Request) bool
}

// MatchFunc is a Matcher based on a function
type MatchFunc func(*http.Request) bool

// Match implements Matcher for the MatchFunc
func (mf MatchFunc) Match(req *http.Request) bool {
	return mf(req)
}

// And logically combines different matchers to a single matcher
// that only matches if all matchers match.
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

// And logically combines different matchers to a single matcher
// that matches if one of the matchers match.
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

// MatchMethod matches based on the request method (like GET, POST etc.)
type MatchMethod string

func (mh MatchMethod) Match(r *http.Request) bool {
	return r.Method == string(mh)
}

// MatchMethod matches based on the request path
type MatchPath string

func (mh MatchPath) Match(r *http.Request) bool {
	return r.URL.Path == string(mh)
}

type matchPathRegex regexp.Regexp

// MatchPathRegex matches the request path against a reqular expression
func MatchPathRegex(r *regexp.Regexp) Matcher {
	m := matchPathRegex(*r)
	return &m
}

func (mhr *matchPathRegex) Match(r *http.Request) bool {
	rg := regexp.Regexp(*mhr)
	return (&rg).MatchString(r.URL.Path)
}

// MatchHost matches based on the request host
type MatchHost string

func (mh MatchHost) Match(r *http.Request) bool {
	return r.URL.Host == string(mh)
}

// MatchHostRegex matches the request host against a reqular expression
func MatchHostRegex(r *regexp.Regexp) Matcher {
	m := matchHostRegex(*r)
	return &m
}

type matchHostRegex regexp.Regexp

func (mhr *matchHostRegex) Match(r *http.Request) bool {
	rg := regexp.Regexp(*mhr)
	return (&rg).MatchString(r.URL.Host)
}

// MatchScheme matches based on the request scheme
type MatchScheme string

func (m MatchScheme) Match(r *http.Request) bool {
	return r.URL.Scheme == string(m)
}

type matchQuery struct {
	Key, Val string
}

// MatchQuery matches based on the request query
func MatchQuery(key, val string) Matcher {
	return &matchQuery{key, val}
}

func (m *matchQuery) Match(r *http.Request) bool {
	return r.URL.Query().Get(m.Key) == m.Val
}

type matchHeader struct {
	Key, Val string
}

// MatchHeader matches based on the request header
func MatchHeader(key, value string) Matcher {
	return &matchHeader{key, value}
}

func (m *matchHeader) Match(r *http.Request) bool {
	return r.Header.Get(m.Key) == m.Val
}
