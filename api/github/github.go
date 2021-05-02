package github

import (
	"fmt"

	"github.com/release-reporter/models"
	"github.com/release-reporter/request"
)

func GetReleasesFor(owner, repo, token string) (r models.Releases, err error) {
	req := request.
		New("GET", fmt.Sprintf("https://api.github.com/repos/%s/%s/releases", owner, repo), nil).
		AddHeaders("Accept", "application/vnd.github.v3+json").
		AddHeaders("Authorization", fmt.Sprintf("token %s", token)).
		Do().
		Decode(&r)

	if err = req.HasError(); err != nil {
		return r, fmt.Errorf("failed to fetch release for owner:%s, repo%s , err: %s", owner, repo, err)
	}

	return
}
