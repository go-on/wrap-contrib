// Copyright (c) 2014 Marc René Arns. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

/*
Package wrapcontrib provides middleware (wraps) and helpers for http://github.com/go-on/wrap.

A walk through some of the middlewares can be found at http://metakeule.github.io/article/tour-de-wrap-middleware.html.

Content

  - body writer: EscapeHTML, GZip, Reader
  - error handling: Catch, Defer
  - caching: ETag, IfNoneMatch, IfMatch
  - combinators: After, Before, Around, Fallback, First, Guard
  - REST: Head, MethodOverride, GETHandler, POSTHandler, PUTHandler, DELETEHandler, PATCHHandler, OPTIONSHandler, HEADHandler
  - dispatching: Dispatch, Map, MethodHandler, And, Or, MatchMethod, MatchPath, MatchHost, MatchScheme, MatchQuery, MatchHeader
  - http.Handler: String, TextString, JSONString, CSSString, HTMLString, JavaScriptString
  - header manipulation: ContentType, RemoveRequestHeader, RemoveResponseHeader, SetRequestHeader, SetResponseHeader
  - integration of 3rd party middleware: wrapnosurf (github.com/justinas/nosurf), wraphttpauth (github.com/abbot/go-http-auth)

Contributions

Yes, please! Make a pull request and let me see. If it is not matured consider adding it to http://github.com/go-on/wrap-contrib-testing.


More Middleware

More (WIP,API may change) middleware can be found at http://github.com/go-on/wrap-contrib-testing

Example

	package main

	import (
		"fmt"
		"github.com/go-on/wrap"
		"github.com/go-on/wrap-contrib/wraps"
		"net/http"
	)

	type catcher struct{}

	func (c catcher) Catch(p interface{}, w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "catched: %s", p)
	}

	type panicker struct{}

	func (p panicker) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
		panic("panic mysterious person found")
	}

	func main() {

		handler := wrap.New(
			wraps.HTMLContentType,
			wraps.GZip,
			wraps.ETag,
			wraps.Before(wraps.String(
			 `<!DOCTYPE html><html lang="en"><head></head><body>`)),
			wraps.After(wraps.String(`</body></html>`)),
			wraps.Catch(catcher{}),
			wraps.Map(
				wraps.MatchQuery("name", "peter"), wraps.String("Hi Peter!"),
				wraps.MatchQuery("name", "mary"), wraps.String("Hello Mary!"),
				wraps.MatchQuery("name", "mister-x"), panicker{},
			),
			wraps.String(`
				<a href="/?name=peter">Peter</a><br />
				<a href="/?name=mary">Mary</a><br>
				<a href="/?name=mister-x">Mister X</a><br>
				`),
		)

		fmt.Println("go to localhost:8080")

		err := http.ListenAndServe(":8080", handler)

		if err != nil {
			println(err)
		}
	}


*/
package wrapcontrib
