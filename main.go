package main

import (
	"fmt"
	"github.com/lipaysamart/build-k8s-assistant-with-cobra/cmd"
	"github.com/lipaysamart/build-k8s-assistant-with-cobra/cmd/audit"
	"github.com/spf13/cobra"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"os"
	"path/filepath"
	"strings"
)

func getCommandName() string {
	if strings.HasPrefix(filepath.Base(os.Args[0]), "kubectl-") {
		return "kubectl\u2002kopilot"
	}
	return "kopilot"
}
func main() {
	opt := cmd.NewOptions()
	err := opt.Complete()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return
	}
	rootCmd := cobra.Command{
		Use: getCommandName(),
		Long: fmt.Sprintf(`You need to set TWO ENVs to run Kopilot.
Set %s to specify your token.
Set %s to specify the language. Valid options like Chinese, French, Spain, etc.`, cmd.EnvToken, cmd.EnvLang),
	}
	rootCmd.AddCommand(
		audit.New(opt),
	)
	_ = rootCmd.Execute()
}
