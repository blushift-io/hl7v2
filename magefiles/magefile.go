package main

import (
	"github.com/magefile/mage/sh"
)

var goDeps = map[string]string{
	"github.com/alvaroloes/enumer": "latest",
}

func Build() {}

func InstallDeps() error {
	for pkg, ver := range goDeps {
		if err := sh.Run("go", "install", pkg+"@"+ver); err != nil {
			return err
		}
	}

	return nil
}

// DocsDev starts the local VitePress documentation development server.
func DocsDev() error {
	return sh.RunV("pnpm", "--prefix", "docs", "run", "dev")
}

// DocsBuild builds the static documentation site.
func DocsBuild() error {
	return sh.RunV("pnpm", "--prefix", "docs", "run", "build")
}
