package main

import (
	"github.com/gyulaieric/stauction/internal/repository"
)

func main() {
	gs := repository.NewGitStore("data/stalzone-database")
	gs.SyncRepo()
	// gs.StartPullWorker(context.Background(), time.Minute)
}
