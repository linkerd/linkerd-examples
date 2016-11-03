from flask import Flask, request, jsonify
import urllib2
import json
import os
import re
import sys


app = Flask(__name__)
pod_ip = os.getenv("POD_IP")

def l5d_headers():
    headers = {k:v for k,v in request.headers.iteritems()
               if re.match('^(l5d-ctx|dtab-local)', k, re.I)}
    return headers

@app.route("/")
def api():
  request = urllib2.Request("http://hello/", headers=l5d_headers())
  return jsonify([
    {'api_result': "api (" + pod_ip + ") calls " + urllib2.urlopen(request).read()}
  ])

if __name__ == "__main__":
    app.run(host='0.0.0.0', port=7779)
