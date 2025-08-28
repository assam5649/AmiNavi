import base64

def process_image(image_bytes: bytes) -> dict:
    print("Received image bytes (base64):", base64.b64encode(image_bytes).decode())

    return {"image": "aGVsbG93b3JsZA=="}