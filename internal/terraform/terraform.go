package terraform

import (
	"bufio"
	"regexp"
	"strings"
)

// GetTFTypes returns the terraform type (datasource or resource) from the filename.
func GetTFTypes(filename string) string {
	reTFTypes := regexp.MustCompile(`^\S+_(datasource|resource).go$`)
	tfTypes := ""

	if reTFTypes.MatchString(filename) {
		tfTypes = reTFTypes.FindStringSubmatch(filename)[1]
	}
	return tfTypes
}

// GetPackageName returns the golang package name from the file content.
func GetPackageName(str string) string {
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

// GetTFName returns the terraform name from the file content looking for comment //tfname: my_tfname.
func GetTFName(str string) (categoryName, resourceName string) {
	reTFName := regexp.MustCompile(`^\/\/(?:\s+)?tfname:\s+([a-z0-9]+)_?(.*)`)

	scanner := bufio.NewScanner(strings.NewReader(str))
	for scanner.Scan() {
		line := scanner.Text()
		if reTFName.MatchString(line) {
			return reTFName.FindStringSubmatch(line)[1], reTFName.FindStringSubmatch(line)[2]
		}
	}

	return "", ""
}
