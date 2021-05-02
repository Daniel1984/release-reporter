package gitscan

import (
	"encoding/json"
	"fmt"
	"strings"

	badger "github.com/dgraph-io/badger/v3"
	"github.com/pkg/errors"
	"github.com/release-reporter/api/github"
	"github.com/release-reporter/logger"
	"github.com/release-reporter/models"
)

func ScanReleases(cfg *models.Config, db *badger.DB, log *logger.Logger) string {
	var sb strings.Builder

	for _, project := range cfg.GitProjects {
		releases, err := github.GetReleasesFor(project.Owner, project.Repo, cfg.GitAuthToken)
		if err != nil {
			log.Info.Printf("failed getting releases for owner:%s | repo:%s | err:%s\n", project.Owner, project.Repo, err)
			continue
		}

		if len(releases) == 0 {
			log.Info.Printf("repo: %s has no releases", project.Repo)
			continue
		}

		// we only care about last/most recent release, which is 1st item in array
		newRelease := releases[0]

		// construct project specific uuid
		uuid := fmt.Sprintf("%s-%s", project.Owner, project.Repo)

		err = db.View(func(txn *badger.Txn) error {
			dbRecord, err := txn.Get([]byte(uuid))
			if err != nil {
				// Record for given uuid not found. Perform insert and return
				newReleaseBytes, err := json.Marshal(newRelease)
				if err != nil {
					return errors.Wrap(err, "failed marshaling newRelease")
				}

				return db.Update(func(txn *badger.Txn) error {
					return txn.Set([]byte(uuid), newReleaseBytes)
				})
			}

			// record for given uuid found, performing comparison/update/notification
			knownReleaseBytes, err := dbRecord.ValueCopy(nil)
			if err != nil {
				return errors.Wrap(err, "failed obtaining record value")
			}

			knownRelease := models.Release{}
			if err := json.Unmarshal(knownReleaseBytes, &knownRelease); err != nil {
				return errors.Wrap(err, "failed unmarshaling record bytes to release object")
			}

			// new release detected
			if knownRelease.ID != newRelease.ID {
				newReleaseBytes, err := json.Marshal(newRelease)
				if err != nil {
					return errors.Wrap(err, "failed marshaling newRelease")
				}

				if err := db.Update(func(txn *badger.Txn) error {
					return txn.Set([]byte(uuid), newReleaseBytes)
				}); err != nil {
					return errors.Wrap(err, "failed updating record")
				}

				sb.WriteString(fmt.Sprintf("*Name:* %s/%s\n", project.Owner, project.Repo))
				sb.WriteString(fmt.Sprintf("*Release ID:* %d\n", newRelease.ID))
				sb.WriteString(fmt.Sprintf("*TagName:* %s\n", newRelease.TagName))
				sb.WriteString(fmt.Sprintf("*TargetCommitish:* %s\n", newRelease.TargetCommitish))
				sb.WriteString(fmt.Sprintf("*Name:* %s\n", newRelease.Name))
				sb.WriteString(fmt.Sprintf("*IsDraft:* %t\n", newRelease.IsDraft))
				sb.WriteString(fmt.Sprintf("*IsPrerelease:* %t\n", newRelease.IsPrerelease))
				sb.WriteString(fmt.Sprintf("*CreatedAt:* %s\n", newRelease.CreatedAt.Format("2006-01-02T15:04:05")))
				sb.WriteString(fmt.Sprintf("*PublishedAt:* %s\n", newRelease.PublishedAt.Format("2006-01-02T15:04:05")))
				sb.WriteString("----------------------------------------------")
			}

			return nil
		})
		if err != nil {
			log.Err.Printf("%s\n", err)
		}
	}

	return sb.String()
}
