package domain


type Bookmark struct {
	UserID string `json:"user_id"`
	HackathonID string `json:"hackathon_id"`
	CreatedAt string `json:"created_at"`
}

type BookmarkRepository interface {
	AddBookmark(userID, hackathonID string) error
	RemoveBookmark(userID, hackathonID string) error
	GetBookmarksByUser(userID string) ([]Bookmark, error)
}

