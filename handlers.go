package main

import (
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

	// Write file
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

	//TODO: Database management
	// Desired Behavior:
	// Holds the the metadata for all uploaded images (pk hash, original name, new name, size, upload date)
	// Should check the new name if the generated name has already been used, if it has get a new one
	// Should check if the newly uploaded image's hash is already in the database, if it is, just return that instead
	// Do I want to deal with user accounts and remembering passed uploads?

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
