package internals

import "golang.org/x/tools/go/vcs"

type Repository struct {
	URL string
}

type RepositorySaver interface {
	Download(path, dest string) (root *vcs.RepoRoot, err error)
}
