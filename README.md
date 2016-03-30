# Linkerd Examples #

**Goals**

- Illustrate how to use linkerd in kubernetes
- Have some fun

**Non-goals**

- Implement a useful application
- Provide exemplary Go code

## The Application ##

Welcome to Gob's Microservice!

![gob](https://media.giphy.com/media/qJxFuXXWpkdEI/giphy.gif)

_Gob's program_, from the television show _Arrested Development_, was
a childish, inane program.  So, naturally, we've turned it into a
microservice web application that can be run at scale!

This application consists of several components:

- _websvc_ -- Gob's frontend -- serves plaintext
- _wordsvc_ -- chooses a word for _web_ when one isn't provided
- _gensvc_ -- given a word and a limit, generates a stream of text

The web service is fairly simple (and entirely plaintext):

```
$ curl -s 'localhost:8080/'                  
Gob's web service!

Send me a request like:

  localhost:8080/gob

You can tell me what to say with:

  localhost:8080/gob?text=WHAT_TO_SAY&limit=NUMBER
```

_websvc_ may call both _wordsvc_ and _gensvc_ to satisfy a request.  

_wordsvc_ and _gensvc_ implement RPC-ish interfaces with HTTP and JSON.

All three services are implemented in Go with no shared code.  They
may be built and run independently.
