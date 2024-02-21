package db

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/varga-lp/data/config"
)

const (
	migrations = `
	CREATE SEQUENCE IF NOT EXISTS futures_id_seq;

	CREATE TABLE IF NOT EXISTS symbols (
		id BIGSERIAL PRIMARY KEY,
		symbol VARCHAR(30) NOT NULL
	);

	CREATE UNIQUE INDEX IF NOT EXISTS index_symbols_on_symbol ON symbols (symbol);

	CREATE TABLE IF NOT EXISTS futures (
		id BIGINT NOT NULL DEFAULT nextval('futures_id_seq'),
		symbol VARCHAR(30) NOT NULL,
		close_time BIGINT NOT NULL,
		low NUMERIC(20, 12) NOT NULL,
		high NUMERIC(20, 12) NOT NULL,
		close NUMERIC(20, 12) NOT NULL,
		volume NUMERIC(20, 12) NOT NULL,
		notr BIGINT NOT NULL,
		PRIMARY KEY (id, symbol)
	) PARTITION BY HASH(symbol);

	CREATE UNIQUE INDEX IF NOT EXISTS index_futures_on_sym_ct ON futures (symbol, close_time);

	CREATE TABLE IF NOT EXISTS futures_p1 PARTITION OF futures FOR VALUES WITH (modulus 10, remainder 0);
	CREATE TABLE IF NOT EXISTS futures_p2 PARTITION OF futures FOR VALUES WITH (modulus 10, remainder 1);
	CREATE TABLE IF NOT EXISTS futures_p3 PARTITION OF futures FOR VALUES WITH (modulus 10, remainder 2);
	CREATE TABLE IF NOT EXISTS futures_p4 PARTITION OF futures FOR VALUES WITH (modulus 10, remainder 3);
	CREATE TABLE IF NOT EXISTS futures_p5 PARTITION OF futures FOR VALUES WITH (modulus 10, remainder 4);
	CREATE TABLE IF NOT EXISTS futures_p6 PARTITION OF futures FOR VALUES WITH (modulus 10, remainder 5);
	CREATE TABLE IF NOT EXISTS futures_p7 PARTITION OF futures FOR VALUES WITH (modulus 10, remainder 6);
	CREATE TABLE IF NOT EXISTS futures_p8 PARTITION OF futures FOR VALUES WITH (modulus 10, remainder 7);
	CREATE TABLE IF NOT EXISTS futures_p9 PARTITION OF futures FOR VALUES WITH (modulus 10, remainder 8);
	CREATE TABLE IF NOT EXISTS futures_p10 PARTITION OF futures FOR VALUES WITH (modulus 10, remainder 9);

	CREATE TABLE IF NOT EXISTS correlations (
		symbol1 VARCHAR(30) NOT NULL,
		symbol2 VARCHAR(30) NOT NULL,
		price_cor NUMERIC(7, 6) NOT NULL,
		volume_cor NUMERIC(7, 6) NOT NULL,
		not_cor NUMERIC(7, 6) NOT NULL,
		avg_cor NUMERIC(7, 6) NOT NULL,
		updated_at TIMESTAMP NOT NULL,
		PRIMARY KEY (symbol1, symbol2)
	);
	`
)

type Conn struct {
	Pool *pgxpool.Pool
}

func (c *Conn) Ready() {
	log.Println("Database connection is ready.")
}

var (
	GConn *Conn
)

func init() {
	pool, err := pgxpool.New(context.Background(), config.DatabaseUrl)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	GConn = &Conn{Pool: pool}
	migrate()
}

func migrate() {
	if _, err := GConn.Pool.Exec(context.Background(), migrations); err != nil {
		log.Fatalf("Unable to run migrations: %v\n", err)
	}
}
