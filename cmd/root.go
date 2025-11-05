// Package cmd for all available commands
package cmd

import (
	"os"
	"runtime"
	"time"

	"immich-compress/immich"

	"github.com/spf13/cobra"
)

var flagsRoot struct {
	flagParallel int
	flagAfter    time.Time
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "immich-compress",
	Short: "Compress existing fotos/videos",
	Long:  `A longer description TODO`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) {
	// 	fmt.Println("Use --help To show commands.")
	// },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.immich-compress.yaml)")
	rootCmd.PersistentFlags().IntVarP(&flagsRoot.flagParallel, "parallel", "p", runtime.NumCPU(), "parallel")
	rootCmd.PersistentFlags().TimeVarP(&flagsRoot.flagAfter, "after", "t", time.Now(), []string{immich.TAG_COMPRESSED_AT_FORMAT}, "after what time we want to recompress")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
