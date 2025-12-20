package libs

type TranscodeAudioParams struct {
	AudioFilePath  string `json:"audio_file_path"`
	OutputFilePath string `json:"output_file_path"`
}

type TranscodingClientType struct {
	TranscodeAudio func(params TranscodeAudioParams) (string, error)
}
