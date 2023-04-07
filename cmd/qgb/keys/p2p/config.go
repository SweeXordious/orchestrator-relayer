package p2p

import (
	"github.com/celestiaorg/orchestrator-relayer/cmd/qgb/base"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
)

func keysConfigFlags(cmd *cobra.Command) *cobra.Command {
	// TODO default value should be given
	cmd.Flags().String(base.FlagHome, "", "The qgb p2p keys home directory")
	return cmd
}

type KeysConfig struct {
	home string
}

func parseKeysConfigFlags(cmd *cobra.Command, serviceName string) (KeysConfig, error) {
	homeDir, err := cmd.Flags().GetString(flags.FlagHome)
	if err != nil {
		return KeysConfig{}, err
	}
	if homeDir == "" {
		var err error
		homeDir, err = base.DefaultServicePath(serviceName)
		if err != nil {
			return KeysConfig{}, err
		}
	}
	return KeysConfig{
		home: homeDir,
	}, nil
}
