import base64
from flask import Flask, request, Response
from model.knit import process_image
import json
import traceback
import sys
import logging

app = Flask(__name__)

@app.route("/health", methods=["GET"])
def health_check():
    return Response(response=json.dumps({"status": "ok"}), status=200, mimetype='application/json')

@app.route("/convert", methods=["POST"])
def convert():
    try:
        if 'file' not in request.files:
            logging.error("No file part in the request (Expected field name: 'file').")
            return "No file part in the request.", 400
        
        image_file = request.files['file']
        image_bytes = image_file.read()

        print(f"DEBUG: Received image data size: {len(image_bytes)} bytes") 

        if not image_bytes:
            logging.error("Empty image data received.")
            return "Error: Empty image data received.", 400

        response_data = process_image(image_bytes)

        return Response(
            response_data,
            mimetype='text/csv',
            status=200
        )
    
    except ValueError as e:
        error_message = str(e) 
        traceback.print_exc(file=sys.stderr) 
        
        return Response(
            json.dumps({"error": "Data validation failed.", "details": error_message}),
            mimetype='application/json',
            status=400
        )

    except Exception as e:
        error_message = str(e) 
        traceback.print_exc(file=sys.stderr) 
        
        return Response(
            json.dumps({"error": "Internal server error occurred.", "details": error_message}),
            mimetype='application/json',
            status=500
        )
        
if __name__ == "__main__":
    app.run(host="0.0.0.0", port=8501)
