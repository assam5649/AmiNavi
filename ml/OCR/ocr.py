import pytesseract
import logging
import numpy as np
from PIL import Image
import cv2

def ocr_image(image_bytes: bytes) -> str:
    logging.info(f"Received image data size: {len(image_bytes)} bytes")
    logging.info(f"First 50 bytes of data: {image_bytes[:50]}")

    np_array = np.frombuffer(image_bytes, np.uint8)
    image = cv2.imdecode(np_array, cv2.IMREAD_COLOR)

    if image is None:
        raise ValueError("Error: 画像バイナリデータをOpenCVでデコードできませんでした。破損またはサポートされていない形式です。")

    pil_img = Image.fromarray(cv2.cvtColor(image, cv2.COLOR_BGR2RGB))

    text = pytesseract.image_to_string(pil_img, lang="eng")
    print(text)

    return text
