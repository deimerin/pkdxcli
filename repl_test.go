package main

import (
	"testing"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "  hello  world  ",
			expected: []string{"hello", "world"},
		},
		{
			input:    " Go           routines    are   lightweight.   ",
			expected: []string{"go", "routines", "are", "lightweight."},
		},
		{
			input:    " Channel        synchronization     in Go      ",
			expected: []string{"channel", "synchronization", "in", "go"},
		},
		{
			input:    " The   standard       library is    vast ",
			expected: []string{"the", "standard", "library", "is", "vast"},
		},
		{
			input:    "Error        handling     is           important.",
			expected: []string{"error", "handling", "is", "important."},
		},
		{
			input:    " Interfaces      provide     polymorphism   ",
			expected: []string{"interfaces", "provide", "polymorphism"},
		},
		// add more cases here
	}

	for _, c := range cases {
		actual := cleanInput(c.input)

		if len(actual) != len(c.expected) {
			t.Errorf("Actual length and expected length are different. %d != %d", len(actual), len(c.expected))
		}

		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]

			if word != expectedWord {
				t.Errorf("Words are different. %v != %v", word, expectedWord)
			}

		}
	}
}
