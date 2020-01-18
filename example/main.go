package main

import (
	"fmt"
	"net/http"

	trending "github.com/ryomak/go-trending"
)

func main() {
	client := trending.NewClient(trending.WithHttpClient(http.DefaultClient))
	repo, _ := client.GetRepository("today", "go")
	fmt.Println(repo)
}
