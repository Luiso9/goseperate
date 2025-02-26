# GoSeperate - Color Separation API

## Overview

GoSeperate is a Go-based REST API that extracts and separates colors from an image into different layers. The processed layers are returned as a ZIP file.

## Features

-  Upload an image and specify the number of colors to extract
-  Extracted layers are saved with filenames based on their colors
-  API returns a ZIP file containing the processed images

---

## 🛠️ Installation & Usage

### **1️⃣ Clone the Repository**

```sh
 git clone https://github.com/Luiso9/goseperate.git
 cd goseperate/backend
```

### **2️⃣ Build and Run with Docker**

#### **Build the Docker Image**

```sh
docker build -t goseperate-app .
```

#### **Run the Container**

```sh
docker run --rm -p 9330:9330 goseperate-app
```

-  API will be available at `http://localhost:9330`

---

## 🛠 API Endpoints

### **🔹 Upload Image & Extract Colors**

```http
POST /extract
```

**Request:** (multipart/form-data)

-  `file`: Image file
-  `colors`: Number of colors to extract

**Response:**

-  JSON with extracted color filenames
-  ZIP file containing processed images

Example:

```json
{
	"extracted_files": ["color_1.png", "color_2.png"]
}
```

---

## ⚡ Contributing

Feel free to submit issues and pull requests! 🚀
