package saver

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/disintegration/imaging"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

type sizeError struct {
	text string
}

func (err sizeError) Error() string {
	return err.text
}

type formatError struct {
	expected string
	real     string
}

func (err formatError) Error() string {
	return fmt.Sprintf("format error: expected %s, but was %s", err.expected, err.real)
}

// create 100x100 preview file for image
func createThumbnail(imageFileName string, previewFileName string, width int, height int) (err error) {
	if width < 1 || height < 1 {
		return &sizeError{"thumbnail size error"}
	}

	var img image.Image
	img, err = imaging.Open(imageFileName)
	if err != nil {
		return
	}
	thumb := imaging.Thumbnail(img, width, height, imaging.CatmullRom)
	err = imaging.Save(thumb, previewFileName)
	if err != nil {
		return
	}
	return
}

// if the last char of dir isn't '/', add it to the end of it
func buildDestinationPath(directory string, filename string) (path string) {
	filename = strings.Replace(filename, "/", "", -1)
	if filename == "" {
		filename = "undefined"
	}
	if directory == "" {
		return filename
	} else {
		path = fmt.Sprintf("%s/%s", directory, filename)
	}
	path = strings.Replace(path, "//", "/", -1)
	return
}

// save file
func SaveImageFromPart(sourceFile *multipart.Part, destinationDirectory string) (err error) {
	var destinationFile *os.File

	// build filepath
	destinationFilePath := buildDestinationPath(destinationDirectory, sourceFile.FileName())

	log.Printf("SAVE FILE TO %s\n", destinationFilePath)
	// create a new file in the desired location
	destinationFile, err = os.Create(destinationFilePath)
	defer func() { _ = destinationFile.Close() }()
	if err != nil {
		return
	}

	// copy the source file to the destination file
	if _, err = io.Copy(destinationFile, sourceFile); err != nil {
		return
	}

	log.Printf("SAVED FROM MULTIPART")

	// create 100x100 preview
	err = createThumbnail(destinationFilePath, strings.Replace(destinationFilePath, sourceFile.FileName(), "min_"+sourceFile.FileName(), 1), 100, 100)
	if err != nil {
		_ = os.Remove(destinationFilePath)
		return
	}

	return
}

// save image file and its preview from base64 source
func SaveImageFromBase64(sourceBase64 string, destinationDirectory string) (err error) {
	log.Println("Save base64")

	// parse base64 string
	// match 1: image type - png, jpg or gif.
	// match 2: base64 code
	re := regexp.MustCompile(`data:image/([a-z]{3,4});base64,(.*)`)

	// save matches to new variable
	b64Groups := re.FindStringSubmatch(sourceBase64)

	// if there is no matches, return error
	if len(b64Groups) == 0 {
		return &formatError{"base64", "undefined"}
	}

	// decode base64 string
	var b []byte
	b, err = base64.StdEncoding.DecodeString(b64Groups[2])
	if err != nil {
		return
	}

	// create reader to decode bytes to image
	r := bytes.NewReader(b)
	var pic image.Image

	// use match 1 to decode with right package depends on image type
	switch b64Groups[1] {
	case "png":
		pic, err = png.Decode(r)
	case "jpg", "jpeg":
		pic, err = jpeg.Decode(r)
	case "gif":
		pic, err = gif.Decode(r)
	default:
		log.Printf("FILE TYPE ERROR: %s\n", b64Groups[1])
		return &formatError{"image(png/jpg/gif)", b64Groups[1]}
	}
	if err != nil {
		return
	}

	// create filename
	fileName := fmt.Sprintf("%d%s.%s", time.Now().Unix(), b64Groups[2][:10], b64Groups[1])
	// build filepath
	destinationFilePath := buildDestinationPath(destinationDirectory, fileName)

	log.Printf("SAVE FILE TO %s\n", destinationFilePath)

	// open file to write
	destinationFile, err := os.OpenFile(destinationFilePath, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		return
	}
	defer destinationFile.Close()

	// write to file using right encoder
	switch b64Groups[1] {
	case "png":
		err = png.Encode(destinationFile, pic)
	case "jpg", "jpeg":
		err = jpeg.Encode(destinationFile, pic, nil)
	case "gif":
		err = gif.Encode(destinationFile, pic, nil)
	}
	if err != nil {
		return
	}

	log.Println("SAVED FROM BASE64")

	// create 100x100 preview
	err = createThumbnail(destinationFilePath, strings.Replace(destinationFilePath, fileName, "min_"+fileName, 1), 100, 100)
	if err != nil {
		os.Remove(destinationFilePath)
		return
	}
	return
}

func SaveImageFromUrl(url string, destinationDirectory string) (err error) {
	var response *http.Response
	response, err = http.Get(url)
	if err != nil {
		return
	}
	defer response.Body.Close()

	contentType := response.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/") {
		return errors.New("content-type error")
	}

	splittedUrl := strings.Split(url, "/")
	fileName := splittedUrl[len(splittedUrl)-1]
	if !strings.ContainsRune(fileName, '.') {
		fileName = fmt.Sprintf("%s.%s", fileName, strings.SplitN(contentType, "/", 2)[1])
	}

	// build filepath
	destinationFilePath := buildDestinationPath(destinationDirectory, fileName)

	log.Printf("SAVE FILE TO %s\n", destinationFilePath)

	// open file to write
	destinationFile, err := os.OpenFile(destinationFilePath, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		return
	}
	defer destinationFile.Close()

	// Use io.Copy to just dump the response body to the file. This supports huge files
	_, err = io.Copy(destinationFile, response.Body)
	if err != nil {
		return
	}

	log.Println("SAVED FROM URL")

	// create 100x100 preview
	err = createThumbnail(destinationFilePath, strings.Replace(destinationFilePath, fileName, "min_"+fileName, 1), 100, 100)
	if err != nil {
		os.Remove(destinationFilePath)
		return
	}
	return
}
