package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func check(e error) {
	if e != nil {
		log.Fatalf("error: %v", e)
	}
}

func listFiles(include string, exclude string) ([]string, error) {
	fileList := []string{}
	err := filepath.Walk(".", func(path string, f os.FileInfo, err error) error {
		if doesFileMatch(path, include, exclude) {
			fileList = append(fileList, path)
		}
		return nil
	})
	return fileList, err
}

func doesFileMatch(path string, include string, exclude string) bool {
	if fi, err := os.Stat(path); err == nil && !fi.IsDir() {
		includeRe := regexp.MustCompile(include)
		excludeRe := regexp.MustCompile(exclude)
		return includeRe.Match([]byte(path)) && !excludeRe.Match([]byte(path))
	}
	return false
}

func findString(path string, find string) ([]string, error) {
	if len(find) > 0 {
		read, readErr := ioutil.ReadFile(path)
		check(readErr)

		re := regexp.MustCompile(find)
		results := re.FindAllString(string(read), -1)
		
		if len(results) > 0 {
			return results, nil
		}
	}

	return []string{}, nil
}

func main() {
	include := os.Getenv("INPUT_INCLUDE")
	exclude := os.Getenv("INPUT_EXCLUDE")
	find := os.Getenv("INPUT_FIND")

	files, filesErr := listFiles(include, exclude)
	check(filesErr)

	finalResultArray := []string {}
	fileFoundCount := 0

	for _, path := range files {
		findingResult, findErr := findString(path, find)
		check(findErr)

		if len(findingResult) > 0 {
			finalResultArray = append(finalResultArray, findingResult...)
			fileFoundCount++
		}
	}
	semiformat := fmt.Sprintf("%q\n", finalResultArray) // Turn the slice into a string that looks like ["one" "two" "three"]
	tokens := strings.Split(semiformat, "\" \"")	// Split this string by spaces
	finalResults := strings.Join(tokens, "\", \"")
	
	fmt.Println(fmt.Sprintf(`::set-output name=fileFoundCount::%d`, fileFoundCount))
	fmt.Println(fmt.Sprintf(`::set-output name=resultArray::%v`, finalResults))
}
