gorequests
==========

A wrapper for the golang stf lib http request/response functions.

This is just a very early version. It currently supports GET, POST, PUT, and DELETE requests. You can access/save files from a response directly and there is a basic function to tell gorequests to retry a request depending on the status code.

TODO
----
- add other verbs
- basic auth
- SSL control over time.
- proper tests
- proper documentation


If you want to help please fork and PR. Any help appreciated.

Old examples i had here are outdated. Check [tests/request_test.go](https://github.com/fvbock/gorequests/blob/master/tests/request_test.go) for basic usage examples.
