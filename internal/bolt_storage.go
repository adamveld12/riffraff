package internal

type Bookmark struct {
	ID   string
	URL  string
	Name string
	Tags []string
}

type User struct {
	ID           uint64
	Username     string
	PasswordHash string
}

type Storage interface {
	AddUser(User) (uint64, error)
	AddBookmark(string, string, []string) error
	GetBookmarks(string) ([]Bookmark, error)
	DeleteBookmark(string) error
}
