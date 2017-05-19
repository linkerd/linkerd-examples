# HelloWorld Identifier

This is an example of an identifier which mutates the request.  It works by
setting a header on the request and then returning an `UnidentifiedRequest`
which causes linkerd to fall back to the next identifier (if there is one).

Note that while this is currently the best/only way to implement things like
request mutation, rate limitting, authentication, etc. the identifier interface
is not very well suited to these use cases.  In the future, we hope to add
better plugin interfaces for adding these kinds of things.
