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
	"gocv.io/x/gocv"
	"bytes"
)

func ExtractColors(imagePath, outputDir string, k int) ([]string, error) {
	img := gocv.IMRead(imagePath, gocv.IMReadUnchanged)
	if img.Empty() {
		return nil, fmt.Errorf("could not read image")
	}
	defer img.Close()

	rows, cols := img.Rows(), img.Cols()
	hasAlpha := img.Channels() == 4
	data := make([]float64, 0, rows*cols*3)
	alphaChannel := make([]uint8, rows*cols)

	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			pixel := img.GetVecbAt(y, x)
			data = append(data, float64(pixel[2]), float64(pixel[1]), float64(pixel[0]))
			if hasAlpha {
				alphaChannel[y*cols+x] = pixel[3]
			} else {
				alphaChannel[y*cols+x] = 255
			}
		}
	}

	points := mat.NewDense(len(data)/3, 3, data)
	centroids, labels := KMeans(points, k)

	_ = os.MkdirAll(outputDir, os.ModePerm)

	extractedFiles := []string{}

	for i, center := range centroids {
		layer := image.NewRGBA(image.Rect(0, 0, cols, rows))
		hexColor := fmt.Sprintf("%02X%02X%02X", uint8(center[0]), uint8(center[1]), uint8(center[2]))

		for y := 0; y < rows; y++ {
			for x := 0; x < cols; x++ {
				index := y*cols + x
				if index < len(labels) && labels[index] == i {
					layer.Set(x, y, color.RGBA{
						R: uint8(center[0]), G: uint8(center[1]), B: uint8(center[2]),
						A: alphaChannel[index], // Maintain transparency if available
					})
				}
			}
		}

		filePath := filepath.Join(outputDir, fmt.Sprintf("%s.png", hexColor))
		file, _ := os.Create(filePath)
		png.Encode(file, layer)
		file.Close()

		extractedFiles = append(extractedFiles, hexColor)
	}

	return extractedFiles, nil
}

// KMeans implementation using basic clustering
func KMeans(data *mat.Dense, k int) ([][]float64, []int) {
	rows, _ := data.Dims()
	labels := make([]int, rows)
	centroids := make([][]float64, k)

	for i := range centroids {
		idx := rand.Intn(rows)
		centroids[i] = mat.Row(nil, idx, data)
	}

	for iter := 0; iter < 10; iter++ {
		for i := 0; i < rows; i++ {
			sample := mat.Row(nil, i, data)
			labels[i] = nearestCentroid(sample, centroids)
		}

		newCenter := mat.NewDense(k, 3, nil)
		counts := make([]int, k)

		for i := 0; i < rows; i++ {
			label := labels[i]
			pixel := mat.Row(nil, i, data)

			for j := 0; j < 3; j++ {
				newCenter.Set(label, j, newCenter.At(label, j)+pixel[j])
			}
			counts[label]++
		}

		for i := 0; i < k; i++ {
			if counts[i] > 0 {
				for j := 0; j < 3; j++ {
					newCenter.Set(i, j, newCenter.At(i, j)/float64(counts[i]))
				}
			}
		}

		for i := range centroids {
			centroids[i] = mat.Row(nil, i, newCenter)
		}
	}

	return centroids, labels
}

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

func euclideanDistance(a, b []float64) float64 {
	sum := 0.0
	for i := range a {
		diff := a[i] - b[i]
		sum += diff * diff
	}
	return sum
}

func GeneratePreview(imagePath string, numColors int) ([]byte, error) {
	img := gocv.IMRead(imagePath, gocv.IMReadColor)
	if img.Empty() {
		return nil, nil
	}
	defer img.Close()

	samples := gocv.NewMat()
	defer samples.Close()
	img.ConvertTo(&samples, gocv.MatTypeCV32F)
	samples = samples.Reshape(1, img.Rows()*img.Cols())

	
	criteria := gocv.NewTermCriteria(gocv.Count|gocv.EPS, 10, 1.0)
	labels := gocv.NewMat()
	centers := gocv.NewMat()
	defer labels.Close()
	defer centers.Close()

	gocv.KMeans(samples, numColors, &labels, criteria, 3, gocv.KMeansRandomCenters, &centers)

	result := gocv.NewMatWithSize(img.Rows(), img.Cols(), gocv.MatTypeCV8UC3)
	defer result.Close()

	for i := 0; i < img.Rows(); i++ {
		for j := 0; j < img.Cols(); j++ {
			clusterIdx := labels.GetIntAt(i*img.Cols()+j, 0)
			rgb := centers.GetVecfAt(clusterIdx, 0)

			result.SetUCharAt3(i, j, uint8(rgb[2]), uint8(rgb[1]), uint8(rgb[0])) // BGR -> RGB
		}
	}

	finalImg, _ := result.ToImage()
	buf := new(bytes.Buffer)
	err := png.Encode(buf, finalImg)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
