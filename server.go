package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	_ "github.com/go-sql-driver/mysql"
	"html/template"
	"log"
	"net/http"
	"regexp"
)

type Todo struct {
	Id    int
	Title string
}

type TodoPageData struct {
	Todos []*Todo
}

type route struct {
	pattern *regexp.Regexp
	verb    string
	handler http.Handler
}

type RegexpHandler struct {
	routes []*route
}

type Server struct {
	db *sql.DB
}

func (h *RegexpHandler) HandleFunc(r string, v string, handler func(http.ResponseWriter, *http.Request)) {
	re := regexp.MustCompile(r)
	h.routes = append(h.routes, &route{re, v, http.HandlerFunc(handler)})
}

func (h *RegexpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, route := range h.routes {
		if route.pattern.MatchString(r.URL.Path) && route.verb == r.Method {
			route.handler.ServeHTTP(w, r)
			return
		}
	}
	http.NotFound(w, r)
}

func (s *Server) homeHandler(res http.ResponseWriter, req *http.Request) {
	var todos []*Todo

	rows, err := s.db.Query("SELECT id, title FROM todos ORDER BY created_at DESC")
	error_check(res, err)
	for rows.Next() {
		todo := &Todo{}
		rows.Scan(&todo.Id, &todo.Title)
		todos = append(todos, todo)
	}
	rows.Close()

	renderTemplate(res, "index", TodoPageData{Todos: todos})
}

func (s *Server) saveHandler(res http.ResponseWriter, req *http.Request) {
	title := req.FormValue("title")

	_, err := s.db.Exec("INSERT INTO todos(title) VALUES(?)", title)
	error_check(res, err)

	http.Redirect(res, req, "/", http.StatusFound)
}

func (s *Server) clearHandler(res http.ResponseWriter, req *http.Request) {
	_, err := s.db.Exec("DELETE FROM todos")
	error_check(res, err)

	http.Redirect(res, req, "/", http.StatusFound)
}

func (s *Server) assets(res http.ResponseWriter, req *http.Request) {
	http.ServeFile(res, req, req.URL.Path[1:])
}

func renderTemplate(res http.ResponseWriter, tmpl string, data TodoPageData) {
	t, _ := template.ParseFiles(tmpl + ".html")
	t.Execute(res, data)
}

func error_check(res http.ResponseWriter, err error) bool {
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return true
	}
	return false
}

func getDatabaseConf() string {
	file, _ := os.Open("config.json")
	decoder := json.NewDecoder(file)
	conf := Configuration{}
	err := decoder.Decode(&conf)

	if err != nil {
		fmt.Println("error:", err)
	}

	// "root:@tcp(127.0.0.1:3306)/sphere"
	return conf.User + ":" + conf.Password + "@tcp(" + conf.Host + ":3306)/" + conf.Database
}

type Configuration struct {
    Host    	string
    Database  string
    User			string
    Password	string
}

func main() {
	confUrl := getDatabaseConf();

	db, err := sql.Open("mysql", confUrl)
	if err != nil {
		log.Fatal(err)
	}
	db.SetMaxIdleConns(100)
	defer db.Close()

	server := &Server{db: db}

	reHandler := new(RegexpHandler)

	reHandler.HandleFunc("/clear/$", "GET", server.clearHandler)
	reHandler.HandleFunc("/save/$", "POST", server.saveHandler)
	reHandler.HandleFunc(".*.[js|css|png|eof|svg|ttf|woff]", "GET", server.assets)
	reHandler.HandleFunc("/", "GET", server.homeHandler)

	fmt.Println("Starting server on port 8100")
	http.ListenAndServe(":8100", reHandler)
}
