package services

import (
	"bytes"
	"io"
	"net/http"
	"encoding/json"
	"fmt"
	"os/exec"
)

func ExtractColors(imagePath, outputDir string, k int) ([]string, error) {
	cmd := exec.Command("./sklearn-env/bin/python3", "scripts/cluster.py", imagePath, fmt.Sprintf("%d", k), outputDir)
	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("failed to extract colors: %v", err)
	}

	var result struct {
		Extracted []string `json:"extracted"`
		Error     string   `json:"error,omitempty"`
	}

	if err := json.Unmarshal(out.Bytes(), &result); err != nil {
		return nil, fmt.Errorf("failed to parse Python output: %v", err)
	}

	if result.Error != "" {
		return nil, fmt.Errorf("python error: %s", result.Error)
	}

	return result.Extracted, nil
}

func GeneratePreview(imagePath string, numColors int) ([]byte, error) {
	pythonAPI := "http://localhost:5037/preview"

	payload := map[string]interface{}{
		"image_path": imagePath,
		"num_colors": numColors,
	}
	jsonData, _ := json.Marshal(payload)

	resp, err := http.Post(pythonAPI, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to call Python API: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("python error: %s", string(body))
	}

	imageData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read Python API response: %v", err)
	}

	return imageData, nil
}
