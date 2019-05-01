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

type semverServiceMock struct {
	version         github.Version
	err             error
	methodCallCount map[string]int
}

func (s semverServiceMock) GetFromPullRequest(int) (github.Version, error) {
	s.methodCallCount["GetFromPullRequest"]++
	return s.version, s.err
}

func (s semverServiceMock) GetFromCommit(string) (github.Version, error) {
	s.methodCallCount["GetFromCommit"]++
	return s.version, s.err
}

func (s semverServiceMock) CreateList() error {
	s.methodCallCount["CreateList"]++
	return s.err
}

func TestCheck(t *testing.T) {
	type scenario struct {
		name             string
		getSemverService func() semverService
		getArgument      func() []string
		test             func(exitCode int, stdout bytes.Buffer, stderr bytes.Buffer, methodCallCount map[string]int)
	}

	scenarios := []scenario{
		{
			"No argument provided",
			func() semverService {
				return semverServiceMock{}
			},
			func() []string {
				return []string{}
			},
			func(exitCode int, stdout bytes.Buffer, stderr bytes.Buffer, methodCallCount map[string]int) {
				assert.Equal(t, 1, exitCode)
				assert.Equal(t, "provide a pull request number or commit sha as first argument\n", stderr.String())
				assert.Len(t, methodCallCount, 0)
			},
		},
		{
			"Argument provided is not supported",
			func() semverService {
				return semverServiceMock{}
			},
			func() []string {
				return []string{"test"}
			},
			func(exitCode int, stdout bytes.Buffer, stderr bytes.Buffer, methodCallCount map[string]int) {
				assert.Equal(t, 1, exitCode)
				assert.Equal(t, "analysis failed, test is not a valid number, nor a valid commit sha\n", stderr.String())
				assert.Len(t, methodCallCount, 0)
			},
		},
		{
			"Failure occurred when fetching label from commit",
			func() semverService {
				return semverServiceMock{err: fmt.Errorf("failure occurred when calling fetching label from commit"), methodCallCount: map[string]int{}}
			},
			func() []string {
				return []string{"a1b2c3d4"}
			},
			func(exitCode int, stdout bytes.Buffer, stderr bytes.Buffer, methodCallCount map[string]int) {
				assert.Equal(t, 1, exitCode)
				assert.Equal(t, "analysis failed, failure occurred when calling fetching label from commit\n", stderr.String())
				assert.Len(t, methodCallCount, 1)
				assert.Equal(t, 1, methodCallCount["GetFromCommit"])
			},
		},
		{
			"Failure occurred when fetching label from pull request",
			func() semverService {
				return semverServiceMock{err: fmt.Errorf("failure occurred when calling fetching label from pull request"), methodCallCount: map[string]int{}}
			},
			func() []string {
				return []string{"123"}
			},
			func(exitCode int, stdout bytes.Buffer, stderr bytes.Buffer, methodCallCount map[string]int) {
				assert.Equal(t, 1, exitCode)
				assert.Equal(t, "analysis failed, failure occurred when calling fetching label from pull request\n", stderr.String())
				assert.Len(t, methodCallCount, 1)
				assert.Equal(t, 1, methodCallCount["GetFromPullRequest"])
			},
		},
		{
			"Fetch version from pull request label",
			func() semverService {
				return semverServiceMock{methodCallCount: map[string]int{}, version: github.MINOR}
			},
			func() []string {
				return []string{"123"}
			},
			func(exitCode int, stdout bytes.Buffer, stderr bytes.Buffer, methodCallCount map[string]int) {
				assert.Equal(t, 0, exitCode)
				assert.Empty(t, stderr.String())
				assert.Equal(t, "minor semver version found\n", stdout.String())
				assert.Len(t, methodCallCount, 1)
				assert.Equal(t, 1, methodCallCount["GetFromPullRequest"])
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

			go func() {
				defer func() {
					if r := recover(); r != nil {
						code = r.(int)
					}

					w.Done()
				}()

				check(msgHandler, semverService, &cobra.Command{}, scenario.getArgument())
			}()

			w.Wait()

			scenario.test(code, stdout, stderr, semverService.(semverServiceMock).methodCallCount)
		})
	}

}
