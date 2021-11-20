package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type album struct {
	Id     int     `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

var albums = []album{
	{Id: 1, Title: "Separuh Aku", Price: 43000},
	{Id: 2, Title: "Yang Terdalam", Artist: "D'masiv", Price: 50000},
	{Id: 3, Title: "Luka di Hati", Artist: "Geisha", Price: 50000},
	{Id: 4, Title: "Salah", Artist: "Geisha", Price: 50000},
}

func main() {
	router := gin.Default()
	router.GET("/albums", getAlbums)
	router.POST("/albums", postAlbums)
	router.GET("/albums/:id", getSingleAlbum)

	router.Run("localhost:8080")
}

func getAlbums(ctx *gin.Context) {
	ctx.IndentedJSON(http.StatusOK, albums)
}

func postAlbums(ctx *gin.Context) {
	var newAlbum album

	if err := ctx.Bind(&newAlbum); err != nil {
		return
	}

	// ctx.Bind(&newAlbum)

	albums = append(albums, newAlbum)
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
		i, _ := strconv.Atoi(id)

		for _, v := range albums {
			if v.Id == i {
				result = v
			}
		}
		if result.Id == 0 {
			ctx.IndentedJSON(http.StatusNotFound, gin.H{
				"message": "Not Found",
			})
			return
		}
		ctx.IndentedJSON(http.StatusOK, result)
	}

}
