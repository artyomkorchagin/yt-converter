package router

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/artyomkorchagin/yt-converter/internal/types"
	"github.com/gin-gonic/gin"
	"github.com/kkdai/youtube/v2"
	ffmpeg_go "github.com/u2takey/ffmpeg-go"
	"go.uber.org/zap"
)

func (h *Handler) getVideo(c *gin.Context) error {
	videoID := c.Param("id")
	client := youtube.Client{}

	video, err := client.GetVideo(videoID)
	if err != nil {
		return types.ErrBadRequest(fmt.Errorf("invalid video ID: %w", err))
	}

	// === 1. Get best VIDEO-ONLY stream (audio channels = 0, has width) ===
	videoFormats := video.Formats.Select(func(f youtube.Format) bool {
		return f.AudioChannels == 0 && f.Width > 0
	})
	if len(videoFormats) == 0 {
		return types.ErrNotFound(fmt.Errorf("no video-only format"))
	}
	videoFormats.Sort() // sorts by resolution, FPS, codec — highest quality first
	bestVideo := &videoFormats[0]

	// === 2. Get best AUDIO-ONLY stream (width = 0, has audio channels) ===
	audioFormats := video.Formats.WithAudioChannels()
	if len(audioFormats) == 0 {
		return types.ErrNotFound(fmt.Errorf("no audio-only format"))
	}
	audioFormats.Sort() // sorts by codec, channels, bitrate — best audio first
	bestAudio := &audioFormats[0]

	// === 3. Create temp dir ===
	tempDir := fmt.Sprintf("./temp_%d", time.Now().UnixNano())
	err = os.MkdirAll(tempDir, 0755)
	if err != nil {
		return types.ErrInternalServerError(fmt.Errorf("failed to create temp dir: %w", err))
	}

	// Clean up temp directory when function exits
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			h.logger.Error("failed to remove temp dir", zap.Error(err))
		}
	}()

	videoFile := filepath.Join(tempDir, "video.mp4")
	audioFile := filepath.Join(tempDir, "audio.m4a")
	outputFile := filepath.Join(tempDir, "final.mp4")

	// === 4. Download video ===
	if err := downloadStream(client, video, bestVideo, videoFile); err != nil {
		return types.ErrInternalServerError(fmt.Errorf("video download failed: %w", err))
	}

	// Check if video file is valid
	if info, err := os.Stat(videoFile); err != nil || info.Size() == 0 {
		return types.ErrInternalServerError(fmt.Errorf("video file invalid"))
	}

	// === 5. Download audio ===
	if err := downloadStream(client, video, bestAudio, audioFile); err != nil {
		return types.ErrInternalServerError(fmt.Errorf("audio download failed: %w", err))
	}

	// Check if audio file is valid
	if info, err := os.Stat(audioFile); err != nil || info.Size() == 0 {
		return types.ErrInternalServerError(fmt.Errorf("audio file invalid"))
	}

	// === 6. Merge with ffmpeg-go
	videoStream := ffmpeg_go.Input(videoFile)
	audioStream := ffmpeg_go.Input(audioFile)

	err = ffmpeg_go.
		Output([]*ffmpeg_go.Stream{videoStream, audioStream}, outputFile,
			ffmpeg_go.KwArgs{
				"c:v":      "copy",
				"c:a":      "aac",
				"strict":   "experimental",
				"movflags": "+faststart",
			}).
		OverWriteOutput().
		Run()

	if err != nil {
		return types.ErrInternalServerError(fmt.Errorf("ffmpeg merge failed: %w", err))
	}

	// === 7. Stream to user ===
	mergedFile, err := os.Open(outputFile)
	if err != nil {
		return types.ErrInternalServerError(fmt.Errorf("failed to open merged file: %w", err))
	}

	// Ensure file is closed before temp dir cleanup
	defer func() {
		if closeErr := mergedFile.Close(); closeErr != nil {
			h.logger.Error("failed to close merged file", zap.Error(closeErr))
		}
	}()

	safeTitle := sanitizeFilename(video.Title)
	fileName := fmt.Sprintf("%s.mp4", safeTitle)

	c.Header("Content-Type", "video/mp4")
	c.Header("Content-Disposition", "attachment; filename=\""+fileName+"\"")

	_, err = io.Copy(c.Writer, mergedFile)
	if err != nil {
		return types.ErrInternalServerError(fmt.Errorf("failed to stream merged file: %w", err))
	}

	return nil
}

// Helper to download a stream to a file
func downloadStream(client youtube.Client, video *youtube.Video, format *youtube.Format, filePath string) error {
	stream, _, err := client.GetStream(video, format)
	if err != nil {
		return types.ErrInternalServerError(fmt.Errorf("error getting stream: %w", err))
	}
	defer stream.Close()

	file, err := os.Create(filePath)
	if err != nil {
		return types.ErrInternalServerError(fmt.Errorf("error creating file: %w", err))
	}
	defer file.Close()

	_, err = io.Copy(file, stream)
	return err
}

func sanitizeFilename(name string) string {
	invalidChars := []string{"\\", "/", ":", "*", "?", "\"", "<", ">", "|", "\n", "\r"}
	result := name
	for _, char := range invalidChars {
		result = strings.ReplaceAll(result, char, "_")
	}
	result = strings.TrimSpace(result)
	if len(result) > 100 {
		result = result[:100]
	}
	if result == "" {
		result = "video"
	}
	return result
}
