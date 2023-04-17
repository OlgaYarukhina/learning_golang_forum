package main

import (
	"database/sql"
	"flag"
	models "forum/pkg"
	"github.com/mattn/go-sqlite3"
	_ "github.com/mattn/go-sqlite3"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type Application struct {
	ErrorLog      *log.Logger
	InfoLog       *log.Logger
	sqlError      sqlite3.Error
	TemplateCache map[string]*template.Template
	Session       map[string]models.Session
	Users         *models.UserModel
	Posts         *models.PostModel
	Categories    *models.CategoryModel
	Comment       *models.CommentModel
}

var app Application

func main() {

	addr := flag.String("addr", ":4000", "Network address HTTP")
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB("./pkg/Forum.db")
	if err != nil {
		errorLog.Fatal(err)
	}

	defer db.Close()

	templateCache, err := newTemplateCache("./ui/html/")

	app = Application{
		ErrorLog:      errorLog,
		InfoLog:       infoLog,
		TemplateCache: templateCache,
		Session:       map[string]models.Session{},
		Users:         &models.UserModel{DB: db},
		Posts:         &models.PostModel{DB: db},
		Categories:    &models.CategoryModel{DB: db},
		Comment:       &models.CommentModel{DB: db},
	}

	flag.Parse()

	infoLog.Printf("Starting forum on port: 4000", *addr)

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}

type neuteredFileSystem struct {
	fs http.FileSystem
}

func (nfs neuteredFileSystem) Open(path string) (http.File, error) {
	f, err := nfs.fs.Open(path)
	if err != nil {
		return nil, err
	}

	s, err := f.Stat()
	if s.IsDir() {
		index := filepath.Join(path, "index.html")
		if _, err := nfs.fs.Open(index); err != nil {
			closeErr := f.Close()
			if closeErr != nil {
				return nil, closeErr
			}

			return nil, err
		}
	}

	return f, nil
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
