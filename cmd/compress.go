package cmd

import (
	"time"

	"immich-compress/compress"

	"github.com/spf13/cobra"
)

var flagsCompress struct {
	flagServer string
	flagAPIKey string
}

// Config holds configuration for compression command
type Config struct {
	Parallel int
	Limit    int
	Server   string
	APIKey   string
	After    time.Time
}

// compressCmd represents the compress command
var compressCmd = &cobra.Command{
	Use:   "compress",
	Short: "Compress existing fotos/videos",
	Long:  `A longer description TODO`,
	RunE: func(cmd *cobra.Command, args []string) error {
		config := compress.Config{
			Parallel: flagsRoot.flagParallel,
			Limit:    flagsRoot.flagLimit,
			Server:   flagsCompress.flagServer,
			APIKey:   flagsCompress.flagAPIKey,
			After:    flagsRoot.flagAfter,
		}
		return compress.Compressing(cmd.Context(), config)
	},
}

func init() {
	rootCmd.AddCommand(compressCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	compressCmd.PersistentFlags().StringVarP(&flagsCompress.flagServer, "server", "s", "", "The immich server address")
	if err := compressCmd.MarkPersistentFlagRequired("server"); err != nil {
		panic(err)
	}
	compressCmd.PersistentFlags().StringVarP(&flagsCompress.flagAPIKey, "api-key", "a", "", "The immich server API key")
	if err := compressCmd.MarkPersistentFlagRequired("api-key"); err != nil {
		panic(err)
	}

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// compressCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
