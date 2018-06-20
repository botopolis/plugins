package buildkite_test

import (
	"fmt"
	"os"
	"sync"

	"github.com/botopolis/bot"
	"github.com/botopolis/plugins/buildkite"
	"github.com/botopolis/slack"
	oslack "github.com/nlopes/slack"
)

type skipped struct {
	emailsPerPipeline map[string][]string
	mu                sync.Mutex
}

func newSkipped() *skipped {
	return &skipped{
		emailsPerPipeline: map[string][]string{},
	}
}
func (s *skipped) Add(pipeline, email string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.emailsPerPipeline[pipeline] = append(s.emailsPerPipeline[pipeline], email)
}

func (s *skipped) Pop(pipeline string) (emails []string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	emails = s.emailsPerPipeline[pipeline]
	delete(s.emailsPerPipeline, pipeline)
	return emails
}

func Example() {
	chat := slack.New(os.Getenv("SLACK_TOKEN"))

	skpd := newSkipped()

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

		if be.Build.State == buildkite.Skipped {
			skpd.Add(be.Pipeline.ID, be.Build.Creator.Email)
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
				Fallback:   fmt.Sprintf("[%s] Build %s", be.Pipeline.Name, be.Build.State),
				Title:      fmt.Sprintf("[%s] Build %s", be.Pipeline.Name, be.Build.State),
				TitleLink:  be.Build.WebURL,
				Text:       fmt.Sprintf("[%s]: %s", be.Build.Branch, be.Build.Commit),
				AuthorName: "buildkite",
				Fields:     fields,
			}},
		}

		emails := append(skpd.Pop(be.Pipeline.ID), be.Build.Creator.Email)
		for _, email := range emails {
			user, ok := chat.Store.UserByEmail(email)
			if !ok {
				continue
			}
			chat.Direct(bot.Message{
				User:   user.ID,
				Params: params,
			})
		}
	})

	robot := bot.New(
		chat,
		bk,
	)

	robot.Run()
}
