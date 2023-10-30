package gosha

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/samber/lo"
	"golang.org/x/exp/slices"
	"golang.org/x/tools/go/packages"
)

func CalculatePackageSHA(
	pkgDir string,
	includeSTDLib bool,
	includeTestFiles bool,
) ([]byte, error) {
	cfg := &packages.Config{
		Mode: packages.NeedDeps |
			packages.NeedImports |
			packages.NeedName |
			packages.NeedEmbedFiles |
			packages.NeedFiles,
		Dir: pkgDir,
	}
	pkgs, err := packages.Load(cfg, ".")
	if err != nil {
		return nil, fmt.Errorf("could not open packages: %w", err)
	}

	err = nil
	packages.Visit(pkgs, nil, func(p *packages.Package) {
		for _, e := range p.Errors {
			err = errors.Join(err, e)
		}
	})

	if err != nil {
		return nil, fmt.Errorf("while loading packages:\n%w", err)
	}

	allPackages := map[string]*packages.Package{}

	for _, pkg := range pkgs {
		visitDeps(pkg, allPackages)
	}

	packageNames := lo.Keys(allPackages)
	slices.Sort(packageNames)

	sum := sha256.New()

	for _, packageName := range packageNames {

		parts := strings.Split(packageName, "/")
		if len(parts) > 0 {
			if !includeSTDLib && !strings.Contains(parts[0], ".") {
				continue
			}
		}

		// fmt.Println("package", packageName)
		sum.Write([]byte(packageName))
		pkg := allPackages[packageName]

		packageFiles := []string{}
		packageFiles = append(packageFiles, pkg.GoFiles...)
		packageFiles = append(packageFiles, pkg.EmbedFiles...)
		packageFiles = append(packageFiles, pkg.OtherFiles...)
		packageFiles = append(packageFiles, pkg.IgnoredFiles...)

		if !includeTestFiles {
			packageFiles = lo.Filter(packageFiles, func(fn string, _ int) bool {
				return !strings.HasSuffix(fn, "_test.go")
			})
		}

		err = sortAndCopyFiles(sum, packageFiles)
		if err != nil {
			return nil, fmt.Errorf("could not sha package %s: %w", pkg.PkgPath, err)
		}
	}

	return sum.Sum(nil), nil

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
