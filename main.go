package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
)

var SQLITE_FILE_NAME = "./shortenedURLs.db"
var SQLITE_DSN = "file:" + SQLITE_FILE_NAME + "?mode=rwc"
var HTTP_PORT = ":8080"
var db *sql.DB

func init(){
	httpInit()
	err:= sqlInit()
	if err != nil {
		log.Fatal(err)
	}
}

func httpInit(){
	http.HandleFunc("/", homePageHandler)
	http.HandleFunc("/help", helpPageHandler)
}

func shutdown(){
	db.Close()
}

func sqlInit() error {
	db, err := sql.Open("sqlite3", SQLITE_DSN)
	if err != nil {
		return err
	}

	sqlStmt := `
	 CREATE TABLE IF NOT EXISTS urls
  (
     id        INTEGER NOT NULL PRIMARY KEY,
     url       TEXT NOT NULL,
     shortcode TEXT NOT NULL UNIQUE
  );`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		return err
	}
	return nil
}

func homePageHandler(w http.ResponseWriter, r *http.Request) {
	response := `
        <h1>shortURL</h1>
        <p class="lead">Shorten your urls here!</p>
<form>
  <div class="form-group">
    <input type="url" class="form-control" id="urlField" placeholder="Enter url">
  </div>
  <button type="submit" class="btn btn-primary">Submit</button>
</form>
`
	fmt.Fprintf(w,
		pageHeader() + response + pageFooter())
		//r.URL.Path[1:])
}

func helpPageHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Need help?")
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := tx.Prepare("insert into foo(id, name) values(?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	for i := 0; i < 100; i++ {
		_, err = stmt.Exec(i, fmt.Sprintf("こんにちわ世界%03d", i))
		if err != nil {
			log.Fatal(err)
		}
	}
	tx.Commit()

	rows, err := db.Query("select id, name from foo")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var name string
		err = rows.Scan(&id, &name)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(id, name)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	stmt, err = db.Prepare("select name from foo where id = ?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	var name string
	err = stmt.QueryRow("3").Scan(&name)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(name)

	_, err = db.Exec("delete from foo")
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec("insert into foo(id, name) values(1, 'foo'), (2, 'bar'), (3, 'baz')")
	if err != nil {
		log.Fatal(err)
	}

	rows, err = db.Query("select id, name from foo")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var name string
		err = rows.Scan(&id, &name)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(id, name)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	log.Fatal(http.ListenAndServe(HTTP_PORT, nil))
	shutdown()
}