// Copyright (c) 2021 Dean Jackson <deanishe@deanishe.net>
// MIT Licence applies http://opensource.org/licenses/MIT
// Created on 2021-02-05

package main

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/text/unicode/norm"
)

// Repo contains Unicode data.
type Repo struct {
	db *sqlx.DB
}

// Open opens a database, creating it if necessary.
func Open(path string) (*Repo, error) {
	var (
		db  *sqlx.DB
		err error
	)
	if db, err = sqlx.Connect("sqlite3", path); err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	// actually create a database connection
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("ping database: %w", err)
	}
	if err = setupDatabase(db); err != nil {
		return nil, fmt.Errorf("set up database: %w", err)
	}
	return &Repo{db: db}, nil
}

// Close closes the database.
func (repo *Repo) Close() error {
	return repo.db.Close()
}

// Search searches the repo.
func (repo *Repo) Search(query string) ([]Codepoint, error) {
	query = repo.prepareQuery(query)
	var (
		cps []Codepoint
		q   = `
			SELECT
				c.hex, c.dec, c.category, c.name, c.aliases, c.html_entity,
				s.rank
			FROM codepoints AS c
			INNER JOIN (
				SELECT
					rowid,
					bm25(search) AS rank
				FROM search
				WHERE search MATCH ?
			) AS s
			ON c.rowid = s.rowid
			ORDER BY s.rank
			LIMIT 200`
		// search for a single character
		q2 = `
			SELECT
				c.hex, c.dec, c.category, c.name, c.aliases, c.html_entity
			FROM codepoints AS c
			WHERE c.dec = ?
		`
		err error
	)
	if runes := []rune(query); len(runes) == 1 {
		err = repo.db.Select(&cps, q2, runes[0])
	} else {
		err = repo.db.Select(&cps, q, query)
	}

	if err != nil {
		return nil, fmt.Errorf("SQL query failed: %w", err)
	}

	log.Printf("[repo] %d result(s) for %q", len(cps), query)
	return cps, nil
}

// Import imports Codepoints into the database.
func (repo *Repo) Import(cps []Codepoint) error {
	var (
		i, j  int
		batch = 5000
		size  = len(cps)
	)
	for {
		if i >= size {
			break
		}
		j = i + batch
		if j >= size-1 {
			j = size - 1
		}
		log.Printf("[repo] importing items %d to %d of %d ...", i, j, size)
		if err := repo.importCodepoints(cps[i:j]); err != nil {
			return err
		}
		i += batch
	}
	return nil
}

func (repo *Repo) importCodepoints(cps []Codepoint) error {
	var (
		q = `INSERT OR REPLACE INTO codepoints (
			hex, dec, category, name, aliases, html_entity
		) VALUES (
			:hex, :dec, :category, :name, :aliases, :html_entity
		);`
		stmt *sqlx.NamedStmt
		tx   *sqlx.Tx
		err  error
	)
	if tx, err = repo.db.Beginx(); err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback()

	if stmt, err = tx.PrepareNamed(q); err != nil {
		return fmt.Errorf("prepare SQL %q: %w", q, err)
	}

	for _, cp := range cps {
		if _, err := stmt.Exec(cp); err != nil {
			return fmt.Errorf("insert codepoint: %#v: %w", cp, err)
		}
		// log.Printf("[repo] inserted %#v", cp)
	}
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}

// setupDatabase initialises and updates the database schema.
func setupDatabase(db *sqlx.DB) error {
	// turn on write-ahead log and store temp tables on disk.
	_, err := db.Exec(`
		PRAGMA journal_mode=WAL;
		PRAGMA temp_store=1;
		PRAGMA foreign_keys=ON;
		PRAGMA encoding='UTF-8';
	`)
	if err != nil {
		return fmt.Errorf("execute PRAGMAs: %w", err)
	}

	tx, err := db.Beginx()
	if err != nil {
		return fmt.Errorf("begin Tx: %w", err)
	}
	defer tx.Rollback()

	var version int
	if err = tx.Get(&version, `PRAGMA user_version;`); err != nil {
		return fmt.Errorf("PRAGMA user_version: %w", err)
	}
	log.Printf("[repo] current schema: %d", version)

	switch version {
	case 0:
		log.Printf("initialising database ...")
		_, err = tx.Exec(`
			-- data table
			CREATE TABLE IF NOT EXISTS codepoints (
				hex TEXT PRIMARY KEY,
				dec INTEGER,
				category TEXT,
				name TEXT,
				aliases TEXT,
				html_entity TEXT
			);
			-- CREATE INDEX IF NOT EXISTS idx_hex ON codepoints(hex);

			-- full-text search
			CREATE VIRTUAL TABLE IF NOT EXISTS search USING FTS5(
				name,
				aliases,
				html_entity,
				hex,
				dec,
				content = 'codepoints',
				tokenize = "unicode61 remove_diacritics 1 separators '-'"
			);

			CREATE TRIGGER IF NOT EXISTS codepoints_ai AFTER INSERT ON codepoints
			BEGIN
				INSERT INTO search (rowid, name, aliases, html_entity, hex, dec)
					VALUES (new.rowid, new.name, new.aliases, new.html_entity, new.hex, new.dec);
			END;
			CREATE TRIGGER IF NOT EXISTS codepoints_ad AFTER DELETE ON codepoints
			BEGIN
				INSERT INTO search (search, rowid, name, aliases, html_entity, hex, dec)
					VALUES ('delete', old.rowid, old.name, old.aliases, old.html_entity, old.hex, old.dec);
			END;
			CREATE TRIGGER IF NOT EXISTS codepoints_au AFTER UPDATE ON codepoints
			BEGIN
				INSERT INTO search (search, rowid, name, aliases, html_entity, hex, dec)
					VALUES ('delete', old.rowid, old.name, old.aliases, old.html_entity, old.hex, old.dec);
				INSERT INTO search (rowid, name, aliases, html_entity, hex, dec)
					VALUES (new.rowid, new.name, new.aliases, new.html_entity, new.hex, new.dec);
			END;

			PRAGMA user_version = 1;
		`)
		if err != nil {
			return fmt.Errorf("migrate to version 1: %w", err)
		}
		fallthrough
	case 1: //  future migration
		break
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("commit schema updates: %w", err)
	}
	var n int
	if err = db.Get(&n, `PRAGMA user_version;`); err != nil {
		return fmt.Errorf("PRAGMA user_version: %w", err)
	}

	if n > version {
		if _, err := db.Exec(`VACUUM;`); err != nil {
			return fmt.Errorf("VACUUM database: %w", err)
		}
	}
	return nil
}

func (repo *Repo) prepareQuery(s string) string {
	if s2 := norm.NFC.String(s); s2 != s {
		log.Printf("[repo] query %X normalised to %X", s, s2)
		s = s2
	}
	if len([]rune(s)) == 1 {
		return s
	}

	if s2, stop := specialQuery(s); s2 != s {
		log.Printf("query transform: %q  -->  %q", s, s2)
		if stop {
			return s2
		}
		s = s2
	}

	// add wildcard to end of query to get more matches
	if strings.HasSuffix(s, " ") || strings.HasSuffix(s, "*") {
		return s
	}
	return s + "*"
}

var rxQuery = regexp.MustCompile(`(?i)^0x([A-F0-9]{4,8}) /(.*)$`)

// parse "details" queries of format U+00AB / [<query>].
func parseQuery(s string) (charName, query string, backup bool) {
	if m := rxQuery.FindStringSubmatch(s); m != nil {
		charName = m[1]
		query = m[2]
		if query == "" {
			backup = true
		}
		query = strings.TrimLeft(query, " ")
		return
	}
	query = s
	return
}
