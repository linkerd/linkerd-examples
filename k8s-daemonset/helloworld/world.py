from flask import Flask, request, abort
import os
import random

app = Flask(__name__)

pod_ip = os.getenv("POD_IP")
target_world = os.getenv("TARGET_WORLD", "world")

failure_rate = 0.0

@app.route("/")
def world():
    r = random.random()
    print r, failure_rate
    if (r < failure_rate):
        abort(500)
    else:
        return target_world + " (" + pod_ip + ")!"

@app.route("/failure", methods=["PUT"])
def set_failure_rate():
    global failure_rate
    failure_rate = float(request.args.get('rate', '0.0'))
    return str(failure_rate)

if __name__ == "__main__":
    app.run(host='0.0.0.0', port=7778, threaded=True)
