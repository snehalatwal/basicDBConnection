package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/go-sql-driver/mysql"
)

type Albumn struct {
	ID     int64
	Title  string
	Artist string
	Price  float32
}

var db *sql.DB

func main() {
	//capture connection properties
	os.Setenv("DBUSER", "root")
	os.Setenv("DBPASS", "some_password")
	cfg := mysql.Config{
		User:   os.Getenv("DBUSER"),
		Passwd: os.Getenv("DBPASS"),
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: "recordings",
	}

	//Get a database handle

	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	// fmt.Println(cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr :=
		db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("connected!")

	albums, err := albumsByArtist("Gerry Mulligan")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Albums found : %v\n", albums)

	alb, err := artistById(1)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Album found : %v\n", alb)

	ablId, err:= addAlbum(Albumn{
		Title : "The Show",
		Artist: "Niall Horan",
		Price :50,
	})
	if err!=nil{
		log.Fatal(err)
	}
	fmt.Printf("ID of added: %v\n", ablId)
}

// returns the multiple rows
func albumsByArtist(name string) ([]Albumn, error) {
	// an albums slice to hold data from returned rows
	var albums []Albumn
	rows, err := db.Query("SELECT * FROM album WHERE artist =?", name)
	if err != nil {
		return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
	}
	defer rows.Close()

	// loop through rows, using scan to assisgn column data to struct file
	for rows.Next() {
		var alb Albumn
		if err := rows.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
			return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
		}
		albums = append(albums, alb)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
	}
	return albums, nil
}

// return album by id ->return the single row
func artistById(id int64) (Albumn, error) {
	//an albumn to hold from returned row
	var alb Albumn
	row := db.QueryRow("SELECT * FROM albumn WHERE id= ?", id)
	if err := row.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
		if err == sql.ErrNoRows {
			return alb, fmt.Errorf("albumnById %d: no such album", id)
		}
	}
	return alb, nil
}

// addAlnum adds the specified album to the database
// return the album ID

func addAlbum(alb Albumn) (int64, error){
	result, err :=db.Exec("INSERT INTO album(title,artist, price) VALUES(?, ?, ?)", alb.Title, alb.Artist, alb.Price)
	if err!=nil{
		return 0, fmt.Errorf("addAlbum: %v", err)
	}
	id, err :=result.LastInsertId()
	if err!=nil{
		return 0, fmt.Errorf("addAlbum: %v",err)

	}
	return id, nil
}
