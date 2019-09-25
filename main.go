package main

import (
	"bufio"
	"log"
	"os"
	"strings"

	"github.com/urfave/cli"
)

func main() {
	var input string
	var output string
	app := cli.NewApp()

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "input, i",
			Usage:       "path to input file",
			Destination: &input,
			Required:    true,
		},
		cli.StringFlag{
			Name:        "output, o",
			Usage:       "path to output file",
			Destination: &output,
			Required:    true,
		},
	}

	app.Action = func(c *cli.Context) error {
		text := extractSMTfromFile(input)
		writeSMTtoFile(text, output)
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}

func writeSMTtoFile(text []string, output string) {
	f, err := os.Create(output)
	check(err)
	for _, s := range text {
		f.WriteString(s)
	}
	f.Close()

}

// Reading files requires checking most calls for errors.
// This helper will streamline our error checks below.
func check(e error) {
	if e != nil {
		panic(e)
	}
}

func extractFunction(text string, functionName string, delimiter string) string {
	var tmp string
	if strings.Contains(text, functionName) {
		//fmt.Println(functionName)
		tmp = text[strings.Index(text, functionName)+len(functionName):]
		end := strings.Index(tmp, delimiter)
		if end > 0 {
			tmp = tmp[:end]
		}
	}
	return tmp
}

func extractSMTfromFile(input string) []string {
	var out []string
	startEvent := "--start-event-name "
	endEvent := "--end-event-name "
	conditionFunction := "--condition-function \""
	delayFunction := "--delay-function \""
	f, err := os.Open(input)
	check(err)

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		text := scanner.Text()
		tmp := extractFunction(text, startEvent, " ")
		if len(tmp) > 0 {
			out = append(out, "Start Event: "+tmp+", ")
		}
		tmp = extractFunction(text, endEvent, " ")
		if len(tmp) > 0 {
			out = append(out, "End Event: "+tmp+", ")
		}
		tmp = extractFunction(text, conditionFunction, "\"")
		if len(tmp) > 0 {
			out = append(out, "Input: "+tmp+", ")
		}
		tmp = extractFunction(text, delayFunction, "\"")
		if len(tmp) > 0 {
			out = append(out, "Output: "+tmp+"\n")
		}

	}
	return out
}
