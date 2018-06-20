package buildkite

import (
	"os"
	"testing"

	"github.com/samsarahq/go/snapshotter"
)

func TestParseBuildEvent(t *testing.T) {
	s := snapshotter.New(t)
	defer s.Verify()

	file, err := os.Open("testdata/build.json")
	if err != nil {
		t.Error(err)
	}

	buildEvent, err := parseBuildEvent(file)
	if err != nil {
		t.Error(err)
	}

	s.Snapshot("buildEvent", buildEvent)
}
