// https://dave.cheney.net/paste/concurrency-made-easy.pdf
package main

import (
	"fmt"
	"strconv"
	"sync"
)

func restore(repos []string) error {
	errChan := make(chan error, 1)
	sem := make(chan int, 4) // four jobs at once
	var wg sync.WaitGroup
	wg.Add(len(repos))
	for _, repo := range repos {
		sem <- 1
		go func() {
			defer func() {
				wg.Done()
				<-sem
			}()
			if err := fetch(repo); err != nil {
				errChan <- err
			}
		}()
	}
	wg.Wait()
	close(sem)
	close(errChan)
	return <-errChan
}

// fake db fetch func
func fetch(repo string) error {
	return nil
}

func main() {
	n := 10
	repos := make([]string, n)
	for i := 0; i < n; i++ {
		repos[i] = strconv.Itoa(i)
	}

	err := restore(repos)
	if err != nil {
		fmt.Println(err.Error())
	}
}
