package pkg

// CType type MEMORY or DATABASE
type CType string

const (
	// MEMORY type for in memory based on sqlite
	MEMORY CType = "MEMORY"
	// DATABASE type based on different types of db
	DATABASE CType = "DATABASE"

	// EMPTYSTRING string
	EMPTYSTRING = ""
)
