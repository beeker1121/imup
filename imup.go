package imup

import (
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
)

// ImageTypes defines the allowed types for an uploaded image.
type ImageTypes []string

// Image types according to the MIME specification.
const (
	GIF  = "image/gif"
	PNG  = "image/png"
	JPEG = "image/jpeg"
	BMP  = "image/bmp"
	WEBP = "image/webp"
	ICO  = "image/vnd.microsoft.icon"
)

var (
	// Convenience type for popular image types.
	PopularTypes = ImageTypes{GIF, PNG, JPEG}

	// Convenience type for all image types.
	AllTypes = ImageTypes{GIF, PNG, JPEG, BMP, WEBP, ICO}
)

// UploadedImage defines an uploaded image.
type UploadedImage struct {
	Type   string
	file   multipart.File
	header *multipart.FileHeader
}

// Options defines the available options for an image upload.
type Options struct {
	MaxFileSize  int64
	AllowedTypes ImageTypes
}

// New returns a new UploadedImage object if the uploaded file could be parsed
// and validated as an image, otherwise it returns an error.
//
// The key parameter should refer to the name of the file input from the
// multipart form.
func New(key string, r *http.Request, opts *Options) (*UploadedImage, error) {
	var err error
	ui := &UploadedImage{}

	// Handle max file size.
	if opts.MaxFileSize > 0 {
		// Check Content-Length header.
		cl, _ := strconv.ParseInt(r.Header.Get("Content-Length"), 10, 64)
		if cl > opts.MaxFileSize {
			return nil, ErrFileSize
		}

		// Wrap r.Body with our limitReader.
		r.Body = newLimitReader(r.Body, opts.MaxFileSize)
	}

	// Try to parse the multipart file from the request.
	if ui.file, ui.header, err = r.FormFile(key); err != nil {
		return nil, err
	}

	// Check if type is allowed.
	if len(opts.AllowedTypes) > 0 {
		if err = isTypeAllowed(ui, opts.AllowedTypes); err != nil {
			return nil, err
		}
	}

	return ui, nil
}

// Save saves the uploaded image to the given location and returns the location
// with the correct image extension added on.
//
// The underlying multipart image file is automatically closed.
func (ui *UploadedImage) Save(filename string) (string, error) {
	// Handle the file extension.
	var ext string
	switch ui.Type {
	case GIF:
		ext = ".gif"
	case PNG:
		ext = ".png"
	case JPEG:
		ext = ".jpg"
	case BMP:
		ext = ".bmp"
	case WEBP:
		ext = ".webp"
	case ICO:
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
	if err = ui.Close(); err != nil {
		return "", err
	}

	return filename, nil
}

// Close closes an uploaded image.
func (ui *UploadedImage) Close() error {
	return ui.file.Close()
}

// isTypeAllowed checks if the given file type is allowed.
func isTypeAllowed(ui *UploadedImage, types ImageTypes) error {
	// Get up to the first 512 bytes of data.
	b := make([]byte, 512)
	_, err := ui.file.Read(b)
	if err != nil {
		return err
	}

	// Reset file pointer.
	if _, err = ui.file.Seek(0, 0); err != nil {
		return err
	}

	// Try to detect the file type.
	ui.Type = http.DetectContentType(b)

	// Validate type.
	for _, t := range types {
		if ui.Type == t {
			return nil
		}
	}

	return ErrDisallowedType
}

// limitReader defines our custom request body ReaderCloser type, which wraps
// the standard io.LimitedReader.
type limitReader struct {
	r *io.LimitedReader
	io.Closer
}

// newLimitReader creates a new limitReader.
func newLimitReader(r io.ReadCloser, maxSize int64) io.ReadCloser {
	return &limitReader{
		r:      &io.LimitedReader{r, maxSize + 1},
		Closer: r,
	}
}

// Read satisfies the io.Reader interface.
//
// ErrFileSize is returned when the limit is exceeded rather than io.EOF like
// the standard io.LimitedReader.
func (l *limitReader) Read(p []byte) (int, error) {
	n, err := l.r.Read(p)
	if l.r.N < 1 {
		return n, ErrFileSize
	}
	return n, err
}
