package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	dtdl "github.com/agwermann/dt-operator/cmd/cli/dtdl"
	"gopkg.in/yaml.v3"
)

func main() {
	allArgs := os.Args
	args := allArgs[1:]

	if len(args) < 2 {
		log.Fatal("Inform DTDL input and output folders path")
	}

	inputFolderPath := args[0]
	outputFolderPath := args[1]

	files, err := ioutil.ReadDir(inputFolderPath)

	if err != nil {
		log.Fatal(err)
	}

	PrepareOutputFolder(outputFolderPath)

	for _, file := range files {
		if !file.IsDir() {
			inputFilename := filepath.Join(inputFolderPath, file.Name())
			outputFileName := strings.Split(file.Name(), ".")[0]
			outputFilename := filepath.Join(outputFolderPath, outputFileName+".yaml")
			if IsJsonFile(inputFilename) {
				fileContent, err := os.ReadFile(inputFilename)
				if err != nil {
					log.Fatal(err)
				}

				twinDefinition := &dtdl.Interface{}
				err = json.Unmarshal(fileContent, twinDefinition)
				if err != nil {
					log.Fatal(err)
				}
				twinYaml, err := yaml.Marshal(twinDefinition)
				err = WriteToFile(outputFilename, twinYaml)
				if err != nil {
					log.Fatal(err)
				}
				return
			}
		}
	}

}

func IsJsonFile(filename string) bool {
	s := strings.Split(filename, ".")
	return len(s) > 1 && s[1] == "json"
}

func PrepareOutputFolder(dirname string) error {
	err := os.Mkdir(dirname, os.ModeDir)

	if err == nil {
		fmt.Println("Output directory " + dirname + " was created")
		return nil
	}

	if os.IsExist(err) {
		info, err := os.Stat(dirname)
		if err != nil {
			return err
		}
		if !info.IsDir() {
			log.Fatal("File is not a directory")
		} else {
			log.Default().Print("Folder already exists")
		}
	}

	return err
}

func WriteToFile(fileName string, data []byte) error {

	err := os.WriteFile(fileName, data, 0664)

	if err != nil {
		fmt.Printf("Error while opening file")
		return err
	}

	return nil
}
