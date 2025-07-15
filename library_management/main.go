package main

import (
	"fmt"
	"library_management/controllers"
	"library_management/models"
	"library_management/services"
	"os"
)

func main() {

	var libraryService services.LibraryManager = services.NewLibrary()

	libraryService.AddBook(models.Book{ID: 1, Title: "The Go Programming Language", Author: "Donovan & Kernighan", Status: "Available"})
	libraryService.AddBook(models.Book{ID: 2, Title: "A Tour of Go", Author: "Go Team", Status: "Available"})
	libraryService.AddBook(models.Book{ID: 3, Title: "Concurrency in Go", Author: "Katherine Cox-Buday", Status: "Available"})

	if lib, ok := libraryService.(*services.Library); ok {
		lib.AddMember(models.Member{ID: 101, Name: "Alice"})
		lib.AddMember(models.Member{ID: 102, Name: "Bob"})
	}

	controller := controllers.NewLibraryController(libraryService)

	for {
		printMenu()
		choice := controller.ReadStringInput("Enter your choice: ")

		switch choice {
		case "1":
			controller.ListAvailableBooksUI()
		case "2":
			controller.BorrowBookUI()
		case "3":
			controller.ReturnBookUI()
		case "4":
			controller.ListBorrowedBooksUI()
		case "5":
			controller.AddBookUI()
		case "6":
			fmt.Println("Exiting the application. Goodbye!")
			os.Exit(0)
		default:
			fmt.Println("Invalid choice. Please try again.")
		}
	}
}

func printMenu() {
	fmt.Println("\n--- Library Menu ---")
	fmt.Println("1. List Available Books")
	fmt.Println("2. Borrow a Book")
	fmt.Println("3. Return a Book")
	fmt.Println("4. List My Borrowed Books")
	fmt.Println("5. Add a New Book (Admin)")
	fmt.Println("6. Exit")
}
