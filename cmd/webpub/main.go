package main

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/readium/r2-streamer-go/parser"
)

func main() {
	if len(os.Args) <= 1 {
		fmt.Fprintln(os.Stderr, "usage:")
		fmt.Fprintln(os.Stderr, "   1: create manifest file")
		fmt.Fprintln(os.Stderr, "      webpub PATH")
		fmt.Fprintln(os.Stderr, "   2: create a copy of ePub with embeded manifest.json")
		fmt.Fprintln(os.Stderr, "      webpub INPUTPATH OUTPUTPATH")
		os.Exit(1)
	}

	var outputFileName string

	inputPath := os.Args[1]

	if len(os.Args) > 2 {
		outputFileName = os.Args[2]
	}

	publication, err := parser.Parse(string(inputPath))

	if err != nil {
		panic(err)
	}

	jsonData, jsonError := json.MarshalIndent(publication, "", "  ")

	if jsonError != nil {
		panic(jsonError)
	}

	if len(outputFileName) > 0 {
		err = embedManifest(inputPath, jsonData, outputFileName)
		if err != nil {
			panic(err)
		}
	} else {
		fmt.Println(string(jsonData))
	}
}

func embedManifest(inputPath string, jsonManifest []byte, outputFileName string) (err error) {

	// create a zip reader
	zipReader, err := zip.OpenReader(inputPath)
	if err != nil {
		return
	}
	defer zipReader.Close()

	newZipFile, err := os.Create(outputFileName)
	if err != nil {
		return
	}
	defer newZipFile.Close()

	zipWriter := zip.NewWriter(newZipFile)
	defer zipWriter.Close()

	manifestFileName := "META-INF/manifest.json"
	manifestWriter, err := zipWriter.Create(manifestFileName)
	if err != nil {
		return
	}

	_, err = manifestWriter.Write(jsonManifest)
	if err != nil {
		return
	}

	err = copyZipFile(zipWriter, zipReader, manifestFileName)
	if err != nil {
		return
	}

	return
}

func copyZipFile(out *zip.Writer, in *zip.ReadCloser, excludeFileName string) error {
	for _, file := range in.File {

		if file.Name == excludeFileName {
			continue
		}

		// copy header, otherwise 'io.Copy(newFile, r)' will fail with an error "zip: checksum error"
		header := file.FileHeader

		newFile, err := out.CreateHeader(&header)
		if err != nil {
			return err
		}

		r, err := file.Open()
		if err != nil {
			return err
		}

		_, err = io.Copy(newFile, r)
		r.Close()
		if err != nil {
			return err
		}
	}

	return nil
}
