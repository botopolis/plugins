package buildkite_test

import (
	"fmt"
	"os"

	"github.com/botopolis/bot"
	"github.com/botopolis/plugins/buildkite"
	"github.com/botopolis/slack"
	oslack "github.com/nlopes/slack"
)

func Example() {
	chat := slack.New(os.Getenv("SLACK_TOKEN"))

	stateColor := func(state buildkite.State) string {
		if state == buildkite.Failed || state == buildkite.Canceled {
			return "danger"
		}

		return "info"
	}

	bk := buildkite.New(os.Getenv("BUILDKITE_TOKEN"))
	bk.OnBuild(func(be buildkite.BuildEvent) {
		if be.Event != buildkite.BuildFinished {
			return
		}

		user, ok := chat.Store.UserByEmail(be.Build.Creator.Email)
		if !ok {
			return
		}

		var fields []oslack.AttachmentField
		for _, job := range be.Build.Jobs {
			fields = append(fields, oslack.AttachmentField{
				Title: job.Name,
				Value: string(job.State),
			})
		}
		params := oslack.PostMessageParameters{
			Attachments: []oslack.Attachment{{
				Color:      stateColor(be.Build.State),
				Fallback:   fmt.Sprintf("Build %s", be.Build.State),
				Title:      fmt.Sprintf("Build %s", be.Build.State),
				TitleLink:  be.Build.WebURL,
				Text:       fmt.Sprintf("[%s]: %s", be.Build.Branch, be.Build.Commit),
				AuthorName: "buildkite",
				Fields:     fields,
			}},
		}

		chat.Direct(bot.Message{
			User:   user.ID,
			Params: params,
		})
	})

	robot := bot.New(
		chat,
		bk,
	)

	robot.Run()
}
