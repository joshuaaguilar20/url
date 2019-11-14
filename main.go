package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"gopkg.in/yaml.v2"
)

//	"gopkg.in/yaml.v2"
type timeHandler struct {
	format string
}
type yamlParse struct {
	ShortURL string `yaml:"shortUrl"`
	LongURL  string `yaml:"longUrl"`
}

//Returns Current Time in Requested Format
func (th *timeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	tm := time.Now().Format(th.format)
	w.Write([]byte("The time is: " + tm))
}

//hello: Sends Response Hello World
func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}

// Elsewhere in our code after we've discovered an error.

//MapHandler: Maps Short Urls to Long and if current URL matches the short Redirects to the long Value
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//Gets Path From Request
		path := r.URL.Path
		if value, ok := pathsToUrls[path]; ok {
			http.Redirect(w, r, value, http.StatusFound)
			return
		}
	}
}

//parses Incomming Flag and Yaml File From CommandLine Argument*
func parseFlag() yamlParse {
	yamlPtr := flag.String("file", "foo", "yaml file withshort cuts to be parsed")
	flag.Parse()
	if *yamlPtr != "foo" {
		yamlFile, err := ioutil.ReadFile(*yamlPtr)
		if err != nil {
			fmt.Println("Error in Yaml File: AT stackTrace: 12222")
		}
		fmt.Println(string(yamlFile))
		var yamlConfig yamlParse
		err = yaml.Unmarshal(yamlFile, &yamlConfig)
		if err != nil {
			fmt.Printf("Error parsing YAML file: %s\n", err)
		}

		return yamlConfig
	}
	var x yamlParse
	return x
}

//StartServer Takes Map of Short Urls to Map to Long Urls. Starts Mux Server with defualt routes and Port
func StartServer(r map[string]string) {
	mux := http.NewServeMux()
	th := &timeHandler{format: time.RFC1123}
	mux.Handle("/time", th)
	mux.HandleFunc("/", MapHandler(r, mux))
	mux.HandleFunc("/hello", hello)
	log.Println("Listening...")
	http.ListenAndServe(":3000", mux)
}
func main() {
	mapUrl := make(map[string]string)
	flag := parseFlag()
	mapUrl[flag.ShortURL] = "https://www.reddit.com"
	fmt.Println(flag)
	StartServer(mapUrl)

}
