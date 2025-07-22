package controllers

import (
	"bufio"
	"fmt"
	"library_management/models"
	"library_management/services"
	"os"
	"strconv"
	"strings"
)

type LibraryController struct {
	service services.LibraryManager
	reader  *bufio.Reader
}

func NewLibraryController(service services.LibraryManager) *LibraryController {
	return &LibraryController{
		service: service,
		reader:  bufio.NewReader(os.Stdin),
	}
}
func (c *LibraryController) readIntInput(prompt string) int {
	for {
		fmt.Print(prompt)
		input, _ := c.reader.ReadString('\n')
		input = strings.TrimSpace(input)
		num, err := strconv.Atoi(input)
		if err == nil {
			return num
		}
		fmt.Println("Invalid input. Please enter a number.")
	}
}

func (c *LibraryController) ReadStringInput(prompt string) string {
	fmt.Print(prompt)
	input, _ := c.reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func (c *LibraryController) AddBookUI() {
	fmt.Println("\n--- Add a New Book ---")
	id := c.readIntInput("Enter Book ID: ")
	title := c.ReadStringInput("Enter Title: ")
	author := c.ReadStringInput("Enter Author: ")
	book := models.Book{ID: id, Title: title, Author: author, Status: "Available"}
	c.service.AddBook(book)
	fmt.Println("Book added successfully!")
}

func (c *LibraryController) ListAvailableBooksUI() {
	fmt.Println("\n--- Available Books ---")
	books := c.service.ListAvailableBooks()
	if len(books) == 0 {
		fmt.Println("No available books.")
		return
	}
	for _, b := range books {
		fmt.Printf("ID: %d, Title: %s, Author: %s\n", b.ID, b.Title, b.Author)
	}
}
func (c *LibraryController) BorrowBookUI() {
	fmt.Println("\n--- Borrow a Book ---")
	bookID := c.readIntInput("Enter Book ID to borrow: ")
	memberID := c.readIntInput("Enter your Member ID: ")
	err := c.service.BorrowBook(bookID, memberID)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Println("Book borrowed successfully!")
	}
}

func (c *LibraryController) ReturnBookUI() {
	fmt.Println("\n--- Return a Book ---")
	bookID := c.readIntInput("Enter Book ID to return: ")
	memberID := c.readIntInput("Enter your Member ID: ")
	err := c.service.ReturnBook(bookID, memberID)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Println("Book returned successfully!")
	}
}

func (c *LibraryController) ListBorrowedBooksUI() {
	fmt.Println("\n--- My Borrowed Books ---")
	memberID := c.readIntInput("Enter your Member ID: ")
	books, err := c.service.ListBorrowedBooks(memberID)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	if len(books) == 0 {
		fmt.Println("You have not borrowed any books.")
		return
	}
	for _, b := range books {
		fmt.Printf("ID: %d, Title: %s, Author: %s\n", b.ID, b.Title, b.Author)
	}
}
