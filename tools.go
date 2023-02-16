package main

import (
	"bufio"
	"io"
	"os"
	"regexp"
	"strings"
)

func checkFileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func getFile(filename string) (string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer f.Close()

	b, err := io.ReadAll(f)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func getTFTypes(filename string) string {
	reTFTypes := regexp.MustCompile(`^\S+_(datasource|resource).go$`)
	tfTypes := ""

	if reTFTypes.MatchString(filename) {
		tfTypes = reTFTypes.FindStringSubmatch(filename)[1]
	}
	return tfTypes
}

func getPackageName(str string) string {
	rePackage := regexp.MustCompile(`^package\s+([a-zA-Z0-9_]+)$`)

	packageName := ""

	scanner := bufio.NewScanner(strings.NewReader(str))
	for scanner.Scan() {
		line := scanner.Text()
		if rePackage.MatchString(line) {
			packageName = rePackage.FindStringSubmatch(line)[1]
		}
	}
	return packageName
}

func getTFName(str string) string {
	reTFName := regexp.MustCompile(`^\/\/(?:\s+)?tfname:(?:\s+)?(\S+)$`)

	tfName := ""

	scanner := bufio.NewScanner(strings.NewReader(str))
	for scanner.Scan() {
		line := scanner.Text()
		if reTFName.MatchString(line) {
			tfName = reTFName.FindStringSubmatch(line)[1]
		}
	}
	return tfName
}
