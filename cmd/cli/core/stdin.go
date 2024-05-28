package core

import (
	"bufio"
	"flag"
	"log"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

func ReadFromStdin() (input string, err error) {
	if flag.NArg() != 0 {
		input = flag.Arg(0)
		return input, nil
	}

	reader := bufio.NewReader(os.Stdin)

	input, err = reader.ReadString(';')
	if err != nil {
		log.Fatalln("failed to read input")
	}

	input = strings.TrimSpace(input)
	logrus.Debugf("File content:\n", input)

	return input, err
}
