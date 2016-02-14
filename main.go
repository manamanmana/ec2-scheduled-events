package main

import (
	"flag"
	"fmt"
	"github.com/manamanmana/ec2-scheduled-events/operations"
	"os"
)

var (
	exitcode int
	region   string
)

func init() {
	flag.StringVar(&region, "region", "", "AWS Region")
}

func main() {
	flag.Parse()
	if region == "" {
		fmt.Fprintln(os.Stderr, "Please specify region with --region option.")
		os.Exit(1)
	}

	var (
		ec2events *operations.EC2Events = operations.NewEC2Event(region)
		err       error
		outputs   []*string
		line      *string
	)

	err = ec2events.Events()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error on EC2Events.Events()")
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(2)
	}

	outputs = ec2events.Results()
	for _, line = range outputs {
		fmt.Fprintln(os.Stdout, *line)
	}

	os.Exit(0)
}
