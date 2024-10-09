package dto

// Post : you can use `type Post map[string]any` if there is no need to define struct fields
type Post struct {
	ID     int    `json:"id"`
	UserID int    `json:"userId"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

type Posts struct {
	Posts []Post `json:"posts"`
}
