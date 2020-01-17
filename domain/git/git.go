// Package git provides the models of Git such as a repository and branch.
package git

// InternalNodeID represents a node ID for GitHub v4 API.
// This is allocated by GitHub.
type InternalNodeID interface{}

// RepositoryID represents a pointer to a repository.
type RepositoryID struct {
	Owner          string
	Name           string
	InternalNodeID InternalNodeID
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

// TagName represents name of a tag.
type TagName string

// Name returns the name of tag.
func (t TagName) Name() string {
	return string(t)
}

// QualifiedName returns RefQualifiedName.
// If the TagName is empty, it returns a zero value.
func (t TagName) QualifiedName() RefQualifiedName {
	if t == "" {
		return RefQualifiedName{}
	}
	return RefQualifiedName{"refs/tags/", string(t)}
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

// ReleaseID represents an ID of release.
type ReleaseID struct {
	Repository RepositoryID
	InternalID int64 // GitHub API will allocate this ID
}

// Release represents a release associated to a tag.
type Release struct {
	ID      ReleaseID
	TagName TagName
	Name    string
}

// ReleaseAsset represents a release asset.
type ReleaseAsset struct {
	Release  ReleaseID
	Name     string
	RealPath string
}
