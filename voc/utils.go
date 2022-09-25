package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/urfave/cli/v2"
)

func UserPrompt(label string) string {
	var s string
	r := bufio.NewReader(os.Stdin)
	for {
		fmt.Fprint(os.Stderr, label+" ")
		s, _ = r.ReadString('\n')
		if s != "" {
			break
		}
	}
	return strings.TrimSpace(s)
}

func selectVinOrThrowError(ctx *cli.Context) error {
	if selectedVin != "" {
		return nil
	}
	if Config.MyCarVIN != "" {
		selectedVin = Config.MyCarVIN
		return nil
	}
	return fmt.Errorf("VIN must be provided either manually or in $HOME/.voc.conf")
}

func commonFlagsVin() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "vin",
			Usage:       "Identify the car by its VIN",
			Value:       "",
			Destination: &selectedVin,
		},
	}
}
