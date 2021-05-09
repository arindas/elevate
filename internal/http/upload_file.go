package http

import (
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"sync"

	"github.com/arindas/elevate/internal/app"
)

type uploadFileRequest struct {
	FileHeaders     []*multipart.FileHeader
	TargetDirectory string

	ResponseWriter http.ResponseWriter
	Request        *http.Request

	err    error
	status int
}

func uploadFileRequestInstance(w http.ResponseWriter, r *http.Request) uploadFileRequest {
	return uploadFileRequest{ResponseWriter: w, Request: r}
}

func (req *uploadFileRequest) verifyMethod() {
	if req.err != nil {
		return
	}

	if req.Request.Method != http.MethodPost {
		req.err = fmt.Errorf("request with method %v. expected %v",
			req.Request.Method, http.MethodPost)
		req.status = http.StatusBadRequest
	}
}

func (req *uploadFileRequest) parseRequest() {
	if req.err != nil {
		return
	}

	req.err = req.Request.ParseMultipartForm(32 << 20)
	req.status = http.StatusBadRequest
}

func (req *uploadFileRequest) obtainTargetDirectory(baseDirectory string) {
	if req.err != nil {
		return
	}

	upload_paths := req.Request.MultipartForm.Value["file_upload_path"]
	if upload_paths != nil {
		req.TargetDirectory = upload_paths[0]
	} else {
		req.TargetDirectory = ""
	}

	req.TargetDirectory = path.Join(baseDirectory, req.TargetDirectory)
	req.err = os.MkdirAll(req.TargetDirectory, os.ModePerm)
	req.status = http.StatusFailedDependency
}

func (req *uploadFileRequest) obtainSelectedFiles() {
	if req.err != nil {
		return
	}

	req.FileHeaders = req.Request.MultipartForm.File["selected_files"]
	if len(req.FileHeaders) != 0 {
		return
	}

	req.err = http.ErrMissingFile
	req.status = http.StatusBadRequest
}

func writeFileToDisk(fh *multipart.FileHeader, targetDirectory string) error {
	filePath := path.Join(targetDirectory, fh.Filename)

	fileOnDisk, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer fileOnDisk.Close()

	var fileInReq multipart.File
	fileInReq, err = fh.Open()
	if err != nil {
		return err
	}
	defer fileInReq.Close()

	_, err = io.Copy(fileOnDisk, fileInReq)

	return err
}

func (req *uploadFileRequest) collectResults(
	filesWritten <-chan string,
	errorChannel <-chan error) {

	for f := range filesWritten {
		log.Printf("File %s written to disk.", f)
	}

	failed := false
	errors := make([]string, 0)

	for e := range errorChannel {
		failed = true
		errors = append(errors, e.Error())
		log.Printf("error: %v", e.Error())
	}

	if failed {
		req.err = fmt.Errorf(
			"errors during writing files: %v",
			errors)
		req.status = http.StatusFailedDependency
	}
}

func (req *uploadFileRequest) commit() {
	if req.err != nil {
		return
	}

	filesWritten := make(chan string,
		len(req.FileHeaders))
	errorChannel := make(chan error,
		len(req.FileHeaders))

	var wg sync.WaitGroup

	for _, fh := range req.FileHeaders {
		wg.Add(1)

		go func(fileHeader *multipart.FileHeader) {
			defer wg.Done()
			err := writeFileToDisk(
				fileHeader, req.TargetDirectory)
			if err != nil {
				errorChannel <- err
			} else {
				filesWritten <- fileHeader.Filename
			}
		}(fh)
	}
	wg.Wait()

	req.collectResults(filesWritten, errorChannel)
}

func (req *uploadFileRequest) respond() {
	if req.err != nil {
		http.Error(req.ResponseWriter, req.err.Error(), req.status)
		log.Printf("error: %s", req.err.Error())
	}

	req.ResponseWriter.Write([]byte("Files uploaded!"))
}

func UploadFileHandler(app app.AppConfig) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req := uploadFileRequestInstance(w, r)

		req.verifyMethod()
		req.parseRequest()
		req.obtainTargetDirectory(app.BaseDirectory)
		req.obtainSelectedFiles()
		req.commit()
		req.respond()
	})
}
