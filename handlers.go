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
	"time"
)

func ReceiveFile(w http.ResponseWriter, r *http.Request) {
	//TODO: Doesn't seem like I can cleanly get docker working on this system, stand up VPS and do it there
	//TODO: I think there's a memory leak somewhere https://golang.org/doc/diagnostics.html
	// Receive file
	file, header, err := r.FormFile("FileFormName")
	if err != nil {
		log.Println(err)
		log.Println("received a malformed request")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	defer file.Close()

	var metaData MetaData

	metaData.hash, err = GetHash(&file)
	if err != nil {
		log.Println(err)
		log.Println("Failed generating the hash")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// Check if the hash already exists, if it does it's likely the same file
	if HasHash(metaData.hash) {
		log.Println(fmt.Sprintf("Uploaded file hash %v already exists on the server", metaData.hash))
		w.Header().Set("Connection", "close")
		w.WriteHeader(http.StatusOK)
		// TODO: Set response URL
		return
	}

	// Check if the uploaded file is larger than what we want to allow to be uploaded
	if header.Size > Config.MaxFileSize {
		log.Println("Uploaded file exceeded the maximum allowed file size")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	var mimeBuf [512]byte // http.DetectContentType will only use the first 512 bytes
	if _, err = io.ReadFull(file, mimeBuf[:]); err != nil && err != io.ErrUnexpectedEOF && err != io.EOF {
		log.Println("An unexpected error occured while reading the file for mimeType detection")
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// get new filename and find the extension for it
	var newFileName string
	ext := getExtension(header.Filename, mimeBuf[:])

	_, err = file.Seek(0, 0) // reset file to copy it
	if err != nil {
		log.Println("Failed resetting file")
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// Ensure there's no name collisions
	const maxRetries = 10
	try := 1
	for {
		newFileName, err = GetRandomFileName(ext)
		if err != nil {
			log.Println("Failed getting a random file name")
			log.Println(err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		if !HasName(newFileName) {
			break
		}

		if try == maxRetries {
			log.Println("Exceeded the maximum number of retries")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		try++
	}

	// Create file
	fp := filepath.Join(Config.UploadDir, newFileName)
	fo, err := os.Create(fp)
	if err != nil {
		log.Println("Failed creating file")
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	defer fo.Close()

	if runtime.GOOS != "windows" { // Chmod is finicky on windows, don't do it.  For development.
		if err := fo.Chmod(0777); err != nil {
			log.Println("Error changing file permissions")
			log.Println(err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}

	// Write file to disk
	if n, err := io.Copy(fo, file); err != nil {
		log.Println("Error copying file")
		log.Println(err)
		if err = os.Remove(fp); err != nil {
			log.Println(err)
		}
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	} else if n != header.Size {
		if err = os.Remove(fp); err != nil {
			log.Println(err)
		}
		log.Printf("Expect %v but only read %v\n", header.Size, n)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// Set the rest of the metadata and write it to the database
	metaData.date = time.Now()
	metaData.n_name = newFileName
	metaData.o_name = header.Filename
	metaData.size = header.Size
	WriteMetadata(metaData)

	w.Header().Set("Connection", "close")
	w.WriteHeader(http.StatusOK)
	// TODO: Set response URL
}

// getExtension either determines the extension through the mime type if possible,
// or simply uses the original file extension.
func getExtension(fileName string, buf []byte) (ext string) {
	mimeExt, _ := mime.ExtensionsByType(http.DetectContentType(buf))
	if mimeExt != nil && len(mimeExt) == 1 {
		ext = mimeExt[0]
	} else {
		fileExt := ""
		if parts := strings.Split(fileName, "."); len(parts) > 1 {
			fileExt = "." + parts[len(parts)-1]
		}
		ext = fileExt
	}

	return ext
}
