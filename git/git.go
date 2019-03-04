// Package git provides models of Git objects.
package git

// RepositoryID represents a pointer to a repository.
type RepositoryID struct {
	Owner string
	Name  string
}

// BranchName represents name of a branch.
type BranchName string

// CommitSHA represents a pointer to a commit.
type CommitSHA string

// CommitMessage represents a message of a commit.
type CommitMessage string

// TreeSHA represents a pointer to a tree.
type TreeSHA string

// File represents a file in a tree.
type File struct {
	Filename   string  // filename (including path separators)
	BlobSHA    BlobSHA // blob SHA
	Executable bool    // if the file is executable
}

// Mode returns mode of the file, i.e. 100644 or 100755.
func (f *File) Mode() string {
	if f.Executable {
		return "100755"
	}
	return "100644"
}

// BlobSHA represents a pointer to a blob.
type BlobSHA string
