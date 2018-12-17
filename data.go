package imdb

import (
	"encoding/json"
	"strconv"
	"strings"
)

type Assigner interface {
	Assign(key string, value string)
}

func toInt(value string) int {
	n, _ := strconv.ParseInt(value, 10, 32)
	return int(n)
}

func toFloat(value string) float32 {
	f, _ := strconv.ParseFloat(value, 32)
	return float32(f)
}

func toStringArray(value string) []string {
	if strings.Index(value, "[\"") == 0 {
		arr := make([]string, 0)
		json.Unmarshal([]byte(value), &arr)
		return arr
	}

	if value == "" {
		return make([]string, 0)
	}

	return strings.Split(value, ",")
}

func trimValue(value string) string {
	if value == "\\N" {
		return ""
	}
	return value
}

type Akas struct {
	ID              string   `json:"id" storm:"id"`            //a tconst, an alphanumeric unique identifier of the title
	Ordering        int      `json:"ordering"`                 //a number to uniquely identify rows for a given titleId
	Title           string   `json:"title"`                    //the localized title
	Region          string   `json:"region"`                   //the region for this version of the title
	Language        string   `json:"language"`                 //the language of the title
	Types           []string `json:"types" storm:"index"`      //Enumerated set of attributes for this alternative title. One or more of the following: "alternative", "dvd", "festival", "tv", "video", "working", "original", "imdbDisplay". New values may be added in the future without warning
	Attributes      []string `json:"attributes" storm:"index"` //Additional terms to describe this alternative title, not enumerated
	IsOriginalTitle bool     `json:"isOriginalTitle"`          //0: not original title; 1: original title
}

func (a *Akas) Assign(key string, value string) {
	value = trimValue(value)

	switch key {
	case "titleId":
		a.ID = value
	case "ordering":
		a.Ordering = toInt(value)
	case "title":
		a.Title = value
	case "region":
		a.Region = value
	case "language":
		a.Language = value
	case "types":
		a.Types = toStringArray(value)
	case "attributes":
		a.Attributes = toStringArray(value)
	case "isOriginalTitle":
		a.IsOriginalTitle = value == "1"
	}
}

// title.basics.tsv.gz - Contains the following information for titles:
type Basic struct {
	ID             string   `json:"id" storm:"id"`                // alphanumeric unique identifier of the title
	TitleType      string   `json:"titleType" storm:"index"`      // the type/format of the title (e.g. movie, short, tvseries, tvepisode, video, etc)
	PrimaryTitle   string   `json:"primaryTitle"`                 // the more popular title / the title used by the filmmakers on promotional materials at the point of release
	OriginalTitle  string   `json:"originalTitle"`                // original title, in the original language
	IsAdult        bool     `json:"isAdult" storm:"index"`        // 0: non-adult title; 1: adult title
	StartYear      string   `json:"startYear" storm:"index"`      // represents the release year of a title. In the case of TV Series, it is the series start year
	EndYear        string   `json:"endYear" storm:"index"`        // TV Series end year. ‘\N’ for all other title types
	RuntimeMinutes string   `json:"runtimeMinutes" storm:"index"` // primary runtime of the title, in minutes
	Genres         []string `json:"genres" storm:"index"`         // includes up to three genres associated with the title
}

func (b *Basic) Assign(key string, value string) {
	value = trimValue(value)

	switch key {
	case "tconst":
		b.ID = value
	case "titleType":
		b.TitleType = value
	case "primaryTitle":
		b.PrimaryTitle = value
	case "originalTitle":
		b.OriginalTitle = value
	case "isAdult":
		b.IsAdult = value == "1"
	case "startYear":
		b.StartYear = value
	case "endYear":
		b.EndYear = value
	case "runtimeMinutes":
		b.RuntimeMinutes = trimValue(value)
	case "genres":
		b.Genres = toStringArray(value)
	}
}

// title.crew.tsv.gz – Contains the director and writer information for all the titles in IMDb. Fields include:
type Crew struct {
	ID        string   `json:"id" storm:"id"`           // alphanumeric unique identifier of the title
	Directors []string `json:"directors" storm:"index"` // director(s) of the given title
	Writers   []string `json:"writers" storm:"index"`   // writer(s) of the given title
}

func (c *Crew) Assign(key string, value string) {
	value = trimValue(value)

	switch key {
	case "tconst":
		c.ID = value
	case "directors":
		c.Directors = toStringArray(value)
	case "writers":
		c.Writers = toStringArray(value)
	}
}

// title.episode.tsv.gz – Contains the tv episode information. Fields include:
type Episode struct {
	ID      string `json:"id" storm:"id"`         // alphanumeric identifier of episode
	TitleID string `json:"titleId" storm:"index"` // alphanumeric identifier of the parent TV Series
	Season  int    `json:"season" storm:"index"`  // season number the episode belongs to
	Episode int    `json:"episode" storm:"index"` // episode number of the tconst in the TV series
}

func (e *Episode) Assign(key string, value string) {
	value = trimValue(value)

	switch key {
	case "tconst":
		e.ID = value
	case "parentTconst":
		e.TitleID = value
	case "seasonNumber":
		e.Season = toInt(value)
	case "episodeNumber":
		e.Episode = toInt(value)
	}
}

// title.principals.tsv.gz – Contains the principal cast/crew for titles
type Principals struct {
	ID         string   `json:"id" storm:"id"`            // alphanumeric unique identifier of the title
	Ordering   int      `json:"ordering"`                 // a number to uniquely identify rows for a given titleId
	NameID     string   `json:"nameId" storm:"index"`     // alphanumeric unique identifier of the name/person
	Category   string   `json:"category" storm:"index"`   // the category of job that person was in
	Job        []string `json:"job" storm:"index"`        // the specific job title if applicable, else '\N'
	Characters []string `json:"characters" storm:"index"` // the name of the character played if applicable, else '\N'
}

func (p *Principals) Assign(key string, value string) {
	value = trimValue(value)

	switch key {
	case "tconst":
		p.ID = value
	case "ordering":
		p.Ordering = toInt(value)
	case "nconst":
		p.NameID = value
	case "category":
		p.Category = value
	case "job":
		p.Job = toStringArray(value)
	case "characters":
		p.Characters = toStringArray(value)
	}
}

// title.ratings.tsv.gz – Contains the IMDb rating and votes information for titles
type Rating struct {
	ID            string  `json:"id" storm:"id"`               // alphanumeric unique identifier of the title
	AverageRating float32 `json:"averageRating" storm:"index"` // weighted average of all the individual user ratings
	Votes         int     `json:"votes" storm:"index"`         // number of votes the title has received
}

func (r *Rating) Assign(key string, value string) {
	value = trimValue(value)

	switch key {
	case "tconst":
		r.ID = value
	case "averageRating":
		r.AverageRating = toFloat(value)
	case "numVotes":
		r.Votes = toInt(value)
	}
}

// name.basics.tsv.gz – Contains the following information for names:
type Name struct {
	ID                string   `json:"id" storm:"id"`                   // alphanumeric unique identifier of the name/person
	PrimaryName       string   `json:"primaryName" storm:"index"`       // name by which the person is most often credited
	BirthYear         string   `json:"birthYear" storm:"index"`         // in YYYY format
	DeathYear         string   `json:"deathYear" storm:"index"`         // in YYYY format if applicable, else '\N'
	PrimaryProfession []string `json:"primaryProfession" storm:"index"` // the top-3 professions of the person
	KnownForTitles    []string `json:"knownForTitles" storm:"index"`    // titles the person is known for
}

func (n *Name) Assign(key string, value string) {
	value = trimValue(value)

	switch key {
	case "nconst":
		n.ID = value
	case "primaryName":
		n.PrimaryName = value
	case "birthYear":
		n.BirthYear = value
	case "deathYear":
		n.DeathYear = value
	case "primaryProfession":
		n.PrimaryProfession = toStringArray(value)
	case "knownForTitles":
		n.KnownForTitles = toStringArray(value)
	}
}
