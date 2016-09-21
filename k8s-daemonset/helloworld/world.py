from flask import Flask
import os

app = Flask(__name__)

pod_ip = os.getenv("POD_IP")

@app.route("/")
def world():
    return "world (" + pod_ip + ")!"

if __name__ == "__main__":
    app.run(host='0.0.0.0', port=7778)
