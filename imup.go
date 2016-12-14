package imup

import (
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

// ImageTypes defines the allowed types for an uploaded image.
type ImageTypes map[string]bool

var (
	GIF  ImageTypes = ImageTypes{"image/gif": true}
	PNG  ImageTypes = ImageTypes{"image/png": true}
	JPG  ImageTypes = ImageTypes{"image/jpg": true}
	JPEG ImageTypes = ImageTypes{"image/jpeg": true}
	BMP  ImageTypes = ImageTypes{"image/bmp": true}
	WEBP ImageTypes = ImageTypes{"image/webp": true}
	ICO  ImageTypes = ImageTypes{"image/vnd.microsoft.icon": true}

	// Convenience type for popular web image types.
	PopularTypes ImageTypes = ImageTypes{
		"image/gif":  true,
		"image/png":  true,
		"image/jpg":  true,
		"image/jpeg": true,
	}

	// Convenience type for all image types.
	AllTypes ImageTypes = ImageTypes{
		"image/gif":                true,
		"image/png":                true,
		"image/jpg":                true,
		"image/jpeg":               true,
		"image/bmp":                true,
		"image/webp":               true,
		"image/vnd.microsoft.icon": true,
	}
)

// UploadedImage defines the uploaded image.
type UploadedImage struct {
	Size   int64
	Type   string
	file   multipart.File
	header *multipart.FileHeader
}

// Options defines the available options for the image upload.
type Options struct {
	MaxFileSize  int64
	AllowedTypes ImageTypes
}

// New returns a new UploadedImage object if the uploaded file could be parsed
// and validated as an image, otherwise it returns an error.
//
// Closing of the saved file must be handled by the user.
func New(key string, r *http.Request, opts *Options) (*UploadedImage, error) {
	var err error
	ui := &UploadedImage{}

	// Try to parse the multipart file from the request.
	if ui.file, ui.header, err = r.FormFile(key); err != nil {
		return nil, err
	}

	// Check if type is allowed.
	if len(opts.AllowedTypes) > 0 {
		if !isTypeAllowed(ui, opts.AllowedTypes) {
			return nil, ErrDisallowedType
		}
	}

	// Check file size.
	if ui.Size, err = ui.file.Seek(0, 2); err != nil {
		return nil, err
	}
	if ui.Size > opts.MaxFileSize {
		return nil, ErrFileSize
	}

	// Reset file pointer.
	if _, err = ui.file.Seek(0, 0); err != nil {
		return nil, err
	}

	return ui, nil
}

// Save saves the uploaded image to the given location and returns the saved
// file location with the file extension added on.
func (ui *UploadedImage) Save(filename string) (string, error) {
	// Handle the file extension.
	var ext string
	switch ui.Type {
	case "image/gif":
		ext = ".gif"
	case "image/png":
		ext = ".png"
	case "image/jpg":
		ext = ".jpg"
	case "image/jpeg":
		ext = ".jpeg"
	case "image/bmp":
		ext = ".bmp"
	case "image/webp":
		ext = ".webp"
	case "image/vnd.microsoft.icon":
		ext = ".ico"
	}
	filename += ext

	// Create output file.
	out, err := os.Create(filename)
	if err != nil {
		return "", err
	}
	defer out.Close()

	// Store the uploaded image to output file.
	_, err = io.Copy(out, ui.file)
	if err != nil {
		return "", err
	}

	return filename, nil
}

// Close closes the uploaded image.
func (ui *UploadedImage) Close() {
	ui.file.Close()
}

// isTypeAllowed checks if the given file type is allowed.
func isTypeAllowed(img *UploadedImage, types ImageTypes) bool {
	// Get up to the first 512 bytes of data.
	b := make([]byte, 512)
	_, err := img.file.Read(b)
	if err != nil {
		return false
	}

	// Reset file pointer.
	if _, err = img.file.Seek(0, 0); err != nil {
		return false
	}

	// Try to detect the file type.
	img.Type = http.DetectContentType(b)

	// Validate type.
	if _, ok := types[img.Type]; ok {
		return true
	}

	return false
}
