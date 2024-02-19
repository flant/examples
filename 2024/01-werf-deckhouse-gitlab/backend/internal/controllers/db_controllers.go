package controllers

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"habr_app/internal/services"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

func RememberController(c *gin.Context) {
	dbType, dbPath := services.GetDBCredentials()

	db, err := sql.Open(dbType, dbPath)
	if err != nil {
		panic(err)
	}

	message := c.Query("message")
	name := c.Query("name")

	_, err = db.Exec("INSERT INTO talkers (message, name) VALUES (?, ?)",
		message, name)
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, gin.H{
		"Name":    name,
		"Message": message,
	})

	defer db.Close()
}

func SayController(c *gin.Context) {
	dbType, dbPath := services.GetDBCredentials()

	db, err := sql.Open(dbType, dbPath)
	if err != nil {
		panic(err)
	}

	result, err := db.Query("SELECT * FROM talkers")
	if err != nil {
		panic(err)
	}

	count := 0
	var data []map[string]string

	for result.Next() {
		count++
		var id int
		var message string
		var name string

		err = result.Scan(&id, &message, &name)
		if err != nil {
			panic(err)
		}

		data = append(data, map[string]string{
			"Name":    name,
			"Message": message})
	}
	if count == 0 {
		c.JSON(http.StatusOK, gin.H{
			"Error": "There are no messages from talkers!",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"Messages": data,
		})
	}
}
