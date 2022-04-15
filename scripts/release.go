package main

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
)

const name = "xmlymft"
const serverbinRelativePath = "components/xmlymft/bins"
const appbinRelativePath = "fyne-cross/bin"
const appDistRelativePath = "fyne-cross/dist"
const releasePath = "release"

// Get main app file data.
func getAppFile(platform string, arch string) []byte {
	ext := ""
	if platform == "windows" {
		ext = ".exe"
	} else if platform == "linux" {
	} else {
		log.Fatalf("unknow platform %s", platform)
		os.Exit(1)
	}
	id := fmt.Sprintf("%s-%s", platform, arch)
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalln(err.Error())
		os.Exit(1)
	}

	file, err := os.Open(
		filepath.Clean(filepath.Join(wd, appbinRelativePath, id, name+ext)),
	)
	if err != nil {
		log.Fatalln(err.Error())
		os.Exit(1)
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatalln(err.Error())
		os.Exit(1)
	}
	return data
}

// Get the server app file data.
func getServerFile(platform string, arch string) []byte {
	ext := ""
	if platform == "windows" {
		platform = "win"
		ext = ".exe"
	} else if platform == "linux" {
		platform = "linux"
	} else {
		log.Fatalf("unknow platform %s", platform)
		os.Exit(1)
	}
	if arch == "amd64" {
		arch = "64"
	} else if arch == "386" {
		arch = "32"
	} else {
		log.Fatalf("unknow arch %s", arch)
		os.Exit(1)
	}
	filename := fmt.Sprintf("xmlymft_%s%s%s", platform, arch, ext)
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalln(err.Error())
		os.Exit(1)
	}

	file, err := os.Open(
		filepath.Clean(filepath.Join(wd, serverbinRelativePath, filename)),
	)
	if err != nil {
		log.Fatalln(err.Error())
		os.Exit(1)
	}
	defer file.Close()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatalln(err.Error())
		os.Exit(1)
	}
	return data
}

func packWindows(arch string) {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalln(err.Error())
		os.Exit(1)
	}

	// Create the archive file.
	archivefile, err := os.Create(filepath.Join(wd, releasePath, name+"-windows"+"-"+arch+".zip"))
	if err != nil {
		log.Fatalln(err.Error())
		os.Exit(1)
	}
	defer archivefile.Close()

	// Create a new zip archive.
	zipWriter := zip.NewWriter(archivefile)
	defer zipWriter.Close()

	// Archive the app file.
	appfile, err := zipWriter.Create(path.Join(name, name+".exe"))
	if err != nil {
		log.Fatalln(err.Error())
		os.Exit(1)
	}
	_, err = appfile.Write(getAppFile("windows", arch))
	if err != nil {
		log.Fatalln(err.Error())
		os.Exit(1)
	}

	// Archive the server file.
	serverfile, err := zipWriter.Create(path.Join(name, "server.exe"))
	if err != nil {
		log.Fatalln(err.Error())
		os.Exit(1)
	}
	_, err = serverfile.Write(getServerFile("windows", arch))
	if err != nil {
		log.Fatalln(err.Error())
		os.Exit(1)
	}
}

func packLinux(arch string) {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalln(err.Error())
		os.Exit(1)
	}
	// Create the archive file.
	archivefile, err := os.Create(filepath.Join(wd, releasePath, name+"-linux"+"-"+arch+".tar.gz"))
	if err != nil {
		log.Fatalln(err.Error())
		os.Exit(1)
	}
	defer archivefile.Close()

	gzipWriter, err := gzip.NewWriterLevel(archivefile, gzip.DefaultCompression)
	if err != nil {
		log.Fatalln(err.Error())
		os.Exit(1)
	}
	defer gzipWriter.Close()

	tarWriter := tar.NewWriter(gzipWriter)
	defer tarWriter.Close()

	// Archive the app file.
	appfile := getAppFile("linux", arch)
	appfileHeader := tar.Header{
		Name: path.Join(name, name),
		Mode: 0644,
		Size: int64(len(appfile)),
	}
	err = tarWriter.WriteHeader(&appfileHeader)
	if err != nil {
		log.Fatalln(err.Error())
		os.Exit(1)
	}
	_, err = tarWriter.Write(appfile)
	if err != nil {
		log.Fatalln(err.Error())
		os.Exit(1)
	}

	// Archive the server file.
	serverfile := getServerFile("linux", arch)
	serverfileHeader := tar.Header{
		Name: path.Join(name, "server"),
		Mode: 0644,
		Size: int64(len(serverfile)),
	}
	err = tarWriter.WriteHeader(&serverfileHeader)
	if err != nil {
		log.Fatalln(err.Error())
		os.Exit(1)
	}
	_, err = tarWriter.Write(serverfile)
	if err != nil {
		log.Fatalln(err.Error())
		os.Exit(1)
	}
}

func main() {
	os.Mkdir(releasePath, 0644)

	// windows
	packWindows("386")
	packWindows("amd64")

	// linux
	packLinux("386")
	packLinux("amd64")
}
