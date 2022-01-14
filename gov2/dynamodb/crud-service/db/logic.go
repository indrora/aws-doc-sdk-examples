package db

type DbConnection interface {
	// Create a link
	Add(link Link)
	// List recent links
	ListRecent() []Link
	// Get a link by its ID, nil if there is no link by that id
	Get(id string) *Link
	// Destroy a link by its ID
	Delete(id string) bool
	// increment the view count on a link by its ID
	Increment(id string)
}
