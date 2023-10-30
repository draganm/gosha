package main

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/samber/lo"
	"github.com/urfave/cli/v2"
	"golang.org/x/exp/slices"
	"golang.org/x/tools/go/packages"
)

func main() {

	app := &cli.App{
		Flags: []cli.Flag{},
		Action: func(c *cli.Context) error {

			packagePath := c.Args().First()

			if packagePath == "" {
				packagePath = "."
			}

			cfg := &packages.Config{
				Mode: packages.NeedDeps |
					packages.NeedImports |
					packages.NeedName |
					packages.NeedEmbedFiles |
					packages.NeedFiles,
				Dir: packagePath,
			}
			pkgs, err := packages.Load(cfg, ".")
			if err != nil {
				return fmt.Errorf("could not open packages: %w", err)
			}

			if packages.PrintErrors(pkgs) > 0 {
				return errors.New("packages.Load returned errors")
			}

			allPackages := map[string]*packages.Package{}

			for _, pkg := range pkgs {
				visitDeps(pkg, allPackages)
			}

			packageNames := lo.Keys(allPackages)
			slices.Sort(packageNames)

			sum := sha256.New()

			for _, packageName := range packageNames {
				sum.Write([]byte(packageName))
				pkg := allPackages[packageName]

				packageFiles := []string{}
				packageFiles = append(packageFiles, pkg.GoFiles...)
				packageFiles = append(packageFiles, pkg.EmbedFiles...)
				packageFiles = append(packageFiles, pkg.OtherFiles...)
				packageFiles = append(packageFiles, pkg.IgnoredFiles...)

				err = sortAndCopyFiles(sum, packageFiles)
				if err != nil {
					return fmt.Errorf("could not sha package %s: %w", pkg.PkgPath, err)
				}
			}

			finalSha := sum.Sum(nil)

			fmt.Println(hex.EncodeToString(finalSha))

			return nil
		},
	}

	app.RunAndExitOnError()

}

func sortAndCopyFiles(destination io.Writer, files []string) error {
	sort.Strings(files)
	for _, f := range files {
		err := copyFile(destination, f)
		if err != nil {
			return err
		}
	}

	return nil
}

func copyFile(destination io.Writer, fileName string) error {
	f, err := os.Open(fileName)
	if err != nil {
		return fmt.Errorf("could not open %s :%w", fileName, err)
	}

	defer f.Close()
	_, err = io.Copy(destination, f)
	if err != nil {
		return fmt.Errorf("failed to copy %s: %w", fileName, err)
	}

	fmt.Println("copied", fileName)

	return nil
}

func visitDeps(pkg *packages.Package, visit map[string]*packages.Package) {
	_, alreadyVisited := visit[pkg.PkgPath]
	if alreadyVisited {
		return
	}
	visit[pkg.PkgPath] = pkg

	for _, dep := range pkg.Imports {
		visitDeps(dep, visit)
	}

}
