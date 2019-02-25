package main

import (
	"encoding/json"
	"github.com/alexandershinov/imageEater/saver"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"strings"
)

type JsonBody struct {
	Base64 []string `json:"base64"`
	Urls   []string `json:"urls"`
}

func LoadImages(w http.ResponseWriter, r *http.Request) {
	log.Println("LOAD IMAGES")
	var err error
	switch strings.Split(r.Header.Get("Content-Type"), ";")[0] {
	case "multipart/form-data":
		log.Println("MULTIPART")
		var reader *multipart.Reader
		reader, err = r.MultipartReader()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		for {
			var f *multipart.Part
			f, err = reader.NextPart()
			if err == io.EOF {
				break
			}
			if f.FileName() != "" {
				err = saver.SaveImageFromPart(f, config.FilesDir)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					log.Println(err.Error())
					return
				}
			}
		}
	case "application/json":
		log.Println("JSON")
		var jsonBody JsonBody
		err = json.NewDecoder(r.Body).Decode(&jsonBody)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		for _, base64file := range jsonBody.Base64 {
			err = saver.SaveImageFromBase64(base64file, config.FilesDir)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Println(err.Error())
				return
			}
		}
		for _, url := range jsonBody.Urls {
			log.Printf("TRY TO SAVE FILE FROM URL %s", url)
			err = saver.SaveImageFromUrl(url, config.FilesDir)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Println(err.Error())
				return
			}
		}
	default:
		log.Println(r.Header.Get("Content-Type"))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Write([]byte("ok"))
}
