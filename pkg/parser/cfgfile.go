package parser

import (
	"bufio"
	"bytes"
	"errors"
	"strings"
)

func GetDivisionFilename(configTSH []byte, divisionName string) (string, error) {
	reader := bytes.NewReader(configTSH)
	bufReader := bufio.NewReader(reader)

	for {
		line, err := bufReader.ReadString('\n')

		if err != nil {
			if line == "" {
				break
			}
		}
		line = strings.TrimSpace(line)

		if strings.HasPrefix(line, "division ") {
			div := strings.Fields(line)
			if div[1] == divisionName {
				return div[2], nil
			}
		}
	}

	return "", errors.New("division " + divisionName + " not found")
}
