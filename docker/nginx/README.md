## Docker image that compiles NGINX from source

The NGINX Headers More module doesn't currently accept wildcard matches
for clearing input headers (it's on a [branch](https://github.com/ghedo/headers-more-nginx-module.git)).
So if we want to use this
capability we need to compile nginx from source and add this module.

We're building a custom docker image from a branch of the
[Headers More](https://github.com/openresty/headers-more-nginx-module) module.
This dockerfile installs the tools necessary
to compile nginx and the headers more module from source and does the compilation.

This image is pre-built and available on docker hub as
[buoyantio/nginx:1.10.2](https://hub.docker.com/r/buoyantio/nginx/tags/),
so there should be no need to rebuild it.

If you would still like to build this image, you'll first need to locally clone
and check out a custom branch
([wildcard_in](https://github.com/ghedo/headers-more-nginx-module/tree/wildcard_in))
of the headers-more-nginx-module repo, as follows:

```
$ git clone https://github.com/ghedo/headers-more-nginx-module.git
$ cd headers-more-nginx-module
$ git checkout wildcard_in
$ cd ..
```

Then, copy that folder into the same place as this Dockerfile, and then build:
```
$ cp headers-more-nginx-module <location-of-this-dockerfile>
$ docker build -t buoyantio/nginx:<tag-name> .
```

The dockerfile was based on instructions from
[NGINX's compilation instructions](https://www.nginx.com/resources/admin-guide/installing-nginx-open-source/)
 as well as the instructions [here](https://github.com/arut/nginx-rtmp-module/wiki/Getting-started-with-nginx-rtmp).