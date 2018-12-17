package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/asdine/storm"
	"github.com/watuwat/imdb"
)

func download(url string) (io.ReadCloser, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}

type info struct {
	link  string
	value imdb.Assigner
	index int
}

type display struct {
	name  string
	index int
	n     int
}

func main() {
	db, err := storm.Open("my.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	sources := []struct {
		link    string
		builder func() imdb.Assigner
	}{
		{
			link: "https://datasets.imdbws.com/name.basics.tsv.gz",
			builder: func() imdb.Assigner {
				return &imdb.Name{}
			},
		},
		{
			link: "https://datasets.imdbws.com/title.akas.tsv.gz",
			builder: func() imdb.Assigner {
				return &imdb.Akas{}
			},
		},
		{
			link: "https://datasets.imdbws.com/title.basics.tsv.gz",
			builder: func() imdb.Assigner {
				return &imdb.Basic{}
			},
		},
		{
			link: "https://datasets.imdbws.com/title.crew.tsv.gz",
			builder: func() imdb.Assigner {
				return &imdb.Crew{}
			},
		},
		{
			link: "https://datasets.imdbws.com/title.episode.tsv.gz",
			builder: func() imdb.Assigner {
				return &imdb.Episode{}
			},
		},
		{
			link: "https://datasets.imdbws.com/title.principals.tsv.gz",
			builder: func() imdb.Assigner {
				return &imdb.Principals{}
			},
		},
		{
			link: "https://datasets.imdbws.com/title.ratings.tsv.gz",
			builder: func() imdb.Assigner {
				return &imdb.Rating{}
			},
		},
	}

	records := make(chan *info, 1)

	gw := sync.WaitGroup{}

	for index, source := range sources {
		gw.Add(1)
		go func(link string, builder func() imdb.Assigner, index int) {
			defer gw.Done()

			input, err := download(link)
			if err != nil {
				panic(err)
			}

			src, err := gzip.NewReader(input)
			if err != nil {
				panic(err)
			}

			values := imdb.Parser(src, builder)
			for value := range values {
				records <- &info{
					link:  link,
					value: value,
					index: index,
				}
			}

		}(source.link, source.builder, index)
	}

	go func() {
		gw.Wait()
		close(records)
	}()

	counts := make(map[int]*display)

	// clean screen
	fmt.Printf("\033[2J")

	for record := range records {
		if _, ok := counts[record.index]; !ok {
			counts[record.index] = &display{
				name: record.link,
				n:    0,
			}
		}

		err = db.Save(record.value)
		if err != nil {
			panic(err)
		}

		count := counts[record.index]
		count.n++

		total := 0
		i := 1
		for key, value := range counts {
			fmt.Printf("\033[%d;%dH%s:   \t%d", key+1, 1, value.name, value.n)
			total += value.n
			i++
		}
		fmt.Printf("\033[%d;%dH%s: %d", i, 1, "total: ", total)
	}
}
