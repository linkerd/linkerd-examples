from flask import Flask, jsonify
from header_utils import l5d_headers
import urllib2
import os


app = Flask(__name__)
pod_ip = os.getenv("POD_IP")

@app.route("/")
def api():
  request = urllib2.Request("http://hello/", headers=l5d_headers())
  return jsonify([
    {'api_result': "api (" + pod_ip + ") calls " + urllib2.urlopen(request).read()}
  ])

if __name__ == "__main__":
    app.run(host='0.0.0.0', port=7779)
