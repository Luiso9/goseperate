import cv2
import numpy as np
import json
from flask import Flask, request, send_file
from sklearn.cluster import MiniBatchKMeans
from io import BytesIO

app = Flask(__name__)

def generate_preview(image_path, num_colors, d, sigmaColor, sigmaSpace):
    image = cv2.imread(image_path, cv2.IMREAD_UNCHANGED)
    if image is None:
        return None, "Could not read image"

    height, width = image.shape[:2]
    total_pixels = height * width
    batch_size = max(512, min(total_pixels // (num_colors * 10), 8192))
    random_state = (width + height) % 100

    if image.shape[-1] == 4:
        alpha_channel = image[:, :, 3]
        bgr_image = image[:, :, :3]
    else:
        bgr_image = image
        alpha_channel = np.ones((height, width), dtype=np.uint8) * 255
    
    lab_image = cv2.cvtColor(bgr_image, cv2.COLOR_BGR2LAB)
    reshaped_image = lab_image.reshape((-1, 3)).astype(np.float32)
    
    kmeans = MiniBatchKMeans(n_clusters=num_colors, random_state=random_state, batch_size=batch_size)
    labels = kmeans.fit_predict(reshaped_image)
    centers = np.uint8(kmeans.cluster_centers_)
    segmented_lab = centers[labels].reshape(lab_image.shape)
    
    segmented_bgr = cv2.cvtColor(segmented_lab, cv2.COLOR_LAB2BGR)
    result_image = cv2.bilateralFilter(segmented_bgr, d, sigmaColor, sigmaSpace)
    
    success, buffer = cv2.imencode('.png', result_image)
    if not success:
        return None, "Failed to encode image"

    return buffer.tobytes(), None

@app.route("/preview", methods=["POST"])
def preview_api():
    try:
        data = request.json
        image_path = data.get("image_path")
        num_colors = data.get("num_colors", 12)
        d = data.get("d", 1)
        sigmaColor = data.get("sigmaColor", 10)
        sigmaSpace = data.get("sigmaSpace", 10)

        preview_buffer, error = generate_preview(image_path, num_colors, d, sigmaColor, sigmaSpace)
        if error:
            return json.dumps({"error": error}), 400

        return send_file(BytesIO(preview_buffer), mimetype="image/png")
    except Exception as e:
        return json.dumps({"error": str(e)}), 500

if __name__ == "__main__":
    app.run(host="0.0.0.0", port=5037)
