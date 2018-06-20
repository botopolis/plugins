package buildkite

import (
	"encoding/json"
	"io"
	"time"
)

type Event string

const (
	BuildScheduled Event = "build.scheduled"
	BuildRunning   Event = "build.running"
	BuildFinished  Event = "build.finished"
)

type State string

const (
	Running   State = "running"
	Scheduled State = "scheduled"
	Passed    State = "passed"
	Failed    State = "failed"
	Blocked   State = "blocked"
	Canceled  State = "canceled"
	Canceling State = "canceling"
	Skipped   State = "skipped"
	NotRun    State = "not_run"
	Finished  State = "finished"
)

// Job corresponds to the buildkite resource:
// https://buildkite.com/docs/rest-api/jobs
type Job struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	Name    string `json:"name"`
	Command string `json:"command"`

	WebURL    string `json:"web_url"`
	LogURL    string `json:"log_url"`
	RawLogURL string `json:"raw_log_url"`

	State      State `json:"state"`
	ExitStatus int64 `json:"exit_status"`

	CreatedAt   time.Time `json:"created_at,string"`
	ScheduledAt time.Time `json:"scheduled_at,string"`
	StartedAt   time.Time `json:"started_at,string"`
	FinishedAt  time.Time `json:"finished_at,string"`
}

// Build corresponds to the buildkite resource:
// https://buildkite.com/docs/rest-api/builds
type Build struct {
	ID      string `json:"id"`
	Message string `json:"message"`
	Branch  string `json:"branch"`
	Commit  string `json:"commit"`

	URL    string `json:"url"`
	WebURL string `json:"web_url"`

	Creator struct {
		AvatarURL string `json:"avatar_url"`
		CreatedAt string `json:"created_at"`
		Email     string `json:"email"`
		ID        string `json:"id"`
		Name      string `json:"name"`
	} `json:"creator"`

	Jobs []Job `json:"jobs"`

	Env      map[string]string `json:"env"`
	MetaData map[string]string `json:"meta_data"`

	State    State    `json:"state"`
	Blocked  bool     `json:"blocked"`
	Number   int64    `json:"number"`
	Pipeline Pipeline `json:"pipeline"`
	Source   string   `json:"source"`

	CreatedAt   time.Time `json:"created_at,string"`
	ScheduledAt time.Time `json:"scheduled_at,string"`
	StartedAt   time.Time `json:"started_at,string"`
	FinishedAt  time.Time `json:"finished_at,string"`
}

// Pipeline corresponds to the buildkite resource:
// https://buildkite.com/docs/rest-api/pipelines
type Pipeline struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description"`

	BuildsURL string `json:"builds_url"`
	BadgeURL  string `json:"badge_url"`
	URL       string `json:"url"`
	WebURL    string `json:"web_url"`

	Env map[string]string

	Repository          string  `json:"repository"`
	DefaultBranch       string  `json:"default_branch"`
	BranchConfiguration *string `json:"branch_configuration"`

	SkipQueuedBranchBuilds          bool    `json:"skip_queued_branch_builds"`
	SkipQueuedBranchBuildsFilter    *string `json:"skip_queued_branch_builds_filter"`
	CancelRunningBranchBuilds       bool    `json:"cancel_running_branch_builds"`
	CancelRunningBranchBuildsFilter *string `json:"cancel_running_branch_builds_filter"`

	RunningBuildsCount   int64 `json:"running_builds_count"`
	RunningJobsCount     int64 `json:"running_jobs_count"`
	ScheduledBuildsCount int64 `json:"scheduled_builds_count"`
	ScheduledJobsCount   int64 `json:"scheduled_jobs_count"`
	WaitingJobsCount     int64 `json:"waiting_jobs_count"`

	CreatedAt time.Time `json:"created_at,string"`
}

// BuildEvent is the payload received from Buildkite for build events
type BuildEvent struct {
	Event    Event    `json:"event"`
	Build    Build    `json:"build"`
	Pipeline Pipeline `json:"pipeline"`
	Sender   struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"sender"`
}

func parseBuildEvent(r io.Reader) (be BuildEvent, err error) {
	decoder := json.NewDecoder(r)
	err = decoder.Decode(&be)
	return be, err
}
