from flask import Flask, request
import urllib2
import os
import re

app = Flask(__name__)

pod_ip = os.getenv("POD_IP")

def l5d_headers():
    headers = {k:v for k,v in request.headers.iteritems()
               if re.match('^(l5d-ctx|dtab-local)', k, re.I)}
    return headers

@app.route("/")
def hello():
    request = urllib2.Request("http://world/", headers=l5d_headers())
    return "Hello (" + pod_ip + ") " + urllib2.urlopen(request).read()

if __name__ == "__main__":
    app.run(host='0.0.0.0', port=7777)
