package main

import (
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/gorilla/mux"
)

type Article struct {
	Id        uint16
	Title     string
	Anons     string
	Full_text string
}

var posts = []Article{}
var showPost = Article{}

func index(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprint(w, "Go is super star!")
	tmpl, err := template.ParseFiles("templates/index.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprint(w, "Error")
	}
	tmpl.ExecuteTemplate(w, "index", nil)
}

func article(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/posts.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprint(w, "Error")
	}
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:8889)/Golang")
	if err != nil {
		panic(err)
	}
	res, err := db.Query("Select * from Articles")
	if err != nil {
		panic(err)
	}
	posts = []Article{}
	for res.Next() {
		var post Article
		err = res.Scan(&post.Id, &post.Title, &post.Anons, &post.Full_text)
		if err != nil {
			panic(err)
		}
		posts = append(posts, post)
	}
	tmpl.ExecuteTemplate(w, "posts", posts)
}

func create(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/create.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprint(w, "Error")
	}
	tmpl.ExecuteTemplate(w, "create", nil)
}

func contacts_page(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Go is contacts page")
}
func save(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	anons := r.FormValue("anons")
	text := r.FormValue("full_text")
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:8889)/Golang")
	if err != nil {
		panic(err)
	}

	defer db.Close()
	if title == "" || anons == "" || text == "" {
		fmt.Fprintf(w, "не все заполнено")
	} else {
		insert, err := db.Query(fmt.Sprintf("Insert into Articles (`Title`, `Anons`, `Full_text`) values('%s','%s','%s')", title, anons, text))
		if err != nil {
			panic(err)
		}
		defer insert.Close()

		http.Redirect(w, r, "/posts", http.StatusSeeOther)
	}
}

func show_Post(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:8889)/Golang")
	if err != nil {
		panic(err)
	}

	defer db.Close()
	tmpl, err := template.ParseFiles("templates/show.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprint(w, "Error")
	}
	res, err := db.Query(fmt.Sprintf("Select * from Articles where `id` = '%s'", vars["id"]))
	if err != nil {
		panic(err)
	}
	showPost = Article{}
	for res.Next() {
		var post Article
		err = res.Scan(&post.Id, &post.Title, &post.Anons, &post.Full_text)
		if err != nil {
			panic(err)
		}
		showPost = post
	}
	tmpl.ExecuteTemplate(w, "show", showPost)
}

func mySQL() {
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:8889)/Golang")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	insert, err := db.Query("Insert into `Users` (`Name`, `Age`) values ('Alex', 25)")
	if err != nil {
		panic(err)
	}

	defer insert.Close()
	fmt.Println("успешно подключено к базе данных")
}

func handleRequest() {
	rtr := mux.NewRouter()
	rtr.HandleFunc("/", index).Methods("GET")
	rtr.HandleFunc("/create", create).Methods("GET")
	rtr.HandleFunc("/contacts", contacts_page)
	rtr.HandleFunc("/save", save).Methods("POST")
	rtr.HandleFunc("/posts", article)
	rtr.HandleFunc("/post/{id:[0-9]+}", show_Post).Methods("GET")

	http.Handle("/", rtr)
	http.Handle("/wwwroot/",
		http.StripPrefix("/wwwroot/", http.FileServer(http.Dir("./wwwroot/"))))
	http.ListenAndServe(":8080", nil)
}

func main() {
	handleRequest()
	mySQL()
}
