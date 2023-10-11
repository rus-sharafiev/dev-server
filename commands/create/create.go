package create

import (
	"archive/zip"
	"fmt"
	"log"
	"os"
	"path"
)

func unzip(source string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	source = path.Join(homeDir, ".dev", "draft", source)

	read, err := zip.OpenReader(source)
	if err != nil {
		return err
	}
	defer read.Close()

	for _, file := range read.File {
		if file.Mode().IsDir() {
			continue
		}
		open, err := file.Open()
		if err != nil {
			return err
		}

		os.MkdirAll(path.Dir(file.Name), os.ModeDir)
		create, err := os.Create(file.Name)
		if err != nil {
			return err
		}
		defer create.Close()
		create.ReadFrom(open)
	}
	return nil
}

func Build() {
	if _, err := os.Stat("build"); err == nil {
		fmt.Println("Build dir already exists")
		return
	}

	if err := unzip("build.zip"); err != nil {
		log.Fatal(err)
	}
}

func Src() {
	if _, err := os.Stat("src"); err == nil {
		fmt.Println("Src dir already exists")
		return
	}

	if err := unzip("src.zip"); err != nil {
		log.Fatal(err)
	}
}
