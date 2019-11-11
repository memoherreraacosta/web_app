// Run PostgreSQL server:
//  docker run -e POSTGRES_PASSWORD="" -p 5432:5432 postgres
// Monitor running processes:
//   watch -n 1 'echo "select pid,query_start,state,query from pg_stat_activity;" | psql -h localhost -U postgres
//
// For all handlers, call to db takes 5 seconds,
//
// Three endpoints:
//  - "/" - take 5 seconds
//  - "/ctx" - take 1 seconds, due to 1 second cancellation policy
//  - "/disconnect" - aborts as soon as client disconnected
//
// To test, run:
//   curl http://localhost:3000/
//   curl http://localhost:3000/ctx
//   curl http://localhost:3000/disconnect

package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"time"

	_ "github.com/lib/pq"
)

type server struct {
	db *sql.DB
}

func (s *server) handler(w http.ResponseWriter, r *http.Request) {
	// slow 5 seconds query
	_, err := s.db.Exec("SELECT pg_sleep(5)")
	if err != nil {
		log.Println("[ERROR]", err)
		w.WriteHeader(http.StatusBadRequest)
	}

	w.Write([]byte("ok"))
}

func (s *server) handlerCtx(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	ctx, _ = context.WithTimeout(ctx, 1*time.Second)

	// slow 5 seconds query
	_, err := s.db.ExecContext(ctx, "SELECT pg_sleep(5)")
	if err != nil {
		log.Println("[ERROR]", err)
		w.WriteHeader(http.StatusBadRequest)
	}

	w.Write([]byte("ok"))
}

func (s *server) handlerDisconnect(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	ctx, cancelFunc := context.WithCancel(ctx)

	// in case of client disconnect, cancel context
	if cn, ok := w.(http.CloseNotifier); ok {
		go func() {
			<-cn.CloseNotify()
			cancelFunc()
		}()
	}

	// slow 5 seconds query
	_, err := s.db.ExecContext(ctx, "SELECT pg_sleep(5)")
	if err != nil {
		log.Println("[ERROR]", err)
		w.WriteHeader(http.StatusBadRequest)
	}

	w.Write([]byte("ok"))
}

func main() {
	db, err := sql.Open("postgres", "user=postgres sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	s := server{db: db}

	http.HandleFunc("/", s.handler)
	http.HandleFunc("/ctx", s.handlerCtx)
	http.HandleFunc("/disconnect", s.handlerDisconnect)
	log.Println("Starting server on :3000...")
	log.Fatal(http.ListenAndServe("localhost:3000", nil))
	if err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
