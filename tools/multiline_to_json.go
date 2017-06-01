package tools

import (
	"bufio"
	"bytes"
	"strings"
)

/* Parses a multi-line input to a JSON string

input:

key1: value1 \n
key2: value2 \n

output:

{ "key" : "value",
  "key" : "value"
}
*/

func MultiLineToJson(multiLineInput string) (string, error) {

	var result bytes.Buffer

	// strip whitespace first
	multiLineInput = strings.TrimSpace(multiLineInput)

	// get the total amount of lines
	totalLines := strings.Count(multiLineInput, "\n")

	// assign a reader
	reader := bufio.NewReader(strings.NewReader(multiLineInput))

	// start constructing the JSON string
	result.WriteString("{")

	for lineCount := 0; lineCount <= totalLines; lineCount++ {

		line, err := (reader.ReadString('\n'))
		if err != nil {
			break
		} else {
			kvPair := strings.Split(line, ": ")

			result.WriteString("\"" + kvPair[0] + "\" : \"" + strings.Trim(kvPair[1], "\n") + "\",")
		}
	}

	// loose the last "," and close with a "}"
	result.Truncate(int(len(result.String()) - 1))
	result.WriteString("}")
	return result.String(), nil
}
