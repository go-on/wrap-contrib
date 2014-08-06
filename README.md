wrap-contrib
============

Collection of middleware (wrappers) for [github.com/go-on/wrap](http://github.com/go-on/wrap)

[![Build Status](https://drone.io/github.com/go-on/wrap-contrib/status.png)](https://drone.io/github.com/go-on/wrap-contrib/latest) [![GoDoc](https://godoc.org/github.com/go-on/wrap-contrib/wraps?status.png)](https://godoc.org/github.com/go-on/wrap-contrib/wraps) [![Coverage Status](https://img.shields.io/coveralls/go-on/wrap-contrib.svg)](https://coveralls.io/r/go-on/wrap-contrib?branch=master) [![Total views](https://sourcegraph.com/api/repos/github.com/go-on/wrap-contrib/counters/views.png)](https://sourcegraph.com/github.com/go-on/wrap-contrib)

Contributions
-------------

Yes, please! Make a pull request and let me see. If it is not matured consider adding it to [go-on/wrap-contrib-testing](http://github.com/go-on/wrap-contrib-testing).

Content
-------

- **body writer**: EscapeHTML, GZip, ReadSeeker
- **error handling**: Error, Catch, Defer
- **caching**: ETag, IfNoneMatch, IfMatch
- **combinators**: After, Before, Around, Fallback, First, Guard
- **REST**: Head, MethodOverride, (GET|POST|PUT|DELETE|PATCH|OPTIONS|HEAD)Handler
- **dispatching**: Dispatch, Map, MethodHandler, And, Or, Match(Method|Path|Host|Scheme|Query|Header)
- **http.Handler**: (Text|JSON|CSS|HTML|JavaScript)String
- **header manipulation**: ContentType, (Set|Remove)(Request|Response)Header
- **integration of 3rd party middleware**: wrapnosurf (github.com/justinas/nosurf), wraphttpauth (github.com/abbot/go-http-auth), wrapsession (github.com/gorilla/sessions)


More Middleware
---------------

More (WIP,API may change) middleware can be found at [github.com/go-on/wrap-contrib-testing](https://github.com/go-on/wrap-contrib-testing)


Router
------

A router (WIP,API may change) can be found at [github.com/go-on/router](https://github.com/go-on/router)
