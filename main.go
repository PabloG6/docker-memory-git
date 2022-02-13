package main

import (
	"archive/tar"
	"bufio"
	"bytes"
	"context"


	"io"

	"io/fs"
	"os"
	"path/filepath"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/memory"
)

const (
	RootPath = "/"
)

func init() {
	customFormatter := new(log.TextFormatter)

	customFormatter.TimestampFormat = "2006-01-02 15:04:05"
	customFormatter.FullTimestamp = true
	log.SetFormatter(customFormatter)
}
func main() {
	storage := memory.NewStorage()
	repo, err := git.Clone(storage, nil, &git.CloneOptions{URL: "https://github.com/PabloG6/ng-config-deploy.git", Depth: 1})
	log.Info("Cloning Repository")
	manageError(err)
	commits, err := repo.CommitObjects()

	aferofs := afero.NewMemMapFs()
	manageError(err)

	aferofs.Mkdir(RootPath, os.ModePerm)
	log.Info("Searching Files")
	commits.ForEach(func(c *object.Commit) error {

		log.Info("Printing files for commit sha ", c.Hash.String())
		files, err := c.Files()
		manageError(err)

		files.ForEach(func(f *object.File) error {
			memfile, err := aferofs.Create(filepath.Join(RootPath, f.Name))
			manageError(err)
			
			fileReader, err := f.Reader()
		
			manageError(err)

			io.Copy(memfile, fileReader)

			return nil
		})

		return nil
	})

	log.Info("checking in memory file path stat")

	manageError(err)
	var writer bytes.Buffer

	tarWriter := tar.NewWriter(&writer)
	defer tarWriter.Close()
	afero.Walk(aferofs, RootPath, func(path string, info fs.FileInfo, err error) error {
		manageError(err)
		log.Info("in memory file path: ", path)

		manageError(err)
		tarHeader, err := tar.FileInfoHeader(info, path)
		manageError(err)
		tarHeader.Name = path
		tarWriter.WriteHeader(tarHeader)

		contents, err := aferofs.Open(path)
		if err != nil {
			manageError(err)
		}

	
		log.Info("writing file contents to tar writer for file: ", path)
		io.Copy(tarWriter, contents)
		return nil
	})

	tarIO := bytes.NewReader(writer.Bytes())
	client, err := client.NewClientWithOpts()
	manageError(err)

	var imageBuildResponse types.ImageBuildResponse;
	imageBuildResponse, err = client.ImageBuild(context.Background(), tarIO, types.ImageBuildOptions{Tags: []string{"hello-world"}})
	manageError(err)


	scanner := bufio.NewScanner(imageBuildResponse.Body)
	for scanner.Scan() {
		text := scanner.Text()
		log.Info(text)
	}
	

}

func manageError(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}
