package main

import (
	"bytes"
	"encoding/json"
	"os/exec"
)

type cmdProbeOutput struct {
	Stream []map[string]any `json:"stream"`
}

func getVideoAspectRatio(filePath string) (string, error) {
	cmd := exec.Command("ffprobe", "-v", "error", "-print_format", "json", "-show_streams", filePath)

	buffer := bytes.Buffer{}

	cmd.Stdout = &buffer
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	var data cmdProbeOutput
	var videoWidth int
	var videoHeight int

	err = json.Unmarshal(buffer.Bytes(), &data)
	if err != nil {
		return "", err
	}

	if rawVideoWidth, ok := data.Stream[1]["width"].(float64); ok {
		videoWidth = int(rawVideoWidth)
	}
	if rawVideoHeight, ok := data.Stream[1]["height"].(float64); ok {
		videoHeight = int(rawVideoHeight)
	}

	if videoWidth/videoHeight*9 == 16 {
		return "16:9", nil
	} else if videoWidth/videoHeight*16 == 9 {
		return "9:16", nil
	} else {
		return "other", nil
	}
}
