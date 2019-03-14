// Package git provides models of Git objects.
package git

// RepositoryID represents a pointer to a repository.
type RepositoryID struct {
	Owner string
	Name  string
}

// IsValid returns true if owner and name is not empty.
func (id RepositoryID) IsValid() bool {
	return id.Owner != "" && id.Name != ""
}

func (id RepositoryID) String() string {
	if !id.IsValid() {
		return ""
	}
	return id.Owner + "/" + id.Name
}

// BranchName represents name of a branch.
type BranchName string

// QualifiedName returns RefQualifiedName.
// If the BranchName is empty, it returns a zero value.
func (b BranchName) QualifiedName() RefQualifiedName {
	if b == "" {
		return RefQualifiedName{}
	}
	return RefQualifiedName{"refs/heads/", string(b)}
}

// RefName represents name of a ref, that is a branch or a tag.
// This may be simple name or qualified name.
type RefName string

// RefQualifiedName represents qualified name of a ref, e.g. refs/heads/master.
type RefQualifiedName struct {
	Prefix string
	Name   string
}

func (r RefQualifiedName) IsValid() bool {
	return r.Prefix != "" && r.Name != ""
}

func (r RefQualifiedName) String() string {
	if !r.IsValid() {
		return ""
	}
	return r.Prefix + r.Name
}

// NewBranch represents a branch.
type NewBranch struct {
	Repository RepositoryID
	BranchName BranchName
	CommitSHA  CommitSHA
}

// CommitSHA represents a pointer to a commit.
type CommitSHA string

// CommitMessage represents a message of a commit.
type CommitMessage string

// NewCommit represents a commit.
type NewCommit struct {
	Repository      RepositoryID
	Message         CommitMessage
	ParentCommitSHA CommitSHA
	TreeSHA         TreeSHA
}

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
