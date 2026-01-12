package main

import "testing"

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input string
		expected []string
	}{
		{
			input:		"	hello  world  ",
			expected:	[]string{"hello", "world"},
		},
		{
			input:		"Hello World",
			expected:	[]string{"hello", "world"},
		},
		{
			input:		"\n\thello world\t\n",
			expected:	[]string{"hello", "world"},
		},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)

		if len(actual) != len(c.expected) {
			t.Errorf("error: actual length not equal to expected length")
		}

		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]

			if word != expectedWord {
				t.Errorf("test %d: expected: %v, actual: %v", i+1, expectedWord, word)
			}
		}
	}
}
