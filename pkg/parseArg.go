package pkg

import (
	"errors"
	"fmt"
)

// disclaimer:
// parseArgStr is modified from github user id: laurent22's project massren
// https://github.com/laurent22/massren/blob/ae4c57da1e09a95d9383f7eb645a9f69790dec6c/main.go#L172
func parseArgStr(editorCmd string) (string, []string, error) {
	const startState = 0
	const quotesState = 1
	const argState = 2

	var args []string
	state := startState
	current := ""
	quote := "\""

	for i := 0; i < len(editorCmd); i++ {
		c := editorCmd[i]

		if state == quotesState {
			if string(c) != quote {
				current += string(c)
			} else {
				args = append(args, current)
				current = ""
				state = startState
			}
			continue
		}

		if c == '"' || c == '\'' {
			state = quotesState
			quote = string(c)
			continue
		}

		if state == argState {
			if c == ' ' || c == '\t' {
				args = append(args, current)
				current = ""
				state = startState
			} else {
				current += string(c)
			}
			continue
		}

		if c != ' ' && c != '\t' {
			state = argState
			current += string(c)
		}
	}

	if state == quotesState {
		return "", []string{}, fmt.Errorf(
			"unclosed quote in command line: %s", editorCmd)
	}

	if current != "" {
		args = append(args, current)
	}

	if len(args) == 0 {
		return "", []string{}, errors.New("empty command line")
	}

	if len(args) == 1 {
		return args[0], []string{}, nil
	}

	return args[0], args[1:], nil
}
