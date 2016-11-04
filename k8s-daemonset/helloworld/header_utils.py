from flask import request
import re

# Context headers (l5d-ctx-*) are generated and read by linkerd instances.
# Applications should forward all context headers in order for all linkerd
# features to work.
def l5d_headers():
    headers = {k:v for k,v in request.headers.iteritems()
               if re.match('^(l5d-ctx|dtab-local|l5d-dtab)', k, re.I)}
    return headers