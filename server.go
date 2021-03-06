package main

import (
	"runtime"
	"sync"
	"html/template"
	"net/http"
	"log"
	"encoding/json"
	"golang.org/x/oauth2"
	"github.com/google/go-github/github"
)

type SearchArgs struct {
	Repositories string
}

type Language struct {
	Name string
	Lines int
}

type Repository struct {
	ID int
	Url string
	Name string
	Owner string
	Languages []Language
}

var tmpl = template.Must(template.ParseFiles("src/github.com/utay/ScalingoTest/index.html",
	"src/github.com/utay/ScalingoTest/search.html"))

// getNewRepositories returns the last X repositories since the ID one.
func getNewRepositories(client github.Client, ID int) []github.Repository {
	optall := &github.RepositoryListAllOptions{Since: ID}
	repos,  _, err := client.Repositories.ListAll(optall)
	if err != nil {
		log.Fatal("Repositories:", err)
	}
	return repos
}

// getLastestID returns the ID of the last repository created.
// It uses the github events list.
func getLastestID(client github.Client) int {
	ID := 0
	for ; ID == 0 ; {
		opt := &github.ListOptions{PerPage: 100}
		events, _, _ := client.Activity.ListEvents(opt)
		for _, event := range events {
			if *event.Type == "CreateEvent" {
				m, _ := event.Payload().(map[string]interface{})
				if m["ref_type"] == "repository" {
					ID = *event.Repo.ID
					break
				}
			}
		}
	}
	return ID
}

// index executes the template of the root page "index.html".
func index(w http.ResponseWriter, r *http.Request) {
	err := tmpl.ExecuteTemplate(w, "index.html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// searchEngine returns the github repositories which are matching with a query.
func searchEngine(query string, repos []github.Repository) []github.Repository {
	results := make([]github.Repository, 0)
	for _, repo := range repos {
		name := *repo.FullName
		i := 0
		k := 0
		for ; i <= len(name) - len(query); i++ {
			j := 0
			for ; j < len(query); j++{
				if name[i + j] != query[j] {
					break
				}
			}
			k = j
			if j == len(query) {
				break
			}
		}
		if i <= len(name) - len(query) && k == len(query) {
			results = append(results, repo)
		} else {
			if len(name) == len(query) {
				for i = 0; i < len(name); i++ {
					if name[i] != query[i] {
						break
					}
				}
				if i == len(name) {
					results = append(results, repo)
				}
			}
		}
	}
	return results
}

// extractData fills the struct Repository with ID, Name, Url, Owner.
// It calls also the ListLanguages function of the github API to fill the Languages field.
// This function is called as a goroutine in the function search.
func extractData(client github.Client, repo github.Repository,
	results []Repository, i int, wg *sync.WaitGroup) {
	results[i].ID = *repo.ID
	results[i].Url = *repo.HTMLURL
	results[i].Name = *repo.FullName
	results[i].Owner = *repo.Owner.Login
	languages, _, _ := client.Repositories.ListLanguages(*repo.Owner.Login, *repo.Name)
	results[i].Languages = make([]Language, len(languages))
	j := 0
	for language, lines := range languages {
		results[i].Languages[j] = Language{language, lines}
		j++
	}
	wg.Done()
}

// search fills an array of the 100 last github repositories,
// filters it whether there is a query or not,
// starts one goroutine per Repository to extract the data,
// and converts to json every structs Repository so that JavaScript can use it.
// This function finally executes the template of the page "search.html".
func search(w http.ResponseWriter, r *http.Request) {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: "90c7a3fbf7e2a3c27c29ddb4600db0db67478ac0"},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	client := github.NewClient(tc)
	min := 100
	ID := getLastestID(*client)
	repos := make([]github.Repository, 0)
	for ; len(repos) < 100 ; {
		repos = append(repos, getNewRepositories(*client, ID - min)...)
		min += 100
	}
	repos = repos[0:100]
	q := r.URL.Query().Get("q")
	if len(q) != 0 {
		repos = searchEngine(q, repos)
	}
	results := make([]Repository, len(repos))
	var wg sync.WaitGroup
	// Start goroutines
	for i, repo := range repos {
		wg.Add(1)
		go extractData(*client, repo, results, i, &wg)
	}
	// Wait goroutines
	wg.Wait()
	var args SearchArgs
	array, _ := json.Marshal(results)
	args.Repositories = string(array)
	tmpl.ExecuteTemplate(w, "search.html", args)
}

// main sets the maximum number of CPUs that can be executing simultaneously,
// handles connexion and listen to port 4242.
func main() {
	runtime.GOMAXPROCS(8)
	http.HandleFunc("/", index)
	http.HandleFunc("/search", search)
	fs := http.FileServer(http.Dir("src/github.com/utay/ScalingoTest/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	err := http.ListenAndServe(":4242", nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
