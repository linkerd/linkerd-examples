from flask import Flask
from header_utils import l5d_headers
import urllib2
import os


app = Flask(__name__)

pod_ip = os.getenv("POD_IP")

@app.route("/")
def hello():
    request = urllib2.Request("http://world/", headers=l5d_headers())
    return "Hello (" + pod_ip + ") " + urllib2.urlopen(request).read()

if __name__ == "__main__":
    app.run(host='0.0.0.0', port=7777, threaded=True)
