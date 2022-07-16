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
