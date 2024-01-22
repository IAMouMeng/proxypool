package data

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
	"strings"
)

type Engine struct {
	Db *sql.DB
}

var GlobalEngine *Engine

type ProxyData struct {
	Proxy   string `json:"proxy"`
	Country string `json:"country"`
	Type    string `json:"type"`
}

func InitDb(dbFileName string) {
	db, err := connectDb(dbFileName)
	if err != nil {
		log.Fatal(err)
	}
	GlobalEngine = &Engine{Db: db}
}

func connectDb(dbFileName string) (*sql.DB, error) {
	if _, err := os.Stat(dbFileName); os.IsNotExist(err) {
		err := createDatabase(dbFileName)
		if err != nil {
			return nil, err
		}
	}

	db, err := sql.Open("sqlite3", dbFileName)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func createDatabase(dbFileName string) error {
	db, err := sql.Open("sqlite3", dbFileName)
	if err != nil {
		return err
	}
	defer db.Close()

	createTableSQL := `
	CREATE TABLE IF NOT EXISTS proxy (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		proxy CHAR(50),
		country CHAR(5),
		type CHAR(10)
	);
	`
	_, err = db.Exec(createTableSQL)
	if err != nil {
		return err
	}

	return nil
}

func (e *Engine) InsertProxy(proxyData ProxyData) error {

	checkExistenceSQL := "SELECT COUNT(*) FROM proxy WHERE proxy = ?"
	var count int
	err := e.Db.QueryRow(checkExistenceSQL, proxyData.Proxy).Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		return fmt.Errorf("the proxy has already exists %s", proxyData.Proxy)
	}

	insertDataSQL := "INSERT INTO proxy (proxy, country, type) VALUES (?, ?, ?)"
	_, err = e.Db.Exec(insertDataSQL, proxyData.Proxy, proxyData.Country, proxyData.Type)
	if err != nil {
		return err
	}
	return nil
}

func (e *Engine) GetProxyList(country string, _type string) ([]ProxyData, error) {
	var query string
	var args []interface{}

	whereClause := []string{}
	if country != "" {
		whereClause = append(whereClause, "country = ?")
		args = append(args, country)
	}
	if _type != "" {
		whereClause = append(whereClause, "type = ?")
		args = append(args, _type)
	}

	if len(whereClause) > 0 {
		query = "SELECT proxy, country, type FROM proxy WHERE " + strings.Join(whereClause, " AND ")
	} else {
		query = "SELECT proxy, country, type FROM proxy"
	}

	rows, err := e.Db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var proxyList []ProxyData

	for rows.Next() {
		var proxyData ProxyData
		err := rows.Scan(&proxyData.Proxy, &proxyData.Country, &proxyData.Type)
		if err != nil {
			return nil, err
		}
		proxyList = append(proxyList, proxyData)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return proxyList, nil
}

func (e *Engine) DeleteProxy(proxy string) error {
	checkExistenceSQL := "SELECT COUNT(*) FROM proxy WHERE proxy = ?"
	var count int
	err := e.Db.QueryRow(checkExistenceSQL, proxy).Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		deleteDataSQL := "DELETE FROM proxy WHERE proxy = ?"
		_, err := e.Db.Exec(deleteDataSQL, proxy)
		if err != nil {
			return err
		}
	}

	return nil
}

func (e *Engine) CountProxiesByTypeAndCountry() (map[string]map[string]int, error) {
	query := "SELECT type, country, COUNT(*) FROM proxy GROUP BY type, country"

	rows, err := e.Db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	counts := make(map[string]map[string]int)

	for rows.Next() {
		var _type, country string
		var count int
		err := rows.Scan(&_type, &country, &count)
		if err != nil {
			return nil, err
		}

		if counts[_type] == nil {
			counts[_type] = make(map[string]int)
		}

		counts[_type][country] = count
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return counts, nil
}
