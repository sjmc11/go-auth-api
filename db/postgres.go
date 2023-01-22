package db

import (
	"auth-api/helpers/env"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/logrusorgru/aurora/v4"
	"log"
)

type PostgresDBConnection struct {
	Db *pgx.Conn
}

var PgConnection PostgresDBConnection

func (p *PostgresDBConnection) InitDbConn() error {

	// &statement_cache_mode=describe
	dbConfig, err := pgx.ParseConfig("postgresql://" + env.Get("DBHOST") + ":" + env.Get("DBPORT") + "/" + env.Get("DBNAME") + "?user=" + env.Get("DBUSER") + "&password=" + env.Get("DBPASSWORD"))
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n\n", err)
	}

	PgConnection.Db, err = pgx.ConnectConfig(context.Background(), dbConfig)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n\n", err)
	}
	
	//return &PostgresDBConnection{
	//	Db: PgConnection.Db,
	//}, nil
	return nil
}

func (p *PostgresDBConnection) PingDB() {
	err := PgConnection.Db.Ping(context.Background())
	if err != nil {
		log.Fatal(aurora.Red(err.Error()))
	}
	fmt.Println(aurora.Green("✓ POSTGRES DB IS ALIVE : " + env.Get("DBPORT") + " @ " + env.Get("DBHOST")))
}
