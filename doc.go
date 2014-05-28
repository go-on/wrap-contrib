// Copyright (c) 2014 Marc René Arns. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

/*
Package wrapcontrib provides middleware (wraps) and helpers for github.com/go-on/wrap.

Status

100% test coverage.
This package is considered complete and the API is stable  (with the exception of the helper subpackage).

More Middleware

More (WIP,API may change) middleware can be found at github.com/go-on/wrap-contrib-testing

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
