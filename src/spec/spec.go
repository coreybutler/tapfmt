package spec

import (
	"fmt"
	"regexp"
	"strings"

	color "github.com/logrusorgru/aurora/v3"
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
	suite := ""

	if len(results.Explanation) > 1 && suite != results.Explanation[0] {
		suite = results.Explanation[0]
		fmt.Printf("\n  %s\n\n", color.Bold(suite))

		suite = results.Explanation[len(results.Explanation)-1]
		fmt.Printf("  %s\n\n", color.Underline(suite))
	} else if len(results.Explanation) == 1 {
		suite = results.Explanation[0]
		fmt.Printf("\n  %s\n\n", color.Underline(suite))
	}

	suitechanged := false

	failedtests := make([]tap.Test, 0)
	for i, test := range results.Tests {
		if test.Skipped {
			re := regexp.MustCompile(`(?i)skip`)
			icon := "\u21B7"
			// icon := "\u2A2F"
			fmt.Printf("    %v %v %s\n", color.Faint(color.Blue(icon)), color.Faint(color.Blue("skipped")), color.Faint(color.Blue(strings.TrimSpace(re.ReplaceAllString(test.DirectiveText, "")))))
		} else if test.Todo {
			re := regexp.MustCompile(`(?i)todo`)
			icon := strings.TrimSpace("🗹")
			fmt.Printf("    %v%v %s\n", color.Yellow(icon), color.Bold(color.Yellow("TO DO:")), color.Yellow(strings.TrimSpace(re.ReplaceAllString(test.DirectiveText, ""))))
		} else {
			if test.Passed {
				fmt.Printf("    %v %v\n", color.Green("\u2713"), color.Faint(test.Description))
			} else {
				failedtests = append(failedtests, test)

				begin := ""
				if len(test.YamlBytes) > 0 && !suitechanged && i > 0 {
					begin = "\n"
				}

				fail(test, false, begin)
			}

			detail(test)
		}

		suitechanged = false

		if len(test.Diagnostics) > 0 && test.TestNumber < results.TotalTests {
			suite = test.Diagnostics[0]
			suitechanged = true
			fmt.Printf("\n  %s\n\n", color.Underline(suite))
		}
	}

	if results.BailOut {
		sep := ""
		if len(strings.TrimSpace(results.BailOutReason)) > 0 {
			sep = ": "
		}
		fmt.Printf("\n  %s%s%s\n", color.Bold(color.Yellow("\u26A0 Aborted")), color.Yellow(sep), color.Yellow(results.BailOutReason))
	}

	f.FailedTests = failedtests
}

func (f *Format) Summary() {
	results := f.Results
	failedtests := f.FailedTests

	if results.TotalTests == 0 {
		fmt.Println("  No tests found\n")
	} else {
		if len(failedtests) > 0 {
			tense := "were"
			action := "failures"
			if len(failedtests) == 1 {
				tense = "was"
				action = "failure"
			}

			fmt.Printf("\n  %s There %s %s %s\n\n", color.Bold(color.BrightRed("Failed Tests:")), tense, color.Bold(color.BrightRed(fmt.Sprintf("%v", len(failedtests)))), action)

			for _, test := range failedtests {
				fail(test, false)
			}
			fmt.Println("")
		}

		fmt.Printf("  total:     %v\n", results.TotalTests)
		fmt.Printf("  %s   %s\n", color.Green("passing:"), fmt.Sprintf("%v", color.Green(results.PassedTests)))
		fmt.Printf("  %s   %s\n", color.Red("failing:"), fmt.Sprintf("%v", color.Red(results.FailedTests)))
		if results.SkippedTests > 0 {
			fmt.Printf("  %s   %s\n", color.Blue("skipped:"), fmt.Sprintf("%v", color.Blue(results.SkippedTests)))
		}
		if results.TodoTests > 0 {
			fmt.Printf("  %s     %s\n", color.Yellow("tasks:"), fmt.Sprintf("%v", color.Yellow(results.TodoTests)))
		}
		fmt.Println("\n")
	}
}

func fail(test tap.Test, info bool, prefix ...string) {
	begin := ""
	if len(prefix) > 0 {
		begin = prefix[0]
	}

	fmt.Printf("%s    %v %v\n", begin, color.Red("\u2A2F"), color.Faint(color.Red(test.Description)))

	if info {
		detail(test)
	}
}

func detail(test tap.Test) {
	if len(test.YamlBytes) > 0 {
		chars := ""
		for i := 0; i < (len(test.Description) + 2); i++ {
			chars = chars + "-"
		}

		fmt.Printf("    %s\n", color.Faint(color.Red(chars)))

		for _, line := range strings.Split(string(test.YamlBytes), "\n") {
			fmt.Printf("  %s\n", color.Cyan(line))
		}
	}
}
