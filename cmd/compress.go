package cmd

import (
	"fmt"
	"strings"

	"immich-compress/compress"

	"github.com/spf13/cobra"
)

// formatSlice converts ImageFormat slice to string slice
func formatSlice[T ~string](formats []T) []string {
	result := make([]string, len(formats))
	for i, format := range formats {
		result[i] = string(format)
	}
	return result
}

var flagsCompress struct {
	flagDiff           int
	flagServer         string
	flagAPIKey         string
	flagAssetType      string
	flagAssetUUIDs     []string
	flagImageQuality   int
	flagImageFormat    string
	flagVideoQuality   int
	flagVideoFormat    string
	flagVideoContainer string
}

// Config holds configuration for compression command
// type Config struct {
// 	Parallel     int
// 	Limit        int
// 	AssetType    string
// 	Server       string
// 	APIKey       string
// 	After        time.Time
// 	ImageQuality int
// 	ImageFormat  compress.ImageFormat
// }

// compressCmd represents the compress command
var compressCmd = &cobra.Command{
	Use:   "compress",
	Short: "Compress existing fotos/videos",
	Long:  `A longer description TODO`,
	RunE: func(cmd *cobra.Command, args []string) error {
		config := compress.Config{
			Parallel:       flagsRoot.flagParallel,
			Limit:          flagsRoot.flagLimit,
			AssetType:      flagsCompress.flagAssetType,
			AssetUUIDs:     flagsCompress.flagAssetUUIDs,
			Server:         flagsCompress.flagServer,
			APIKey:         flagsCompress.flagAPIKey,
			After:          flagsRoot.flagAfter,
			DiffPercent:    flagsCompress.flagDiff,
			ImageQuality:   flagsCompress.flagImageQuality,
			ImageFormat:    (compress.ImageFormat)(strings.ToLower(strings.TrimSpace(flagsCompress.flagImageFormat))),
			VideoContainer: (compress.VideoContainer)(strings.ToLower(strings.TrimSpace(flagsCompress.flagVideoContainer))),
			VideoFormat:    (compress.VideoFormat)(strings.ToLower(strings.TrimSpace(flagsCompress.flagVideoFormat))),
			VideoQuality:   flagsCompress.flagVideoQuality,
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
	compressCmd.PersistentFlags().StringVarP(&flagsCompress.flagAssetType, "type", "i", "ALL", "Asset type to compress (IMAGE, VIDEO, ALL)")
	compressCmd.PersistentFlags().StringArrayVarP(&flagsCompress.flagAssetUUIDs, "uuid", "u", []string{}, "Asset UUID")
	compressCmd.PersistentFlags().IntVarP(&flagsCompress.flagImageQuality, "image-quality", "q", 80, "Image quality for compression (1-100)")
	compressCmd.PersistentFlags().StringVarP(&flagsCompress.flagImageFormat, "image-format", "f", string(compress.JXL), fmt.Sprintf("Image format for compression (%v)", strings.Join(formatSlice(compress.ImageFormatsAvailable), ", ")))
	compressCmd.PersistentFlags().IntVarP(&flagsCompress.flagVideoQuality, "video-quality", "Q", 25, "Video quality for compression (1-100). Lower is higher quality")
	compressCmd.PersistentFlags().StringVarP(&flagsCompress.flagVideoFormat, "video-format", "F", string(compress.AV1), fmt.Sprintf("Video format for compression (%v)", strings.Join(formatSlice(compress.VideoFormatsAvailable), ", ")))
	compressCmd.PersistentFlags().StringVarP(&flagsCompress.flagVideoContainer, "video-container", "C", string(compress.MKV), fmt.Sprintf("Video container format (%v)", strings.Join(formatSlice(compress.VideoContainersAvailable), ", ")))
	compressCmd.PersistentFlags().IntVarP(&flagsCompress.flagDiff, "diff-percents", "D", 8, "If size diff is lower than this percent files will not be replaced with new.")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// compressCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
