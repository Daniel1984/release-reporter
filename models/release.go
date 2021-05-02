package models

import "time"

// Release struct will hold generic release info
type Release struct {
	ID              uint64    `json:"id"`
	TagName         string    `json:"tag_name"`
	TargetCommitish string    `json:"target_commitish"`
	Name            string    `json:"name"`
	IsDraft         bool      `json:"draft"`
	IsPrerelease    bool      `json:"prerelease"`
	CreatedAt       time.Time `json:"created_at"`
	PublishedAt     time.Time `json:"published_at"`
}

// Releases will hold a list of all releases for given repo
type Releases []Release
