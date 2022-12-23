package json

import (
	base "encoding/json"
	"fmt"
	"strings"

	tap "github.com/mpontillo/tap13"
)

type Format struct {
	Results     *tap.Results
	FailedTests []tap.Test
}

func Formatter() *Format {
	return &Format{}
}

func (f *Format) Format(results *tap.Results) {
	f.Results = results
	// data := map[string]interface{}{}

	r := Results{
		Version: results.TapVersion,
		Summary: ResultSummary{
			Total:         results.TotalTests,
			Passed:        results.PassedTests,
			Failed:        results.FailedTests,
			Skipped:       results.SkippedTests,
			Todo:          results.TodoTests,
			Expected:      results.ExpectedTests,
			BailOut:       results.BailOut,
			BailOutReason: results.BailOutReason,
		},
		Tests: make([]TestResult, 0),
	}

	var suite string
	var group string

	if len(results.Explanation) > 1 && suite != results.Explanation[0] {
		suite = results.Explanation[0]
		group = results.Explanation[len(results.Explanation)-1]
	} else if len(results.Explanation) == 1 {
		group = results.Explanation[0]
	}

	for _, test := range results.Tests {
		t := TestResult{
			Suite:       suite,
			Group:       group,
			Number:      test.TestNumber,
			Passed:      test.Passed,
			Description: test.Description,
		}

		if len(t.Description) == 0 && len(test.DirectiveText) > 0 {
			items := strings.Split(test.DirectiveText, " ")
			t.Directive = items[0]
			t.Description = strings.Join(items[1:], " ")
		}

		if len(test.YamlBytes) > 0 {
			t.Metadata = strings.TrimSpace(string(test.YamlBytes))
		}

		r.Tests = append(r.Tests, t)

		if !t.Passed && t.Directive == "" {
			if r.Summary.Failures == nil {
				r.Summary.Failures = make([]TestResult, 0)
			}

			r.Summary.Failures = append(r.Summary.Failures, t)
		}

		if len(test.Diagnostics) > 0 && test.TestNumber < results.TotalTests {
			group = test.Diagnostics[0]
		}
	}

	data, err := base.MarshalIndent(r, "", "  ")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(data))
}

func (f *Format) Summary() {}

type Results struct {
	Version int           `json:"version"`
	Summary ResultSummary `json:"summary"`
	Tests   []TestResult  `json:"results"`
}

type ResultSummary struct {
	Total         int          `json:"total"`
	Passed        int          `json:"passed"`
	Failed        int          `json:"failed"`
	Skipped       int          `json:"skipped"`
	Todo          int          `json:"todo"`
	Expected      int          `json:"expected"`
	BailOut       bool         `json:"bailout,omitempty"`
	BailOutReason string       `json:"bailout_reason,omitempty"`
	Failures      []TestResult `json:"failures,omitempty"`
}

type TestResult struct {
	Suite       string `json:"suite,omitempty"`
	Group       string `json:"group,omitempty"`
	Number      int    `json:"test_number"`
	Passed      bool   `json:"passed"`
	Directive   string `json:"directive,omitempty"`
	Description string `json:"description"`
	Metadata    string `json:"info,omitempty"`
}
