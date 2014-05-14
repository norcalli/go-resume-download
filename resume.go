package resume

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

func NewRequestAtOffset(method, url string, body io.Reader, offset int64) (*http.Request, error) {
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Add("Range", fmt.Sprintf("bytes=%d-", offset))
	return request, err
}

// func DownloadToFile(client *http.Client, url, filename string) (*os.File, *http.Response, error) {
func DoAndCreateFile(client *http.Client, request *http.Request, filename string) (*os.File, *http.Response, error) {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
	if err != nil {
		return nil, nil, err
	}
	response, err := DondAppendToWriter(client, request, file)
	if err != nil {
		file.Close()
		return nil, nil, err
	}
	return file, response, nil
}

// func DownloadToFile(client *http.Client, url, filename string) (*os.File, *http.Response, error) {
func GetAndCreateFile(client *http.Client, url, filename string) (*os.File, *http.Response, error) {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
	if err != nil {
		return nil, nil, err
	}
	response, err := GetAndAppendToWriter(client, url, file)
	if err != nil {
		file.Close()
		return nil, nil, err
	}
	return file, response, nil
}

// func DownloadToWriter(client *http.Client, url string, file *os.File) (*os.File, *http.Response, error) {
func DoAndAppendToWriter(client *http.Client, request *http.Request, file *os.File) (*http.Response, error) {
	// Assert that the file is at the end and get fileSize at the same time.
	fileSize, err := file.Seek(0, 2)
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	} else {
		request.Header.Add("Range", fmt.Sprintf("bytes=%d-", fileSize))
	}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	// TODO: How to handle stuff that doesn't accept ranges?
	// if response.StatusCode != 206 {
	if response.Header.Get("Accept-Ranges") == "" {
		response.Body.Close()
		// file.Seek(0, 0)
		return nil, errors.New("Url doesn't accept Ranges (resuming download).")
	}
	return response, nil
}

// func DownloadToWriter(client *http.Client, url string, file *os.File) (*os.File, *http.Response, error) {
func GetAndAppendToWriter(client *http.Client, url string, file *os.File) (*http.Response, error) {
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	// Assert that the file is at the end and get fileSize at the same time.
	fileSize, err := file.Seek(0, 2)
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	} else {
		request.Header.Add("Range", fmt.Sprintf("bytes=%d-", fileSize))
	}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	// TODO: How to handle stuff that doesn't accept ranges?
	// if response.StatusCode != 206 {
	if response.Header.Get("Accept-Ranges") == "" {
		response.Body.Close()
		// file.Seek(0, 0)
		return nil, errors.New("Url doesn't accept Ranges (resuming download).")
	}
	return response, nil
}
