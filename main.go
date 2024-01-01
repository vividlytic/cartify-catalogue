package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"cartify/catalogue/infrastructure/repository"
	"cartify/catalogue/interfaces"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

func main() {

	flag.Parse()

	dsn := os.Getenv("DATABASE")
	if dsn == "" {
		dsn = "catalogue_user:default_password@tcp(127.0.0.1:3306)/booksdb"
	}

	db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	sequenceRepository := repository.NewSequenceRepository(db)
	bookRepository := repository.NewBookRepository(db, sequenceRepository)

	server := interfaces.NewServer(interfaces.ServerParams{
		BookRepository: bookRepository,
	})

	var stop = make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM)
	signal.Notify(stop, syscall.SIGINT)

	go func() {
		sig := <-stop
		log.Printf("%v", sig)
		server.GracefulStop()
	}()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("server listening at %v", lis.Addr())

	if err = server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	log.Println("shutdown success, bye!")

}
