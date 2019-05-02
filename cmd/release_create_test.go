package cmd

import (
	"bytes"
	"fmt"
	"sync"
	"testing"

	"github.com/antham/versem/github"

	"github.com/stretchr/testify/assert"

	"github.com/spf13/cobra"
)

func TestReleaseCreate(t *testing.T) {
	type scenario struct {
		name              string
		arguments         []string
		getSemverService  func() semverService
		getReleaseService func() releaseService
		test              func(exitCode int, stdout bytes.Buffer, stderr bytes.Buffer, semverServiceMethodCallCount map[string]int, releaseServiceMethodCallCount map[string]int)
	}

	scenarios := []scenario{
		{
			"No argument provided",
			[]string{},
			func() semverService {
				return semverServiceMock{methodCallCount: map[string]int{}}
			},
			func() releaseService {
				return releaseServiceMock{methodCallCount: map[string]int{}}
			},
			func(exitCode int, stdout bytes.Buffer, stderr bytes.Buffer, semverServiceMethodCallCount map[string]int, releaseServiceMethodCallCount map[string]int) {
				assert.Equal(t, 1, exitCode)
				assert.Equal(t, "provide a full commit sha as first argument\n", stderr.String())
				assert.Len(t, semverServiceMethodCallCount, 0)
				assert.Len(t, releaseServiceMethodCallCount, 0)
			},
		},
		{
			"Failure occurred when fetching label",
			[]string{"8a5ed8235d18fb0243493b82baf5d988459d24db"},
			func() semverService {
				return semverServiceMock{err: fmt.Errorf("failure occurred when fetching label"), methodCallCount: map[string]int{}}
			},
			func() releaseService {
				return releaseServiceMock{methodCallCount: map[string]int{}}
			},
			func(exitCode int, stdout bytes.Buffer, stderr bytes.Buffer, semverServiceMethodCallCount map[string]int, releaseServiceMethodCallCount map[string]int) {
				assert.Equal(t, 1, exitCode)
				assert.Equal(t, "failure occurred when fetching label\n", stderr.String())
				assert.Len(t, semverServiceMethodCallCount, 1)
				assert.Equal(t, 1, semverServiceMethodCallCount["GetFromCommit"])
				assert.Len(t, releaseServiceMethodCallCount, 0)
			},
		},
		{
			"Failure occurred when creating release",
			[]string{"8a5ed8235d18fb0243493b82baf5d988459d24db"},
			func() semverService {
				return semverServiceMock{methodCallCount: map[string]int{}}
			},
			func() releaseService {
				return releaseServiceMock{err: fmt.Errorf("failure occurred when creating release"), methodCallCount: map[string]int{}}
			},
			func(exitCode int, stdout bytes.Buffer, stderr bytes.Buffer, semverServiceMethodCallCount map[string]int, releaseServiceMethodCallCount map[string]int) {
				assert.Equal(t, 1, exitCode)
				assert.Equal(t, "failure occurred when creating release\n", stderr.String())
				assert.Len(t, semverServiceMethodCallCount, 1)
				assert.Equal(t, 1, semverServiceMethodCallCount["GetFromCommit"])
				assert.Len(t, releaseServiceMethodCallCount, 1)
				assert.Equal(t, 1, releaseServiceMethodCallCount["CreateNext"])
			},
		},
		{
			"Create release",
			[]string{"8a5ed8235d18fb0243493b82baf5d988459d24db"},
			func() semverService {
				return semverServiceMock{methodCallCount: map[string]int{}}
			},
			func() releaseService {
				return releaseServiceMock{methodCallCount: map[string]int{}, tag: github.Tag{Major: 1}}
			},
			func(exitCode int, stdout bytes.Buffer, stderr bytes.Buffer, semverServiceMethodCallCount map[string]int, releaseServiceMethodCallCount map[string]int) {
				assert.Equal(t, 0, exitCode)
				assert.Equal(t, "tag 1.0.0 created\n", stdout.String())
				assert.Len(t, semverServiceMethodCallCount, 1)
				assert.Equal(t, 1, semverServiceMethodCallCount["GetFromCommit"])
				assert.Len(t, releaseServiceMethodCallCount, 1)
				assert.Equal(t, 1, releaseServiceMethodCallCount["CreateNext"])
			},
		},
		{
			"Skip tag creation on norelease version",
			[]string{"8a5ed8235d18fb0243493b82baf5d988459d24db"},
			func() semverService {
				return semverServiceMock{methodCallCount: map[string]int{}, version: github.NORELEASE}
			},
			func() releaseService {
				return releaseServiceMock{methodCallCount: map[string]int{}, tag: github.Tag{Major: 1}}
			},
			func(exitCode int, stdout bytes.Buffer, stderr bytes.Buffer, semverServiceMethodCallCount map[string]int, releaseServiceMethodCallCount map[string]int) {
				assert.Equal(t, 0, exitCode)
				assert.Equal(t, "label norelease found, skip tag creation\n", stdout.String())
				assert.Len(t, semverServiceMethodCallCount, 1)
				assert.Equal(t, 1, semverServiceMethodCallCount["GetFromCommit"])
				assert.Len(t, releaseServiceMethodCallCount, 0)
			},
		},
	}

	for _, scenario := range scenarios {
		scenario := scenario
		t.Run(scenario.name, func(*testing.T) {

			var code int
			var stdout bytes.Buffer
			var stderr bytes.Buffer
			var w sync.WaitGroup

			msgHandler := messageHandler{
				func(exitCode int) {
					panic(exitCode)
				},
				&stdout,
				&stderr,
			}

			w.Add(1)

			semverService := scenario.getSemverService()
			releaseService := scenario.getReleaseService()

			go func() {
				defer func() {
					if r := recover(); r != nil {
						code = r.(int)
					}

					w.Done()
				}()

				releaseCreate(msgHandler, semverService, releaseService, &cobra.Command{}, scenario.arguments)
			}()

			w.Wait()

			scenario.test(code, stdout, stderr, semverService.(semverServiceMock).methodCallCount, releaseService.(releaseServiceMock).methodCallCount)
		})
	}

}
