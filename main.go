package main

import (
	"errors"
	"fmt"
	"github.com/urfave/cli"
	"io/ioutil"
	"log"
	"os"
	"time"
)

type SharedIgnoreFile struct {
	file *os.File
}

func (s SharedIgnoreFile) GenerateIgnoreFiles() {
	for _, name := range s.IgnoreFilesList() {
		file, err := os.Create(name)
		if err != nil {
			log.Println(fmt.Sprintf("cannot generate %v file. Error: %v", name, err))
			continue
		}
		defer file.Close()

		ignoreContent, err := os.ReadFile(s.file.Name())
		if err != nil {
			log.Fatalln(errors.New(fmt.Sprintf("cannot read the content of .ignore file. Error: %v", err)))
		}

		file.Write(ignoreContent)
		file.Sync()
	}
}

func (s SharedIgnoreFile) IgnoreFilesList() []string {
	return []string{
		".dockerignore", ".gitignore",
	}
}

func GenerateIgnoreFiles(ctx *cli.Context) error {
	currentWd, err := os.Getwd()
	if err != nil {
		return err
	}

	var sharedIgnoreFile SharedIgnoreFile

	files, err := ioutil.ReadDir(currentWd)
	for _, dirEntry := range files {
		if dirEntry.Name() == ".ignore" {
			file, err := os.OpenFile(dirEntry.Name(), os.O_RDONLY, os.ModePerm)
			if err != nil {
				return errors.New("cannot open an .ignore file")
			}
			sharedIgnoreFile.file = file
		}
	}

	if sharedIgnoreFile.file == nil {
		return errors.New("cannot find an .ignore file in current directory")
	}

	sharedIgnoreFile.GenerateIgnoreFiles()
	return nil
}

func main() {
	startedAt := time.Now()

	app := &cli.App{
		Name:   "IgnoreMaker",
		Usage:  "generate ignore files from one shared source",
		Action: GenerateIgnoreFiles,
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatalln(err)
	}

	finishedAt := time.Now()
	log.Printf("Done! %v s.", finishedAt.Sub(startedAt).Seconds())
}
