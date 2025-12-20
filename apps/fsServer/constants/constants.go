package constants

type SortField string

const (
	ModifiedDate SortField = "modifiedDate"
	Name         SortField = "name"
)

type SortOrder string

const (
	Ascending  SortOrder = "asc"
	Descending SortOrder = "desc"
)
