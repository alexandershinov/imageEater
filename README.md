# imageEater

A simple REST API with one single method that loads images

## Installation

### Without Docker
```
git clone https://github.com/alexandershinov/imageEater
cd imageEater
go get -d -v ./...
go get github.com/stretchr/testify/assert
```
then
```
go build
./imageEater
```
or
```
go run ./
```

### With docker-compose
```
git clone https://github.com/alexandershinov/imageEater
cd imageEater
docker-compose run --build
```

## Use API

You can send requests to the API by url http://localhost:4000

This application contains one method (`POST /images`) for uploading images to the server. The method can be used in three different ways:
1. multipart/form-data request with `files` field which contents list of files;
1. request with JSON body with `base64` field contents list of images in base64 format;
1. request with JSON body with `urls` field contents list of urls of images;

You can combine options 2 and 3 using both fields in the request body.
When uploading a file to the server, a preview of an image of 100x100 is created.

