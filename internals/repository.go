package internals

type Repository struct {
	URL string
}

type RepositorySaver interface {
	Download(path, dest string) (fullLocalPath string, err error)
}
