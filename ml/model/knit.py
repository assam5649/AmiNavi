import base64
import cv2
import numpy as np
import os
import io
import pandas as pd
import json
import logging
from tensorflow import keras

try:
    MODEL_PATH = "/app/model/knit3.keras"


    load_model = keras.models.load_model(MODEL_PATH)
    with open('/app/model/class_indices.json', 'r') as f:
        class_indices = json.load(f)
    class_label = class_indices

except Exception as e:
    import traceback
    import sys
    traceback.print_exc(file=sys.stderr)
    sys.exit(1)

def process_image(image_str: str) -> bytes:
    logging.info(f"Received image data size: {len(image_str)} bytes")
    logging.info(f"First 50 bytes of data: {image_str[:50]}")
    
    global load_model, class_label

    np_array = np.frombuffer(image_str, np.uint8)
    image = cv2.imdecode(np_array, cv2.IMREAD_COLOR)

    if image is None:
        raise ValueError("Error: 画像バイナリデータをOpenCVでデコードできませんでした。破損またはサポートされていない画像形式です。")

    gray = cv2.cvtColor(image, cv2.COLOR_BGR2GRAY)
    _, thresh = cv2.threshold(gray, 200, 255, cv2.THRESH_BINARY_INV)

    contours, _ = cv2.findContours(thresh, cv2.RETR_LIST, cv2.CHAIN_APPROX_SIMPLE)

    cell_contours = []
    min_area = 10
    max_area = 10000

    for cnt in contours:
        area = cv2.contourArea(cnt)
        if min_area < area < max_area:
            peri = cv2.arcLength(cnt, True)
            approx = cv2.approxPolyDP(cnt, 0.04 * peri, True)

            if len(approx) == 4:
                cell_contours.append(cnt)

    num_cells = len(cell_contours)
    print(f"検出されたセルの総数: {num_cells}")

    cell_info = []
    for cnt in cell_contours:
        M = cv2.moments(cnt)
        if M["m00"] != 0:
            cx = int(M["m10"] / M["m00"])
            cy = int(M["m01"] / M["m00"])
            x, y, w, h = cv2.boundingRect(cnt)
            cell_info.append({'center': (cx, cy), 'bbox': (x, y, w, h)})

    def coords_counter(coords, threshold):
        if not coords:
            return []

        sorted_coords = sorted(list(set(coords)))
        clusters = []
        current_cluster = [sorted_coords[0]]

        for i in range(1, len(sorted_coords)):
            if sorted_coords[i] - current_cluster[-1] < threshold:
                current_cluster.append(sorted_coords[i])
            else:
                clusters.append(int(np.mean(current_cluster)))
                current_cluster = [sorted_coords[i]]
        clusters.append(int(np.mean(current_cluster)))
        return clusters

    x_coords = [info['center'][0] for info in cell_info]
    y_coords = [info['center'][1] for info in cell_info]

    x_clusters = coords_counter(x_coords, threshold=10)
    y_clusters = coords_counter(y_coords, threshold=10)

    num_cols = len(x_clusters)
    num_rows = len(y_clusters)

    print(f"検出されたグリッドの横の数: {num_cols}")
    print(f"検出されたグリッドの縦の数: {num_rows}")

    amizu = [[None for _ in range(num_cols)] for _ in range(num_rows)]

    for info in cell_info:
        cx, cy = info['center']
        x_idx = np.argmin([abs(cx - x_c) for x_c in x_clusters])
        y_idx = np.argmin([abs(cy - y_c) for y_c in y_clusters])
        amizu[y_idx][x_idx] = info['bbox']

    print("\nグリッドの配列:")

    for row in amizu:
        print([ '-' if cell is None else cell for cell in row ])

    output_image = image.copy()
    cv2.drawContours(output_image, cell_contours, -1, (0, 255, 0), 2)

    if x_clusters:
        for x in x_clusters:
            cv2.line(output_image, (x, 0), (x, output_image.shape[0]), (255, 0, 0), 2)

    if y_clusters:
        for y in y_clusters:
            cv2.line(output_image, (0, y), (output_image.shape[1], y), (0, 0, 255), 2)

    amizu_predict = [[None for _ in range(num_cols)] for _ in range(num_rows)]

    for row_idx, row in enumerate(amizu):
        for col_idx, bbox in enumerate(row):
            if bbox is not None:
                x, y, w, h = bbox
                cell_image = image[y:y+h, x:x+w]
                
                img = cv2.cvtColor(cell_image, cv2.COLOR_BGR2RGB)
                img = cv2.resize(img, (224, 224))
                img_array = img / 255.0

                img_array = img_array.astype(np.float32) 

                img_array = np.expand_dims(img_array, axis=0)

                print(f"DEBUG: Input Shape = {img_array.shape}, Dtype = {img_array.dtype}")

                predictions = load_model.predict(img_array)
                probabilities = predictions[0].flatten()
                predicted_class_index = np.argmax(probabilities)

                max_probability = probabilities[predicted_class_index]

                if predicted_class_index < len(class_label):
                    predicted_class_name = class_label[str(predicted_class_index)]
                else:
                    predicted_class_name = '-'
                    print(f"Warning: Cell ({row_idx+1},{col_idx+1}) predicted index {predicted_class_index} out of bounds (Max Index: {len(class_label)-1}).")
                
                amizu_predict[row_idx][col_idx] = predicted_class_name
                
                print(f"cell ({row_idx+1}, {col_idx+1}): {predicted_class_name} (確率: {max_probability:.2f})")
                
    print("\n結果:")
    for row in amizu_predict:
        print(row)

    df_amizu = pd.DataFrame(amizu_predict)
    df_amizu.fillna('-', inplace=True)

    header_info = [f"#rows,{num_rows}", f"#cols,{num_cols}"]
    header_df = pd.DataFrame(header_info)

    final_df = pd.concat([header_df, df_amizu], ignore_index=True)

    csv_buffer = io.StringIO()
    final_df.to_csv(csv_buffer, index=False, header=False, encoding='utf-8')

    csv_string = csv_buffer.getvalue()

    go_byte_array = csv_string.encode('utf-8')

    print("Return csv bytes (base64):", base64.b64encode(go_byte_array).decode())

    return go_byte_array

