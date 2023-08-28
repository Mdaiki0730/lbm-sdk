package main

import (
	"os"

	"github.com/Finschia/finschia-rdk/l2app"
	"github.com/Finschia/finschia-rdk/l2app/rollupd/cmd"
	"github.com/Finschia/finschia-rdk/server"
	svrcmd "github.com/Finschia/finschia-rdk/server/cmd"
)

func main() {
	var _ = os.Args // To avoid linter "imported but not used" false positive
	rootCmd, _ := cmd.NewRootCmd()

	if err := svrcmd.Execute(rootCmd, l2app.DefaultNodeHome); err != nil {
		switch e := err.(type) {
		case server.ErrorCode:
			os.Exit(e.Code)

		default:
			os.Exit(1)
		}
	}
}
