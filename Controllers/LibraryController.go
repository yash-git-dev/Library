package controllers

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"os"
	constants "root/Constants"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type Book struct {
	Name            string `json:"name"`
	Author          string `json:"author"`
	PublicationYear int    `json:"publication_year"`
}

func BooksGET(c *gin.Context) {
	var totalData []Book

	username := c.GetString("username")

	isAdmin := checkAdminRole(username)

	regularData, err := readCSV("regularUser.csv")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read regularUser.csv:" + err.Error()})
		return
	}

	totalData = append(totalData, regularData...)

	if isAdmin {
		adminData, err := readCSV("adminUser.csv")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read adminUser.csv:" + err.Error()})
			return
		}

		totalData = append(totalData, adminData...)
	}

	c.JSON(http.StatusOK, gin.H{"books": totalData})
}

func readCSV(filename string) ([]Book, error) {
	var books []Book

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)

	if _, err := reader.Read(); err != nil {
		return nil, err
	}

	lines, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	for _, line := range lines {
		year, err := strconv.Atoi(line[2])
		if err != nil {
			return nil, err
		}

		book := Book{
			Name:            line[0],
			Author:          line[1],
			PublicationYear: year,
		}
		books = append(books, book)
	}

	return books, nil
}

func checkAdminRole(username string) bool {
	for _, userDB := range constants.Users {
		if username == userDB.Username && userDB.IsAdmin {
			return true
		}
	}

	return false
}

func BookPOST(c *gin.Context) {
	var book Book
	if err := c.ShouldBindJSON(&book); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if book.Name == "" || book.Author == "" || book.PublicationYear <= 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid book data"})
		return
	}

	username := c.GetString("username")

	isAdmin := checkAdminRole(username)

	if !isAdmin {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if err := writeToCSV(book); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to add book: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Book added successfully"})

}

func writeToCSV(book Book) error {
	file, err := os.OpenFile("regularUser.csv", os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file:%v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("failed to read file:%v", err)
	}

	var lastIndex int
	for i := len(records) - 1; i >= 0; i-- {
		if len(records[i]) > 0 {
			lastIndex = i
			break
		}
	}

	records[lastIndex] = []string{book.Name, book.Author, strconv.Itoa(book.PublicationYear)}

	_, err = file.Seek(0, 0)
	if err != nil {
		return fmt.Errorf("failed to seek at initial point of file:%v", err)
	}

	err = file.Truncate(0)
	if err != nil {
		return fmt.Errorf("failed to truncate file:%v", err)
	}

	writer := csv.NewWriter(file)

	if err := writer.WriteAll(records); err != nil {
		return fmt.Errorf("failed to write file:%v", err)
	}

	return nil
}

func BookDELETE(c *gin.Context) {
	bookName := c.Param("bookName")
	if bookName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid URI data"})
		return
	}

	username := c.GetString("username")

	isAdmin := checkAdminRole(username)

	if !isAdmin {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if err := deleteFromCSV(bookName); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete book" + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Book deleted successfully"})

}

func deleteFromCSV(bookName string) error {
	file, err := os.OpenFile("regularUser.csv", os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	lines, err := reader.ReadAll()
	if err != nil {
		return err
	}

	var updatedLines [][]string
	for _, line := range lines {
		if strings.EqualFold(line[0], bookName) {
			continue
		}
		updatedLines = append(updatedLines, line)
	}

	file.Truncate(0)
	file.Seek(0, 0)
	writer := csv.NewWriter(file)
	err = writer.WriteAll(updatedLines)
	if err != nil {
		return err
	}
	writer.Flush()

	return nil
}
