import cv2
import numpy as np
import sys
import json
import os
from sklearn.cluster import MiniBatchKMeans

def rgb_to_hex(color):
    return "{:02x}{:02x}{:02x}".format(color[0], color[1], color[2])

def extract_colors(image_path, num_colors, output_dir):
    image = cv2.imread(image_path, cv2.IMREAD_UNCHANGED)
    if image is None:
        print(json.dumps({"error": "Could not read image"}))
        sys.exit(1)
    
    is_transparent = image.shape[-1] == 4 if len(image.shape) == 3 else False
    height, width = image.shape[:2]
    
    if is_transparent:
        alpha_channel = image[:, :, 3]
        rgb_image = cv2.cvtColor(image, cv2.COLOR_BGRA2RGB)
    else:
        rgb_image = cv2.cvtColor(image, cv2.COLOR_BGR2RGB)
        alpha_channel = np.ones((height, width), dtype=np.uint8) * 255
    
    reshaped_image = rgb_image.reshape((-1, 3)).astype(np.float32)
    
    random_state = (width + height) % 100
    batch_size = max(702, min(height * width // (num_colors * 10), 8192))
    
    kmeans = MiniBatchKMeans(n_clusters=num_colors, random_state=random_state, batch_size=batch_size)
    labels = kmeans.fit_predict(reshaped_image)
    centers = np.uint8(kmeans.cluster_centers_)
    
    labels_reshaped = labels.reshape(height, width)
    
    os.makedirs(output_dir, exist_ok=True)
    
    extracted_files = []
    
    for i, center in enumerate(centers):
        hex_color = rgb_to_hex(center)
        
        mask = (labels_reshaped == i)
        
        result = np.zeros((height, width, 4), dtype=np.uint8)
        result[:, :, :3] = center
        result[:, :, 3] = np.where(mask, alpha_channel, 0)
        
        output_path = os.path.join(output_dir, f"{hex_color}.png")
        cv2.imwrite(output_path, cv2.cvtColor(result, cv2.COLOR_RGBA2BGRA))
        extracted_files.append(hex_color)
    
    print(json.dumps({"extracted": extracted_files}))
    sys.exit(0)

if __name__ == "__main__":
    if len(sys.argv) != 4:
        print(json.dumps({"error": "Usage: python3 cluster.py <image_path> <num_colors> <output_dir>"}))
        sys.exit(1)
    
    image_path = sys.argv[1]
    num_colors = int(sys.argv[2])
    output_dir = sys.argv[3]
    
    extract_colors(image_path, num_colors, output_dir)