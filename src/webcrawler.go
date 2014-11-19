package main

import (
	"fmt"
)

type Fetcher interface {
	// Fetch returns the body of URL and
	// a slice of URLs found on that page.
	Fetch(url string) (body string, urls []string, err error)
}

type FetchResult struct {
	Url string
	Depth int
	Urls []string
	Body string
	Err error
}

// Crawl uses fetcher to recursively crawl
// pages starting with url, to a maximum of depth.
func Crawl(url string, depth int, fetcher Fetcher) {
	if depth <= 0 {
		return
	}

	ch := make(chan FetchResult)
	urlsInProgressOrFetched := make(map[string]bool)
	inProgress := 0

	startFetch := func(url string, depth int) {
		fmt.Printf("START %s depth %d\n", url, depth)
		go func(){
			body, urls, err := fetcher.Fetch(url)
			ch <- FetchResult{Url: url, Depth: depth, Urls: urls, Body: body, Err: err}
		}()
		urlsInProgressOrFetched[url] = true
		inProgress++
	}

	startFetch(url, depth)

	for inProgress > 0 {
		fmt.Println("inProgress", inProgress)
		fd := <- ch
		inProgress--

		if (fd.Err != nil) {
			fmt.Printf("Error fetching %s: %s\n", fd.Url,fd.Err)
		} else {
			fmt.Printf("found: %s %q\n", fd.Url, fd.Body)
			if fd.Depth > 0 {
				for _, url := range fd.Urls {
					if (!urlsInProgressOrFetched[url]) {
						startFetch(url, depth - 1)
					} else {
						fmt.Printf("Skip %s - already fetched\n", url)
					}
				}
			} else {
				fmt.Printf("Skip %s found urls %s - depth too low\n", fd.Url, fd.Urls)
			}
		}
	}
}

func main() {
	Crawl("http://golang.org/", 1, fetcher)
}

// fakeFetcher is Fetcher that returns canned results.
type fakeFetcher map[string]*fakeResult

type fakeResult struct {
	body string
	urls []string
}

func (f fakeFetcher) Fetch(url string) (string, []string, error) {
	if res, ok := f[url]; ok {
		return res.body, res.urls, nil
	}
	return "", nil, fmt.Errorf("not found: %s", url)
}

// fetcher is a populated fakeFetcher.
var fetcher = fakeFetcher{
	"http://golang.org/": &fakeResult{
		"The Go Programming Language",
		[]string{
			"http://golang.org/pkg/",
			"http://golang.org/cmd/",
		},
	},
	"http://golang.org/pkg/": &fakeResult{
		"Packages",
		[]string{
			"http://golang.org/",
			"http://golang.org/cmd/",
			"http://golang.org/pkg/fmt/",
			"http://golang.org/pkg/os/",
		},
	},
	"http://golang.org/pkg/fmt/": &fakeResult{
		"Package fmt",
		[]string{
			"http://golang.org/",
			"http://golang.org/pkg/",
		},
	},
	"http://golang.org/pkg/os/": &fakeResult{
		"Package os",
		[]string{
			"http://golang.org/",
			"http://golang.org/pkg/",
		},
	},
}
