package imup

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"io"
	"mime/multipart"
	"net/http/httptest"
	"os"
	"testing"
)

const (
	testPNGName = "testpng.png"
	testTXTName = "testtxt.txt"
)

// createTestPNG creates the test PNG image.
func createTestPNG() error {
	// Create a new image.
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))

	// Draw a square.
	for x := 10; x < 90; x++ {
		for y := 10; y < 90; y++ {
			img.Set(x, y, color.RGBA{255, 0, 0, 255})
		}
	}

	// Create the test PNG file.
	file, err := os.Create(testPNGName)
	if err != nil {
		return err
	}
	defer file.Close()

	// Save the image.
	png.Encode(file, img)

	return nil
}

// deleteTestPNG deletes the test PNG image.
func deleteTestPNG() {
	os.Remove(testPNGName)
}

// createTestTXT create the test TXT file.
func createTestTXT() error {
	// Create the test TXT file.
	file, err := os.Create(testTXTName)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write data to the file.
	file.Write([]byte(`Lorem ipsum dolor sit amet`))

	return nil
}

// deleteTestTXT deletes the test TXT file.
func deleteTestTXT() {
	os.Remove(testTXTName)
}

func TestNew(t *testing.T) {
	err := createTestPNG()
	if err != nil {
		t.Fatal(err)
	}
	defer deleteTestPNG()

	// Get file handler for test PNG image.
	file, err := os.Open(testPNGName)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	// Create a multipart writer.
	var b bytes.Buffer
	mw := multipart.NewWriter(&b)

	// Create a new form file.
	fw, err := mw.CreateFormFile("file", testPNGName)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = io.Copy(fw, file); err != nil {
		t.Fatal(err)
	}
	mw.Close()

	// Create a new test request.
	req := httptest.NewRequest("POST", "http://127.0.0.1/", &b)

	// Set the multipart/form-data Content-Type.
	req.Header.Set("Content-Type", mw.FormDataContentType())

	// Upload the image.
	_, err = New("file", req, &Options{
		MaxFileSize:  1024,
		AllowedTypes: PopularTypes,
	})
	if err != nil {
		t.Fatal(err)
	}
}
