## Docker image that compiles NGINX from source

The nginx Headers More module doesn't currently accept wildcard matches
for clearing input headers (it's on a branch somehwere). So if we want to use this
capability we need to compile nginx from source and add this module.

We're building a custom docker image from a branch of the [Headers More](https://github.com/openresty/headers-more-nginx-module) module. This dockerfile installs the tools necessary
to compile nginx and the headers more module from source, then copies in our
custom nginx config.

This image is pre-built and available on docker hub as[buoyantio/nginx:1.10.2](https://hub.docker.com/r/buoyantio/nginx/tags/),
so there should be no need to rebuild it.
If you would still like to build this, get a copy of the
[wildcard_in](https://github.com/ghedo/headers-more-nginx-module/tree/wildcard_in) branch, and adjust this line in the dockerfile to get the module files
from wherever you've saved the branch.

```
ADD headers-more-nginx-module /headers-more-nginx-module
```

The dockerfile was based on instructions from:
[NGINX's compilation instructions](https://www.nginx.com/resources/admin-guide/installing-nginx-open-source/)
 as well as the instructions [here](https://github.com/arut/nginx-rtmp-module/wiki/Getting-started-with-nginx-rtmp).