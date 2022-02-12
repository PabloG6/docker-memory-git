package main

import (
	log "github.com/sirupsen/logrus"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/memory"
)
func init() {
	customFormatter := new(log.TextFormatter)
	
    customFormatter.TimestampFormat = "2006-01-02 15:04:05"
	customFormatter.FullTimestamp = true
	log.SetFormatter(customFormatter)
}
func main() {
	storage := memory.NewStorage()
	repo, err := git.Clone(storage, nil, &git.CloneOptions{URL: "https://github.com/skunight/nestjs-redis.git"})
	log.Info("Cloning Repository")
	manageError(err)
	commits, err := repo.CommitObjects();

	manageError(err)

	log.Info("Searching Files")
	commits.ForEach(func(c *object.Commit) error {
		
		log.Info("Printing files for commit sha ", c.Hash.String())
		files, err := c.Files()
		manageError(err)

		files.ForEach(func(f *object.File) error {
			contents := f.Name
			log.Info("File: ", contents)
			return nil
		})
		return nil;
	})
}

func manageError(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}