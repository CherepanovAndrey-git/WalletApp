package database

//
//import (
//	"database/sql"
//	"log"
//	"os"
//	"testing"
//
//	_ "github.com/lib/pq"
//)
//
//var testQueries *Queries
//
//func TestMain(m *testing.M) {
//	conn, err := sql.Open("postgres", "postgres://postgres:postgres@db:5432/wallet?sslmode=disable")
//	if err != nil {
//		log.Fatal("cannot connect to db:", err)
//	}
//
//	testQueries = New(conn)
//
//	os.Exit(m.Run())
//}
