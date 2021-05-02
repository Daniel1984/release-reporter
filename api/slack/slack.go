package slack

import (
	"bytes"
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/release-reporter/request"
)

// SendMessage takes message and url as string arguments and calls slack webhook with valid payload
func SendMessage(msg, url string) (err error) {
	var slackPld = struct {
		Text string `json:"text"`
	}{
		Text: msg,
	}

	jb, err := json.Marshal(slackPld)
	if err != nil {
		return errors.Wrap(err, "failed marshaling slack msg payload")
	}

	req := request.
		New("POST", url, bytes.NewReader(jb)).
		AddHeaders("Content-type", "application/json").
		Do()

	if err = req.HasError(); err != nil {
		return errors.Wrap(err, "failed posting slack message")
	}

	return
}
