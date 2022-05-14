package main

import (
	"archive/zip"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync"
)

func appendFile(filename string, zipw *zip.Writer) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("Failed to open %s: %s", filename, err)
	}
	defer file.Close()

	wr, err := zipw.Create(filename)
	if err != nil {
		msg := "Failed to create entry for %s in zip file: %s"
		return fmt.Errorf(msg, filename, err)
	}

	if _, err := io.Copy(wr, file); err != nil {
		return fmt.Errorf("Failed to write %s to zip: %s", filename, err)
	}

	return nil
}

func appendAllFiles(s []string, zipw *zip.Writer, wg *sync.WaitGroup) {
	defer wg.Done()

	for i := 0; i < len(s); i++ {
		if err := appendFile(s[i], zipw); err != nil {
			log.Fatalf("Failed to add file %s to zip: %s", s[i], err)
		}
	}
}

func getFilesFromCmd() []string {
	var names string
	flag.StringVar(&names, "src", "guest", "names")
	flag.Parse()
	s := strings.Split(names, ",")
	return s
}

func main() {
	var wg sync.WaitGroup

	flags := os.O_WRONLY | os.O_CREATE | os.O_TRUNC
	file, err := os.OpenFile("archive.zip", flags, 0644)
	if err != nil {
		log.Fatalf("Failed to open zip for writing: %s", err)
	}
	defer file.Close()

	zipw := zip.NewWriter(file)
	defer zipw.Close()

	wg.Add(1)

	go appendAllFiles(getFilesFromCmd(), zipw, &wg)

	fmt.Println("Waiting for goroutines to finish...")
	wg.Wait()
	fmt.Println("Done!")
}
