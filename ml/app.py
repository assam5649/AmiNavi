import base64
from flask import Flask, request, jsonify
from model import process_image

app = Flask(__name__)

@app.route("/convert", methods=["POST"])
def convert():
    image_bytes = request.data
    response_data = process_image(image_bytes)
    return jsonify(response_data)

if __name__ == "__main__":
    app.run(host="0.0.0.0", port=8501)
