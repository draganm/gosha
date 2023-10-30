package main

import (
	"encoding/hex"
	"fmt"

	"github.com/draganm/gosha/gosha"
	"github.com/urfave/cli/v2"
)

func main() {

	cliFlags := struct {
		includeStdlib    bool
		includeTestFiles bool
	}{}
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:        "include-stdlib",
				EnvVars:     []string{"INCLUDE_STDLIB"},
				Destination: &cliFlags.includeStdlib,
			},
			&cli.BoolFlag{
				Name:        "include-testfiles",
				EnvVars:     []string{"INCLUDE_TESTFILES"},
				Destination: &cliFlags.includeTestFiles,
			},
		},
		Action: func(c *cli.Context) error {

			packagePath := c.Args().First()

			if packagePath == "" {
				packagePath = "."
			}

			finalSHA, err := gosha.CalculatePackageSHA(packagePath, cliFlags.includeStdlib, cliFlags.includeTestFiles)
			if err != nil {
				return fmt.Errorf("could not calculate package sha: %w", err)
			}
			fmt.Println(hex.EncodeToString(finalSHA))

			return nil
		},
	}

	app.RunAndExitOnError()

}
