package main

import (
	"bufio"
	"context"
	"fmt"
	"os"

	"github.com/bytedance/sonic"
)

func readFlags() []string {
	flags := make([]string, 0)
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			break
		}
		flags = append(flags, line)
	}
	return flags
}

func main() {
	if len(os.Args) < 3 {
		fmt.Printf("usage: %s <url> <token>\n", os.Args[0])
		os.Exit(1)
	}
	url, token := os.Args[1], os.Args[2]
	sender := NewClient(url, token)

	flagInfoWriter := bufio.NewWriter(os.Stdout)
	errWriter := bufio.NewWriter(os.Stderr)
	defer errWriter.Flush()
	defer flagInfoWriter.Flush()

	flags := readFlags()

	flagMap, err := sender.SendFlags(context.Background(), flags)
	if err != nil {
		errWriter.Write([]byte(err.Error()))
		return
	}

	data, err := sonic.Marshal(flagMap)
	if err != nil {
		errWriter.Write([]byte(err.Error()))
		return
	}
	flagInfoWriter.Write(data)
}
