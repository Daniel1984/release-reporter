package models

import "time"

// GitProject used to map project related data from config.yml
type GitProject struct {
	Owner string
	Repo  string
}

// Config will be used for mapping config.yml content to for easy access
type Config struct {
	Webhook         string
	DBPath          string
	ResponseType    string
	GitAuthToken    string
	GitProjects     []GitProject
	CheckIntervalMs time.Duration
}
