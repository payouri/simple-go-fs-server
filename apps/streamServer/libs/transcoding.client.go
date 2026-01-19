package libs

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
)

type BuildTranscodingClientParams struct {
}

type TranscodeAudioParams struct {
	AudioFilePath  string `json:"audio_file_path"`
	OutputFilePath string `json:"output_file_path"`
}

type AudioFormat int

const (
	MP3 AudioFormat = iota
	AAC
	FLAC
)

var AudioOutputFormat = map[AudioFormat]string{
	MP3:  "mp3",
	AAC:  "aac",
	FLAC: "flac",
}

type InMemoryTranscodeAudioParams struct {
	OutputFormat  AudioFormat
	InputStream   io.ReadCloser
	ContentLength int64
}

type TranscodingClientType struct {
	TranscodeAudioSync     func(params TranscodeAudioParams) error
	InMemoryTranscodeAudio func(params InMemoryTranscodeAudioParams) (*io.PipeReader, error)
}

func buildTranscodeAudioSync(dependencies BuildTranscodingClientParams) func(params TranscodeAudioParams) error {
	return func(params TranscodeAudioParams) error {
		return nil
	}
}

func buildInMemoryTranscodeAudio(dependencies BuildTranscodingClientParams) func(params InMemoryTranscodeAudioParams) (*io.PipeReader, error) {
	return func(params InMemoryTranscodeAudioParams) (*io.PipeReader, error) {
		reader, writer := io.Pipe()

		// Build the FFmpeg command
		cmd := exec.Command(
			"ffmpeg",
			"-loglevel", "error", // Log level
			"-vn", // Disable video
			// "-ar", "44100", // Set audio sample rate
			"-i", "pipe:0", // Read from stdin
			"-codec:a", "libmp3lame",
			"-b:a", "32k",
			"-f", AudioOutputFormat[params.OutputFormat], "pipe:1", // Output format (e.g., "flac", "mp3", "aac")
			// "-ac", "2",
			// "-c:a", "mp3", // Audio codec (e.g., "aac", "libmp3lame", "flac")
			// "-b:a", "192k", // Bitrate
			// "-", // Write to stdout
		)
		cmd.Stdin = params.InputStream // FFmpeg reads from the pipe
		cmd.Stdout = writer
		cmd.Stderr = os.Stderr

		// Start FFmpeg
		err := cmd.Start()
		if err != nil {
			return nil, fmt.Errorf("failed to start FFmpeg: %v", err)
		}

		go func() {
			defer writer.Close() // Important to close the writer when done
			defer params.InputStream.Close()

			if err := cmd.Wait(); err != nil {
				log.Printf("ffmpeg command finished with error: %v\n", err)
			} else {
				log.Println("ffmpeg command finished successfully")
			}
		}()
		// Return the stdout reader (transcoded audio stream)
		return reader, nil
	}
}

func BuildTranscodingClient(dependencies BuildTranscodingClientParams) TranscodingClientType {
	transcodeAudioSync := buildTranscodeAudioSync(dependencies)
	transcodeAudio := buildInMemoryTranscodeAudio(dependencies)

	return TranscodingClientType{
		TranscodeAudioSync:     transcodeAudioSync,
		InMemoryTranscodeAudio: transcodeAudio,
	}
}

var TranscodingClient = BuildTranscodingClient(BuildTranscodingClientParams{})
