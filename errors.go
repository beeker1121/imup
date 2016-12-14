package imup

import "errors"

var (
	// ErrDisallowedType is returned when the uploaded
	// file type is not allowed.
	ErrDisallowedType = errors.New("Image type is not allowed")

	// ErrFileSize is returned when the uploaded file
	// size exceeds the max file size limit.
	ErrFileSize = errors.New("File size exceeds max allowed size")
)
