package git

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
