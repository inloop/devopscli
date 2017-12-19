package cmd

import "testing"

func TestTagForRefName(t *testing.T) {

	cases := map[string]string{
		"master":       "latest",
		"develop":      "unstable",
		"release/xxx":  "xxx",
		"release-aaxx": "release-aaxx",
		"blah":         "blah",
	}

	for input, expectedOutput := range cases {
		output := tagForRefName(input)
		if expectedOutput != output {
			t.Fatalf("%s != %s", expectedOutput, output)
		}
	}
}
