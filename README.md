# imup [![GoDoc](http://img.shields.io/badge/godoc-reference-blue.svg)](http://godoc.org/github.com/beeker1121/imup) [![License](http://img.shields.io/badge/license-mit-blue.svg)](https://raw.githubusercontent.com/beeker1121/imup/master/LICENSE) [![Go Report Card](https://goreportcard.com/badge/github.com/beeker1121/imup)](https://goreportcard.com/report/github.com/beeker1121/imup) [![Build Status](https://travis-ci.org/beeker1121/imup.svg?branch=master)](https://travis-ci.org/beeker1121/imup)

imup is an image upload handler written in Go.

Managing image uploads over HTTP in Go can be difficult when taking into account file type and request length checks, handling request cancellation, and so on. This package was built to abstract those details away and provide a simple API to make dealing with image uploads easier.

**Special Thanks** to [@vcabbage](https://github.com/vcabbage) for the brilliant max file size solution!

This project was created for [MailDB.io](https://maildb.io/), check us out!

## Features

- Handles all MIME image types defined in the standard `http` lib:  
  GIF, PNG, JPEG, BMP, WEBP, ICO
- Set allowable image formats
- Supports max file size limit
- Reads data only up to max file size limit, will not eat up bandwidth
- Handles spoofed Content-Type header

## Installation

Fetch the package from GitHub:

```sh
go get github.com/beeker1121/imup
```

Import to your project:

```go
import "github.com/beeker1121/imup"
```

## Usage

```go
func handler(w http.ResponseWriter, r *http.Request) {
	// Parse the uploaded file.
	ui, err = imup.New("file", r, &imup.Options{
		MaxFileSize:  1 * 1024 * 1024,   // 1 MB
		AllowedTypes: imup.PopularTypes,
	})
	if err != nil {
		...
	}

	// Save the image.
	filename, err := ui.Save("images/test")
	if err != nil {
		...
	}

	fmt.Println(filename) // images/test.png
}
```

## License

MIT license