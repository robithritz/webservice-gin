package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"robithritz/web-service-gin/morestrings"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
)

type album struct {
	Id     int      `json:"id"`
	Title  string   `json:"title"`
	Artist *string  `json:"artist"`
	Price  *float64 `json:"price"`
}
type simpleMessage struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
}

var db *pgxpool.Pool
var err error

func main() {
	fmt.Println(morestrings.ReverseRunes("!oG ,olleH"))

	connStr := "postgres://postgres@localhost:5432/albums_db?sslmode=disable"

	db, err = pgxpool.Connect(context.Background(), connStr)

	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	pingError := db.Ping(context.Background())
	if pingError != nil {

		fmt.Println(pingError)
		os.Exit(0)
	}
	fmt.Println("Connected!")

	router := gin.Default()
	router.GET("/albums", getAlbums)
	router.POST("/albums", postAlbums)
	router.GET("/albums/:id", getSingleAlbum)

	router.Run("localhost:8080")
}

func getAlbums(ctx *gin.Context) {
	// var userid int
	// var username string
	var listData []album
	rows, err := db.Query(context.Background(), "select id, title, artist, price from master_album order by id ASC")

	if err != nil {
		fmt.Println(err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var obj album
		if err := rows.Scan(&obj.Id, &obj.Title, &obj.Artist, &obj.Price); err != nil {
			log.Fatal(err)
		}
		listData = append(listData, obj)
	}

	ctx.IndentedJSON(http.StatusOK, listData)
}

func postAlbums(ctx *gin.Context) {
	var newAlbum album
	decoder := json.NewDecoder(ctx.Request.Body)

	err := decoder.Decode(&newAlbum)
	if err != nil {
		fmt.Println(err)
	}

	// var artist = *newAlbum.Artist
	// var price = *newAlbum.Price
	// fmt.Println(artist)
	err = db.QueryRow(context.Background(), "insert into master_album (title, artist, price) values($1, $2, $3) returning id", newAlbum.Title, newAlbum.Artist, newAlbum.Price).Scan(&newAlbum.Id)
	if err != nil {
		fmt.Println(err)
	}

	// albums = append(albums, newAlbum)
	ctx.IndentedJSON(http.StatusCreated, newAlbum)

}

func getSingleAlbum(ctx *gin.Context) {
	id := ctx.Param("id")

	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "id required",
		})
	} else {
		var result album

		err := db.QueryRow(context.Background(), "select id, title, artist, price from master_album where id=$1", id).Scan(&result.Id, &result.Title, &result.Artist, &result.Price)

		if err != nil {
			fmt.Println(err)
			var resp simpleMessage
			if strings.Contains(err.Error(), "no rows") {
				resp = simpleMessage{
					Status:  false,
					Message: "data not found",
				}

				ctx.IndentedJSON(http.StatusNotFound, resp)
			} else {
				resp = simpleMessage{
					Status:  false,
					Message: "something went wrong",
				}
				ctx.IndentedJSON(http.StatusBadGateway, resp)
			}
			return
		}
		ctx.IndentedJSON(http.StatusOK, result)
	}

}
