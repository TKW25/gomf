package main

import (
	"fmt"
	"io"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func ReceiveFile(w http.ResponseWriter, r *http.Request) {
	// Receive file
	file, header, err := r.FormFile("FileFormName")
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	// Check if the uploaded file is larger than what we want to allow to be uploaded
	if header.Size > Config.MaxFileSize {
		log.Println("Uploaded file exceeded the maximum allowed file size")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	var mimeBuf [512]byte // http.DetectContentType will only use the first 512 bytes
	_, err = io.ReadFull(file, mimeBuf[:])

	if err != nil && err != io.ErrUnexpectedEOF && err != io.EOF {
		log.Println("An unexpected error occured while reading the file for mimeType detection")
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// get new filename and find the extension for it
	var newFileName string
	if ext := getExtension(strings.Split(header.Filename, "."), mimeBuf[:]); ext != "" {
		_, err := file.Seek(0, 0)
		if err != nil {
			log.Println("Failed resetting file")
			log.Println(err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		newFileName, _ = GetRandomFileName(ext)
		fmt.Println(newFileName)
	} else {
		log.Println("Failed to get an extension")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	//TODO: Database management

	// Create file
	var fp string = filepath.Join(Config.UploadDir, newFileName)
	fo, err := os.Create(fp)
	if err != nil {
		log.Println("Failed creating file")
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	defer fo.Close()

	if runtime.GOOS != "windows" { // Chmod is finicky on windows, don't do it.  Primarily for development.
		if err := fo.Chmod(0777); err != nil {
			log.Println("Error changing file permissions")
			log.Println(err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}

	// Write file
	if n, err := io.Copy(fo, file); err != nil {
		_ = os.Remove(fp)
		log.Println("Error copying file")
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	} else if n != header.Size {
		_ = os.Remove(fp)
		log.Printf("Expect %v but only read %v", header.Size, n)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Connection", "close")
	w.WriteHeader(http.StatusOK)
	// TODO: Set response URL
}

func getExtension(fileName []string, buf []byte) (ext string) {
	var fileExt string = ""
	if len(fileName) > 1 {
		fileExt = "." + fileName[len(fileName)-1]
	}
	mimeExt, err := mime.ExtensionsByType(http.DetectContentType(buf))

	if err != nil {
		log.Println(err)
		return ""
	} else if mimeExt == nil && fileExt == "" {
		log.Println("ERROR: Found no possible file extensions")
		//TODO: What do I do in this case?
		return ""
	} else if mimeExt == nil {
		log.Println("WARN: Failed to find a valid mimeType, using passed in extention")
		ext = fileExt
	} else if fileExt == "" {
		if len(mimeExt) == 1 {
			log.Println("Using detected extension")
			ext = mimeExt[0]
		} else {
			log.Println("WARN: Found multiple potential mimeTypse and no fileExt")
			//TODO: What should be done here
			return ""
		}
	} else {
		log.Println("Using provided extension")
		ext = fileExt
	}
	return ext
}
