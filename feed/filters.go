package feed

type FilterRule struct {
	url  string
	expr string
}

type FilterFile struct {
	filePath   string
	rawContent []byte
}
