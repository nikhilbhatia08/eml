package parser

var (
	Keywords = map[string]bool{
		"div":     true,
		"h1":     true,
	}

	config = map[string]bool {
		"styles": true,
		"content": true,
	}
)

var (
	KEYWORD_TYPE      = 1
	CONFIG_TYPE       = 2
	GENERAL_TYPE      = 3
)