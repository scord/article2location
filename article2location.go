package main

import (
	"code.google.com/p/go.net/html"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

type country struct {
	Name         string
	Cca2         string
	Capital      string
	AltSpellings []string
	Demonym      string
	Region       string
	Subregion    string
}

func (c *country) MakeLowercase() *country {

	nc := &country{}

	nc.Name = strings.ToLower(c.Name)
	nc.Cca2 = strings.ToLower(c.Cca2)
	nc.Capital = strings.ToLower(c.Capital)
	nc.Demonym = strings.ToLower(c.Demonym)
	nc.Region = strings.ToLower(c.Region)
	nc.Subregion = strings.ToLower(c.Subregion)
	nc.AltSpellings = c.AltSpellings

	for i := range c.AltSpellings {
		nc.AltSpellings[i] = strings.ToLower(c.AltSpellings[i])
	}

	return nc
}
func main() {
	fmt.Println(countrySubject(os.Args[1]))
}

func getCountries() ([]country, error) {
	file, err := ioutil.ReadFile("countries.json")
	if err != nil {
		fmt.Println(err)
	}
	var countries []country
	err = json.Unmarshal(file, &countries)
	if err != nil {
		return nil, err
	}

	return countries, nil
}

func countrySubject(url string) string {

	countries, err := getCountries()

	a, err := getarticle(url)
	if err != nil {
		fmt.Println(err)
	}

	for _, c := range countries {

		a = strings.Replace(a, c.Name, strings.Replace(c.Name, " ", "", -1), -1)
		a = strings.Replace(a, c.Capital, strings.Replace(c.Capital, " ", "", -1), -1)
		a = strings.Replace(a, c.Demonym, strings.Replace(c.Demonym, " ", "", -1), -1)
		a = strings.Replace(a, c.Region, strings.Replace(c.Region, " ", "", -1), -1)
		a = strings.Replace(a, c.Subregion, strings.Replace(c.Subregion, " ", "", -1), -1)

		for _, s := range c.AltSpellings {
			a = strings.Replace(a, s, strings.Replace(s, " ", "", -1), -1)
		}

	}

	a = strings.ToLower(a)

	words := wordCount(a)

	topscore := 0
	topcountry := ""

	for _, c := range countries {
		lc := c.MakeLowercase()
		score := 0
		score += 3 * words[strings.Replace(lc.Name, " ", "", -1)]
		score += 2 * words[strings.Replace(lc.Capital, " ", "", -1)]
		score += 3 * words[strings.Replace(lc.Demonym, " ", "", -1)]
		score += words[strings.Replace(lc.Region, " ", "", -1)]
		score += words[strings.Replace(lc.Subregion, " ", "", -1)]

		for _, spelling := range lc.AltSpellings {
			if len(spelling) >= 2 &&
				spelling != "in" &&
				spelling != "is" &&
				spelling != "as" &&
				spelling != "to" &&
				spelling != "at" &&
				spelling != "be" &&
				spelling != "of" &&
				spelling != "it" &&
				spelling != "mr" &&
				spelling != "ms" &&
				spelling != "by" &&
				spelling != "me" &&
				spelling != "my" &&
				spelling != "no" &&
				spelling != "mp" &&
				spelling != "la" &&
				spelling != "on" {
				score += 2 * words[strings.ToLower(strings.Replace(spelling, " ", "", -1))]
			}
		}

		if score > topscore {
			topcountry = c.Name
			topscore = score

		}

	}
	return topcountry

}

func wordCount(s string) map[string]int {
	s = strings.Replace(s, ",", " ", -1)
	s = strings.Replace(s, ".", " ", -1)
	s = strings.Replace(s, "[", " ", -1)
	s = strings.Replace(s, "]", " ", -1)
	s = strings.Replace(s, "\"", " ", -1)
	s = strings.Replace(s, "'", " ", -1)
	s = strings.Replace(s, "/", " ", -1)
	s = strings.Replace(s, "-", " ", -1)

	words := strings.Fields(s)

	c := make(map[string]int)

	for _, word := range words {
		c[word] += 1
	}

	return c
}

func getarticle(url string) (string, error) {

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	var b []byte

	b, err = ioutil.ReadAll(resp.Body)
	s := string(b)

	doc, err := html.Parse(strings.NewReader(s))
	if err != nil {
		log.Fatal(err)
	}

	a := url

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.TextNode && n.Parent.Data == "h1" {
			a += n.Data
		}
		if n.Type == html.TextNode && n.Parent.Data == "p" {
			a += n.Data
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	return a, nil
}
