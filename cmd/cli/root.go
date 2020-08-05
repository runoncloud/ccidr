package cli

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/runoncloud/ccidr/pkg/ccidr"
	"github.com/runoncloud/ccidr/pkg/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"strings"
)

func RootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "ccidr",
		Short:         "",
		Long:          `.`,
		SilenceErrors: true,
		SilenceUsage:  true,
		PreRun: func(cmd *cobra.Command, args []string) {
			viper.BindPFlags(cmd.Flags())
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			log := logger.NewLogger()
			log.Info("")

			finishedCh := make(chan bool, 1)
			go func() {
				for {
					select {
					case <-finishedCh:
						fmt.Printf("\r")
						return
					}
				}
			}()

			defer func() {
				finishedCh <- true
			}()

			if err := ccidr.RunCommand(cmd, args); err != nil {
				return errors.Cause(err)
			}
			return nil
		},
	}
	cobra.OnInitialize(initConfig)

	var region, service string
	var remote bool

	cmd.Flags().StringVarP(&region, "region", "r", "",
		"Only selects services or ips for a specific region")

	cmd.Flags().StringVarP(&service, "service", "s", "",
		"Only selects ips for a specific service")

	cmd.Flags().BoolVarP(&remote, "remote", "R", false,
		"Fetch the values directly from the source over HTTP")

	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	return cmd
}

func InitAndExecute() {
	if err := RootCmd().Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func initConfig() {
	viper.AutomaticEnv()
}
