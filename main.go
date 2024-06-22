package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"os/exec"
	"strings"
)

const (
	DefaultVoice = "ru_RU-ruslan-medium.onnx"
	DefaultPort  = "8080"
	ModelsDir    = "models"
	PiperPath    = "./bin/piper/piper"
)

func runExecutable(input string, voice string) (io.Reader, error) {
	voice = ModelsDir + "/" + voice
	cmd := exec.Command(PiperPath, "--model", voice, "--output_file", "-")

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	err = cmd.Start()
	if err != nil {
		return nil, err
	}

	_, err = io.WriteString(stdin, input)
	if err != nil {
		return nil, err
	}
	stdin.Close()
	return stdoutPipe, nil
}

func handleGetRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "audio/wav")
	w.WriteHeader(http.StatusOK)

	inputText := strings.TrimSpace(r.URL.Query().Get("text"))
	if inputText == "" {
		http.Error(w, "Missing Text Parameter.", http.StatusBadRequest)
	}

	voice := r.URL.Query().Get("voice")
	if voice == "" {
		voice = DefaultVoice
	}

	if voice != DefaultVoice {
		http.Error(w, "Voice not supported.", http.StatusBadRequest)
	}

	stdoutPipe, err := runExecutable(inputText, voice)
	if err != nil {
		http.Error(w, "Error running executable", http.StatusInternalServerError)
		return
	}

	_, err = io.Copy(w, stdoutPipe)
	if err != nil {
		http.Error(w, "Error streaming audio data", http.StatusInternalServerError)
		return
	}
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/api/tts", handleGetRequest).Methods("GET")

	http.Handle("/", r)

	fmt.Println("Server listening on port " + DefaultPort)
	http.ListenAndServe(":"+DefaultPort, nil)
}
