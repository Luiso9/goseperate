package services

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math/rand"
	"os"
	"path/filepath"

	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/gonum/stat"
	"gocv.io/x/gocv"
)

// ExtractColors processes an image and saves layers with hex filenames
func ExtractColors(imagePath, outputDir string, k int) error {
	// Load the image
	img := gocv.IMRead(imagePath, gocv.IMReadColor)
	if img.Empty() {
		return fmt.Errorf("could not read image")
	}
	defer img.Close()

	rows, cols := img.Rows(), img.Cols()
	data := make([]float64, 0, rows*cols*3)

	// Convert image pixels to data points
	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			pixel := img.GetVecbAt(y, x) // BGR format
			data = append(data, float64(pixel[2]), float64(pixel[1]), float64(pixel[0])) // Convert to RGB
		}
	}

	// Run K-Means clustering
	points := mat.NewDense(len(data)/3, 3, data)
	centroids, labels := KMeans(points, k)

	// Create output directory
	os.MkdirAll(outputDir, os.ModePerm)

	// Generate extracted layers
	for i, center := range centroids {
		layer := image.NewRGBA(image.Rect(0, 0, cols, rows))

		// Convert to HEX color
		hexColor := fmt.Sprintf("#%02X%02X%02X", uint8(center[0]), uint8(center[1]), uint8(center[2]))

		for y := 0; y < rows; y++ {
			for x := 0; x < cols; x++ {
				index := y*cols + x
				if index < len(labels) && labels[index] == i {
					layer.Set(x, y, color.RGBA{
						R: uint8(center[0]), G: uint8(center[1]), B: uint8(center[2]), A: 255,
					})
				}
			}
		}

		// Save with hex color name
		filePath := filepath.Join(outputDir, fmt.Sprintf("%s.png", hexColor))
		file, _ := os.Create(filePath)
		png.Encode(file, layer)
		file.Close()
	}

	return nil
}

// KMeans runs K-Means clustering on image data
func KMeans(data *mat.Dense, k int) ([][]float64, []int) {
	rows, _ := data.Dims()
	labels := make([]int, rows)
	centroids := make([][]float64, k)

	// Initialize centroids randomly
	for i := range centroids {
		idx := rand.Intn(rows)
		centroids[i] = mat.Row(nil, idx, data)
	}

	// Iterate clustering
	for iter := 0; iter < 10; iter++ {
		// Assign points to nearest centroid
		for i := 0; i < rows; i++ {
			sample := mat.Row(nil, i, data)
			labels[i] = nearestCentroid(sample, centroids)
		}

		// Update centroids
    for i := 0; i < 3; i++ {
      colVec := newCenter.colView(i)
      colSlice := make([]float64, colVec.Len())

      for j := 0; j < colVec.Len(); j++ {
        colSlice[j] = colVec.AtVec(j)
      }

      centroids[i] = stat.Mean(colSlice, nil)
    }
  }
	return centroids, labels
}

// nearestCentroid finds the closest centroid
func nearestCentroid(sample []float64, centroids [][]float64) int {
	bestIdx := 0
	bestDist := euclideanDistance(sample, centroids[0])

	for i, center := range centroids[1:] {
		dist := euclideanDistance(sample, center)
		if dist < bestDist {
			bestIdx = i + 1
			bestDist = dist
		}
	}

	return bestIdx
}

// euclideanDistance calculates the distance between two points
func euclideanDistance(a, b []float64) float64 {
	sum := 0.0
	for i := range a {
		diff := a[i] - b[i]
		sum += diff * diff
	}
	return sum
}
