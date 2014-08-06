// Copyright (c) 2014 Marc René Arns. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

/*
Package wraps provides middleware for http://github.com/go-on/wrap.

A walk through some of the middlewares can be found at http://metakeule.github.io/article/tour-de-wrap-middleware.html.

Content

  body writer:
    - EscapeHTML
    - GZip
    - ReadSeeker

  error handling:
    - Error
    - Catch
    - Defer

  caching:
    - ETag
    - IfNoneMatch
    - IfMatch

  combinators:
    - After
    - Before
    - Around
    - Fallback
    - First
    - Guard

  dispatching:
    - Dispatch
    - Map
    - MethodHandler
    - And
    - Or
    - Match(Method|Path|Host|Scheme|Query|Header)

  REST:
    - Head
    - MethodOverride
    - (GET|POST|PUT|DELETE|PATCH|OPTIONS|HEAD)Handler

  http.Handler:
    - (Text|JSON|CSS|HTML|JavaScript)String

  header manipulation:
    - ContentType
    - (Remove|Set)(Request|Response)Header

Integration of 3rd party middleware	(http://godoc.org/github.com/go-on/wrap-contrib/third-party)

    - wrapnosurf (github.com/justinas/nosurf)
    - wraphttpauth (github.com/abbot/go-http-auth)
    - wrapsession (github.com/gorilla/sessions)

Contributions

Yes, please! Make a pull request and let me see. If it is not matured consider adding it to http://github.com/go-on/wrap-contrib-testing.


More Middleware

More (WIP,API may change) middleware can be found at http://github.com/go-on/wrap-contrib-testing

*/
package wraps
