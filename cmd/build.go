//// +build dev

// build.go automates proper versioning of usersms binaries
// and installer scripts.
// Use it like:   go run build.go
// The result binary will be located in bin/app
// You can customize the build with the -goos, -goarch, and
// -goarm CLI options:   go run build.go -goos=windows
//
// This program is NOT required to build usersms from source
// since it is go-gettable. (You can run plain `go build`
// from respective cmd sub-directories to get a binary).
package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"

	"github.com/tomogoma/go-typed-errors"
	"github.com/tomogoma/usersms/pkg/config"
	"github.com/tomogoma/usersms/pkg/fileutils"
)

func main() {
	var goos, goarch, goarm string
	var help bool
	flag.StringVar(&goos, "goos", "",
		"GOOS\tThe operating system for which to compile\n"+
			"\t\tExamples are linux, darwin, windows, netbsd.")
	flag.StringVar(&goarch, "goarch", "",
		"GOARCH\tThe architecture, or processor, for which to compile code.\n"+
			"\t\tExamples are amd64, 386, arm, ppc64.")
	flag.StringVar(&goarm, "goarm", "",
		"GOARM\tFor GOARCH=arm, the ARM architecture for which to compile.\n"+
			"\t\tValid values are 5, 6, 7.")
	flag.BoolVar(&help, "help", false, "Show this help message")
	flag.Parse()
	if help {
		flag.Usage()
		os.Exit(0)
	}
	if err := buildMicroservice(goos, goarch, goarm); err != nil {
		log.Fatalf("buildMicroservice error: %v", err)
	}
	if err := installVars(); err != nil {
		log.Fatalf("write installer script error: %v", err)
	}
	if err := buildGcloud(); err != nil {
		log.Fatalf("build GCloud error: %v", err)
	}
}

func installVars() error {
	content := `#!/usr/bin/env bash
NAME="` + config.Name + `"
VERSION="` + config.VersionFull + `"
DESCRIPTION="` + config.Description + `"
CANONICAL_NAME="` + config.CanonicalName() + `"
CONF_DIR="` + config.DefaultConfDir() + `"
CONF_FILE="` + config.DefaultConfPath() + `"
INSTALL_DIR="` + config.DefaultInstallDir() + `"
INSTALL_FILE="` + config.DefaultInstallPath() + `"
UNIT_NAME="` + config.DefaultSysDUnitName() + `"
UNIT_FILE="` + config.DefaultSysDUnitFilePath() + `"
DOCS_DIR="` + config.DefaultDocsDir() + `"
`
	return ioutil.WriteFile("install/vars.sh", []byte(content), 0755)
}

func buildMicroservice(goos, goarch, goarm string) error {
	docsDir := path.Join("install", "docs", config.VersionMajorPrefixed(), config.Name, "docs")
	if err := compileDocs(docsDir); err != nil {
		return err
	}
	args := []string{"build", "-o", "bin/app", "./cmd/micro"}
	cmd := exec.Command("go", args...)
	cmd.Env = os.Environ()
	for _, env := range []string{
		"GOOS=" + goos,
		"GOARCH=" + goarch,
		"GOARM=" + goarm,
	} {
		cmd.Env = append(cmd.Env, env)
	}
	if out, err := cmd.CombinedOutput(); err != nil {
		return errors.Newf("build: %s - %v", out, err)
	}
	return nil
}

func buildGcloud() error {
	confDir := config.DefaultConfDir("cmd", "gcloud", "conf")

	if err := os.MkdirAll(confDir, 0755); err != nil {
		return errors.Newf("create conf dir: %v", err)
	}

	docsDir := path.Join(config.DefaultDocsDir(), config.VersionMajorPrefixed(), config.Name, "docs")
	if err := compileDocs(docsDir); err != nil {
		return err
	}

	err := fileutils.CopyIfDestNotExists(path.Join("install", "conf.yml"), config.DefaultConfPath())
	if err != nil {
		return err
	}
	if err := cleanGCloudConfFile(); err != nil {
		return errors.Newf("clean gcloud config file: %v", err)
	}

	return nil
}

func compileDocs(docsDir string) error {

	subjDir := path.Join("pkg", "handler", "http")
	headerFile := path.Join(subjDir, "apidoc_header.md")
	APIDocConfFile := path.Join(subjDir, "apidoc.json")

	apiDoc := struct {
		Name        string      `json:"name"`
		Version     string      `json:"version"`
		Description string      `json:"description"`
		Title       string      `json:"title"`
		Header      interface{} `json:"header"`
	}{
		Name:        config.Name,
		Version:     config.VersionFull,
		Description: config.Description,
		Title:       config.CanonicalName(),
		Header: struct {
			Title    string `json:"title"`
			FileName string `json:"filename"`
		}{
			Title:    "Introduction",
			FileName: headerFile,
		},
	}

	apiDocB, err := json.Marshal(apiDoc)
	if err != nil {
		return errors.Newf("Marshal API doc config: %v", err)
	}

	err = ioutil.WriteFile(APIDocConfFile, apiDocB, 0655)
	if err != nil {
		return errors.Newf("Write API doc file: %v", err)
	}

	args := []string{"-i", subjDir, "-c", subjDir, "-o", docsDir}
	cmd := exec.Command("apidoc", args...)
	if out, err := cmd.CombinedOutput(); err != nil {
		return errors.Newf("generate http docs: %v: %s", err, out)
	}
	return nil
}

func cleanGCloudConfFile() error {
	newPath := config.DefaultConfPath()
	confContent, err := ioutil.ReadFile(config.DefaultConfPath())
	if err != nil {
		return errors.Newf("read file for transform: %v", err)
	}
	confContentClean := bytes.Replace(confContent, []byte(config.SysDConfDir()+"/"), []byte("conf/"), -1)
	err = ioutil.WriteFile(newPath, confContentClean, 0644)
	if err != nil {
		return errors.Newf("write transformed file: %v", err)
	}
	return nil
}
