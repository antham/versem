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

func TestLabelCreate(t *testing.T) {
	type scenario struct {
		name             string
		getSemverService func() semverService
		test             func(exitCode int, stdout bytes.Buffer, stderr bytes.Buffer, methodCallCount map[string]int)
	}

	scenarios := []scenario{
		{
			"Failure occurred when creating labels",
			func() semverService {
				return semverServiceMock{err: fmt.Errorf("failure occurred when creating labels"), methodCallCount: map[string]int{}}
			},
			func(exitCode int, stdout bytes.Buffer, stderr bytes.Buffer, methodCallCount map[string]int) {
				assert.Equal(t, 1, exitCode)
				assert.Equal(t, "failure occurred when creating labels\n", stderr.String())
				assert.Len(t, methodCallCount, 1)
				assert.Equal(t, 1, methodCallCount["CreateList"])
			},
		},
		{
			"Create semver labels on repository",
			func() semverService {
				return semverServiceMock{methodCallCount: map[string]int{}, version: github.MINOR}
			},
			func(exitCode int, stdout bytes.Buffer, stderr bytes.Buffer, methodCallCount map[string]int) {
				assert.Equal(t, 0, exitCode)
				assert.Empty(t, stderr.String())
				assert.Equal(t, "semver labels created\n", stdout.String())
				assert.Len(t, methodCallCount, 1)
				assert.Equal(t, 1, methodCallCount["CreateList"])
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

				labelCreate(msgHandler, semverService, &cobra.Command{}, []string{})
			}()

			w.Wait()

			scenario.test(code, stdout, stderr, semverService.(semverServiceMock).methodCallCount)
		})
	}

}
