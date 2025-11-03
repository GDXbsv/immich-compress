package cmd

import (
	"immich-compress/compress"

	"github.com/spf13/cobra"
)

var flagsCompress struct {
	flagServer string
	flagApiKey string
}

// compressCmd represents the compress command
var compressCmd = &cobra.Command{
	Use:   "compress",
	Short: "Compress existing fotos/videos",
	Long:  `A longer description TODO`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return compress.Compress(cmd.Context(), flagsRoot.flagParallel, flagsCompress.flagServer, flagsCompress.flagApiKey)
	},
}

func init() {
	rootCmd.AddCommand(compressCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	compressCmd.PersistentFlags().StringVarP(&flagsCompress.flagServer, "server", "s", "", "The immich server address")
	compressCmd.MarkFlagRequired("server")
	compressCmd.PersistentFlags().StringVarP(&flagsCompress.flagApiKey, "api-key", "a", "", "The immich server API key")
	compressCmd.MarkFlagRequired("api-key")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// compressCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
