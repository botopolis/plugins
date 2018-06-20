package buildkite

import (
	"net/http"

	"github.com/botopolis/bot"
)

const (
	tokenHeader = "X-Buildkite-Token"
	eventHeader = "X-Buildkite-Event"
)

type Plugin struct {
	bot            *bot.Robot
	buildHooks     []func(BuildEvent)
	BuildkiteToken string
}

func New(token string) *Plugin {
	return &Plugin{
		BuildkiteToken: token,
	}
}

func (p Plugin) OnBuild(fn func(BuildEvent)) {
	p.buildHooks = append(p.buildHooks, fn)
}

func (p *Plugin) Load(r *bot.Robot) {
	p.bot = r
	r.Router.HandleFunc("/webhooks/buildkite", p.handler).Methods("POST")
}

func (p *Plugin) handler(w http.ResponseWriter, r *http.Request) {
	if token := r.Header.Get(tokenHeader); p.BuildkiteToken != token {
		p.bot.Logger.Error("Received a POST to webhooks/buildkite with an invalid token")
		return
	}

	switch Event(r.Header.Get(eventHeader)) {
	case BuildRunning, BuildScheduled, BuildFinished:
		buildEvent, err := parseBuildEvent(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		for _, fn := range p.buildHooks {
			go fn(buildEvent)
		}
	}

	w.WriteHeader(http.StatusOK)
}
