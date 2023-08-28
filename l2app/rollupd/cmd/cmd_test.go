package cmd_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Finschia/finschia-rdk/l2app"
	"github.com/Finschia/finschia-rdk/l2app/rollupd/cmd"
	svrcmd "github.com/Finschia/finschia-rdk/server/cmd"
	"github.com/Finschia/finschia-rdk/x/genutil/client/cli"
)

func TestInitCmd(t *testing.T) {
	t.Skipf("🔬 The rollkit/cosmos-sdk also remains faulty.")
	rootCmd, _ := cmd.NewRootCmd()
	rootCmd.SetArgs([]string{
		"init",       // Test the init cmd
		"l2app-test", // Moniker
		fmt.Sprintf("--%s=%s", cli.FlagOverwrite, "true"), // Overwrite genesis.json, in case it already exists
	})

	require.NoError(t, svrcmd.Execute(rootCmd, l2app.DefaultNodeHome))
}
