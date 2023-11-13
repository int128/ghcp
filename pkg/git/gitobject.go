package git

// CommitSHA represents a pointer to a commit.
type CommitSHA string

// CommitMessage represents a message of a commit.
type CommitMessage string

// NewCommit represents a commit.
type NewCommit struct {
	Repository      RepositoryID
	Message         CommitMessage
	Author          *CommitAuthor // optional
	Committer       *CommitAuthor // optional
	ParentCommitSHA CommitSHA     // optional
	TreeSHA         TreeSHA
}

// CommitAuthor represents an author of commit.
type CommitAuthor struct {
	Name  string
	Email string
}

// TreeSHA represents a pointer to a tree.
type TreeSHA string

// File represents a file in a tree.
type File struct {
	Filename   string  // filename (including path separators)
	BlobSHA    BlobSHA // blob SHA
	Executable bool    // if the file is executable
	Deleted    bool    // if the file is deleted
}

// Mode returns mode of the file, i.e. 100644 or 100755.
func (f *File) Mode() string {
	if f.Executable {
		return "100755"
	}
	return "100644"
}

// NewTree represents a tree.
type NewTree struct {
	Repository  RepositoryID
	BaseTreeSHA TreeSHA
	Files       []File
}

// BlobSHA represents a pointer to a blob.
type BlobSHA string

// NewBlob represents a blob.
type NewBlob struct {
	Repository RepositoryID
	Content    string // base64 encoded content
}
