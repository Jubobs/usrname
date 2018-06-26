// TODO:
// * way to enumerate sites? => provide a map of checkers?
// * use URL rather than string concatenation
// * actually flesh out IsValid methods

// func(username string) URL
// func(url URL, func(Response) error): error

package main

import (
	"fmt"
	"github.com/jubobs/username-checker/sites"
	"os"
)

func main() {
	username := os.Args[1]
	checkers := sites.All()
	client := sites.NewClient()

	// see https://www.safaribooksonline.com/library/view/ultimate-go-programming/9780134757476/ugpg_04_10_03_00.html
	// and https://gist.github.com/Jubobs/3987c6f9f902489356ccd12422c64e1d
	n := len(checkers)
	ch := make(chan string, n)
	waitChecks := n

	for _, checker := range checkers {
		go func(c sites.NameChecker) {
			if err := c.Check(client, username); err != nil {
				ch <- err.Error()
			} else {
				ch <- fmt.Sprintf("%s is available and valid on %s", username, c.Name())
			}

		}(checker)
	}

	for waitChecks > 0 {
		r := <-ch
		fmt.Println(r)
		waitChecks--
	}
}
