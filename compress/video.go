package compress

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"

	"immich-compress/immich"

	"github.com/google/uuid"
)

type VideoConfig struct {
	Container VideoContainer
	Format    VideoFormat
	Quality   int
}

type VideoContainer string

const (
	MKV VideoContainer = "mkv"
	MP4 VideoContainer = "mp4"
)

var VideoContainersAvailable = []VideoContainer{MKV, MP4}

type VideoFormat string

const (
	AV1  VideoFormat = "av1"
	HEVC VideoFormat = "hevc"
	H264 VideoFormat = "h264"
)

var VideoFormatsAvailable = []VideoFormat{AV1, HEVC, H264}

func (c *VideoConfig) compress(client *immich.ClientSimple, asset immich.AssetResponseDto) (*os.File, error) {
	uuid, err := uuid.Parse(asset.Id)
	if err != nil {
		return nil, fmt.Errorf("failed to parse uuid '%s': %w", asset.Id, err)
	}

	// Create temporary input file
	fileIn, err := os.Create(filepath.Join(os.TempDir(), fmt.Sprintf("%s%s", uuid.String(), filepath.Ext(asset.OriginalPath))))
	if err != nil {
		return nil, fmt.Errorf("failed to create temp input file: %w", err)
	}
	defer fileIn.Close()
	defer os.Remove(fileIn.Name())

	// Create temporary output file
	fileOut, err := os.Create(filepath.Join(os.TempDir(), fmt.Sprintf("%s.%s", filepath.Base(asset.OriginalPath), c.Container)))
	if err != nil {
		return nil, fmt.Errorf("failed to create temp output file: %w", err)
	}
	defer fileOut.Close()

	// Download video to temporary input file
	resp, err := client.AssetDownload(uuid)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch video: %w", err)
	}
	defer resp.Body.Close()

	// Copy downloaded video to input file
	_, err = io.Copy(fileIn, resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to save video to temp file: %w", err)
	}
	fileIn.Close() // Close to ensure all data is written

	args := make([]string, 20)
	args = append(args,
		"-i", fileIn.Name(),
	)

	switch c.Format {
	case AV1:
		args = append(args,
			"-c:v", "libsvtav1",
			"-pix_fmt", "yuv420p10le", // This flag enables 10-bit encoding
			"-preset", "5", // Was 8. Lower is slower but better compression.
			"-c:a", "libopus",
			"-b:a", "128k",
		)
	case HEVC:
		args = append(args,
			"-c:v", "libx265",
			// Use the "main10" profile for 10-bit encoding
			"-profile:v", "main10",
			"-pix_fmt", "yuv420p10le", // Explicitly set 10-bit pixel format
			// "VideoPreset": libx265 uses names, not numbers.
			// 'slow' is a good equivalent to svt-av1's preset '5'.
			// Other options: medium (default), fast, faster, etc.
			"-preset", "slow",
			"-c:a", "libopus",
			"-b:a", "128k",
		)
	case H264:
		args = append(args,
			"-c:v", "libx264",
			// "Profile": "high" & 8-bit color
			// We use 8-bit for maximum compatibility.
			"-profile:v", "high",
			"-pix_fmt", "yuv420p", // 8-bit pixel format
			// "VideoPreset": libx264 uses names.
			// 'slow' is a great balance of quality and encoding time.
			"-preset", "slow",
			// "AudioEncoder": Use "aac" for MP4 containers.
			// "libfdk_aac" is the highest-quality AAC encoder.
			// If you don't have it, you can use the built-in "aac".
			"-c:a", "aac",
			"-b:a", "128k", // Set audio bitrate
		)

	default:
		return nil, fmt.Errorf("unsupported output format: %s", c.Format)
	}

	// -i: input file
	// -c:v libsvtav1: Use the SVT-AV1 video codec
	// -crf 30: Constant Rate Factor (quality). Lower is better quality,
	//          higher is smaller file. 25-35 is a good range.
	// -preset 8: SVT-AV1 speed preset. 0 (slowest, best quality)
	//            to 12 (fastest, lowest quality). 5-8 is a good balance.
	// -c:a libopus: Use the Opus audio codec, a great companion for AV1.
	// -b:a 128k: Set audio bitrate to 128kbps.
	args = append(args, []string{
		"-crf", strconv.Itoa(c.Quality), // Was 30. Lower is higher quality.
		fileOut.Name(),
	}...)

	cmd := exec.Command("ffmpeg", args...)

	// Run the command and capture its output
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("ffmpeg failed with output '' %w", string(output), err)
	}

	return fileOut, nil
}
