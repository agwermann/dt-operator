package pkg

import (
	"fmt"
	"log"
	"os"
	"strings"
)

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
