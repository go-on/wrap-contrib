wrap-contrib
============

Collection of middleware (wrappers) for [github.com/go-on/wrap](http://github.com/go-on/wrap)

[![Build Status](https://secure.travis-ci.org/go-on/wrap-contrib.png)](http://travis-ci.org/go-on/wrap-contrib) [![GoDoc](https://godoc.org/github.com/go-on/wrap-contrib?status.png)](https://godoc.org/github.com/go-on/wrap-contrib) [![Coverage Status](https://img.shields.io/coveralls/go-on/wrap-contrib.svg)](https://coveralls.io/r/go-on/wrap-contrib?branch=master) [![Project status](http://img.shields.io/status/stable.png?color=green)](#) [![Total views](https://sourcegraph.com/api/repos/github.com/go-on/wrap-contrib/counters/views.png)](https://sourcegraph.com/github.com/go-on/wrap-contrib)

Contributions
-------------

Yes, please! Make a pull request and let me see. If it is not matured consider adding it to [go-on/wrap-contrib-testing](http://github.com/go-on/wrap-contrib-testing).

Content
-------

- **body writer**: EscapeHTML, GZip, Reader
- **error handling**: Catch, Defer
- **caching**: ETag, IfNoneMatch, IfMatch
- **security**: Guard, CSRF (via nosurf)
- **combinators**: After, Before, Around, Fallback, First
- **REST**: Head, MethodOverride, GETHandler, POSTHandler, PUTHandler, DELETEHandler, PATCHHandler, OPTIONSHandler, HEADHandler
- **dispatching**: Dispatch, Map, MethodHandler, And, Or, MatchMethod, MatchPath, MatchHost, MatchScheme, MatchQuery, MatchHeader
- **http.Handler**: String, TextString, JSONString, CSSString, HTMLString, JavaScriptString
- **header manipulation**: ContentType, RemoveRequestHeader, RemoveResponseHeader, SetRequestHeader, SetResponseHeader


More Middleware
---------------

More (WIP,API may change) middleware can be found at [github.com/go-on/wrap-contrib-testing](https://github.com/go-on/wrap-contrib-testing)


Router
------

A router (WIP,API may change) can be found at [github.com/go-on/router](https://github.com/go-on/router)
