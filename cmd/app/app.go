package app

import (
	"fmt"
	"token-price/internal/memory"

	"token-price/internal/config"
	"token-price/internal/router"
	"token-price/pkg/log"

	"github.com/spf13/cobra"
)

var (
	cfgFile string
)

func NewScanCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "token-price",
		Short:   "get token usd price",
		Version: "1.0.0",
		RunE:    RunCommand,
	}

	cmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "etc/dev.yaml", "config file path")

	return cmd
}

func RunCommand(cmd *cobra.Command, args []string) error {
	c, err := config.Load(cfgFile)
	if err != nil {
		return fmt.Errorf("load config error: %s", err)
	}

	log.InitLogger(&c.Log)
	log.Infof("config: %+v", c)

	memory.InitCache()

	e := router.Router(c)
	if err = e.Run(c.Http.Addr); err != nil {
		return err
	}

	return e.Run(c.Http.Addr)
}
