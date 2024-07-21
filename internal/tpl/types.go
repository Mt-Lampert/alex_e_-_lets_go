package tpl

type DbSnippet struct {
	ID      string
	Title   string
	Content string
	Created string
	Expires string
}

// a convenient wrapper for template data (usually from the DB)
type TemplateData struct {
	CSRFToken       string
	CurrentYear     string
	Flash           string
	Form            any
	Snippet         DbSnippet
	Snippets        []DbSnippet
	IsAuthenticated bool
	URL             string
}

// vim: ts=4 sw=4 fdm=indent
