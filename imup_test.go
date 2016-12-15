package imup

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"io"
	"mime/multipart"
	"net/http"
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

// createRequest creates a test request.
func createRequest(filename string) (*http.Request, error) {
	// Get file handler.
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Create a multipart writer.
	var b bytes.Buffer
	mw := multipart.NewWriter(&b)

	// Create a new form file.
	fw, err := mw.CreateFormFile("file", filename)
	if err != nil {
		return nil, err
	}
	if _, err = io.Copy(fw, file); err != nil {
		return nil, err
	}
	mw.Close()

	// Create a new test request.
	req := httptest.NewRequest("POST", "http://127.0.0.1/", &b)

	// Set the multipart/form-data Content-Type.
	req.Header.Set("Content-Type", mw.FormDataContentType())

	return req, nil
}

func TestNewAndSave(t *testing.T) {
	err := createTestPNG()
	if err != nil {
		t.Fatal(err)
	}
	defer deleteTestPNG()

	// Create a new test request.
	req, err := createRequest(testPNGName)
	if err != nil {
		t.Fatal(err)
	}

	// Upload the image.
	ui, err := New("file", req, &Options{
		MaxFileSize:  1024,
		AllowedTypes: PopularTypes,
	})
	if err != nil {
		t.Fatal(err)
	}

	// Save the image.
	filename, err := ui.Save("testsave")
	if err != nil {
		t.Fatal(err)
	}

	if err = os.Remove(filename); err != nil {
		t.Fatal(err)
	}
}

func TestContentLength(t *testing.T) {
	err := createTestPNG()
	if err != nil {
		t.Fatal(err)
	}
	defer deleteTestPNG()

	// Create a new test request.
	req, err := createRequest(testPNGName)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Length", "1024")

	// Upload the image.
	_, err = New("file", req, &Options{
		MaxFileSize:  1024,
		AllowedTypes: PopularTypes,
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestErrFileSize(t *testing.T) {
	err := createTestPNG()
	if err != nil {
		t.Fatal(err)
	}
	defer deleteTestPNG()

	// Create a new test request.
	req, err := createRequest(testPNGName)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Length", "1025")

	// Upload the image.
	_, err = New("file", req, &Options{
		MaxFileSize:  1024,
		AllowedTypes: PopularTypes,
	})
	if err != ErrFileSize {
		t.Errorf("Expected ErrFileSize, got %s", err)
	}
}

func TestErrFileSizeSpoofed(t *testing.T) {
	err := createTestPNG()
	if err != nil {
		t.Fatal(err)
	}
	defer deleteTestPNG()

	// Create a new test request.
	req, err := createRequest(testPNGName)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Length", "128")

	// Upload the image.
	_, err = New("file", req, &Options{
		MaxFileSize:  128,
		AllowedTypes: PopularTypes,
	})
	if err != ErrFileSize {
		t.Errorf("Expected ErrFileSize, got %s", err)
	}
}

func TestErrDisallowedType(t *testing.T) {
	err := createTestTXT()
	if err != nil {
		t.Fatal(err)
	}
	defer deleteTestTXT()

	// Create a new test request.
	req, err := createRequest(testTXTName)
	if err != nil {
		t.Fatal(err)
	}

	// Upload the image.
	_, err = New("file", req, &Options{
		MaxFileSize:  1024,
		AllowedTypes: PopularTypes,
	})
	if err != ErrDisallowedType {
		t.Errorf("Expected ErrDisallowedType, got %s", err)
	}
}
