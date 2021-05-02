## release-reporter

### Notes
Tool used to track public git project releases and report relevant info to
slack channel over webhook.

### Requirements
1. slack channel and webhook with permission to post messages. More info [here](https://slack.com/intl/en-lt/help/articles/115005265063-Incoming-webhooks-for-Slack)
2. to avoid rate limiting, github access token has to be provided. More info
   [here](https://docs.github.com/en/github/authenticating-to-github/creating-a-personal-access-token). Only one scope is required: `repo > public_repo`
3. `CheckIntervalMs` from `config.yml` represents the check frequency in milliseconds. Currently
   hardcoded to once per 12h. Runs the check on initial launch.

### Running program
- `CGO_ENABLED=0 GOOS=linux go build -o releasereporter main.go`
- `nohup releasereporter &`
