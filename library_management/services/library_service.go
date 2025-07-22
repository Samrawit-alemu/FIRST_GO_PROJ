package services

import (
	"errors"
	"library_management/models"
)

type LibraryManager interface {
	AddBook(book models.Book)
	RemoveBook(bookID int)
	BorrowBook(bookID int, memberID int) error
	ReturnBook(bookID int, memberID int) error
	ListAvailableBooks() []models.Book
	ListBorrowedBooks(memberID int) ([]models.Book, error)
}

type Library struct {
	books   map[int]models.Book
	members map[int]models.Member
}

func NewLibrary() *Library {
	return &Library{
		books:   make(map[int]models.Book),
		members: make(map[int]models.Member),
	}
}

func (l *Library) AddBook(book models.Book) {
	l.books[book.ID] = book
}

func (l *Library) RemoveBook(bookID int) {
	delete(l.books, bookID)
}

func (l *Library) BorrowBook(bookID int, memberID int) error {

	book, ok := l.books[bookID]
	if !ok {
		return errors.New("book not found")
	}

	if book.Status == "Borrowed" {
		return errors.New("book is already borrowed")
	}

	member, ok := l.members[memberID]
	if !ok {
		return errors.New("member not found")
	}

	book.Status = "Borrowed"
	l.books[bookID] = book
	member.BorrowedBooks = append(member.BorrowedBooks, book)
	l.members[memberID] = member
	return nil
}

func (l *Library) ReturnBook(bookID int, memberID int) error {
	book, ok := l.books[bookID]
	if !ok {
		return errors.New("book not found")
	}

	member, ok := l.members[memberID]
	if !ok {
		return errors.New("member not found")
	}

	found := false
	newBorrowedBooks := []models.Book{}
	for _, borrowedBook := range member.BorrowedBooks {
		if borrowedBook.ID == bookID {
			found = true
		} else {
			newBorrowedBooks = append(newBorrowedBooks, borrowedBook)
		}
	}

	if !found {
		return errors.New("this member did not borrow this book")
	}
	member.BorrowedBooks = newBorrowedBooks
	l.members[memberID] = member
	book.Status = "Available"
	l.books[bookID] = book

	return nil // Success
}

func (l *Library) ListAvailableBooks() []models.Book {
	availableBooks := []models.Book{}
	for _, book := range l.books {
		if book.Status == "Available" {
			availableBooks = append(availableBooks, book)
		}
	}
	return availableBooks
}

func (l *Library) ListBorrowedBooks(memberID int) ([]models.Book, error) {
	member, ok := l.members[memberID]
	if !ok {
		return nil, errors.New("member not found")
	}
	return member.BorrowedBooks, nil
}

func (l *Library) AddMember(member models.Member) {
	l.members[member.ID] = member
}
