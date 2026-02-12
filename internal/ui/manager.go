package ui

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type ModelType string

const (
	ModelTypeASR ModelType = "ASR"
	ModelTypeLLM ModelType = "LLM"
)

type ModelInfo struct {
	Name        string
	Description string
	Type        ModelType
	URL         string
	Filename    string
	Size        string
}

var AvailableModels = []ModelInfo{
	{
		Name:        "Whisper Small (Recommended)",
		Description: "Balanced speed and accuracy (~460MB)",
		Type:        ModelTypeASR,
		URL:         "https://huggingface.co/ggerganov/whisper.cpp/resolve/main/ggml-small.bin",
		Filename:    "ggml-small.bin",
		Size:        "460MB",
	},
	{
		Name:        "Whisper Base",
		Description: "Faster, less accurate (~140MB)",
		Type:        ModelTypeASR,
		URL:         "https://huggingface.co/ggerganov/whisper.cpp/resolve/main/ggml-base.bin",
		Filename:    "ggml-base.bin",
		Size:        "140MB",
	},
	{
		Name:        "TinyLlama 1.1B Chat",
		Description: "Fast, good for basic cleanup (~637MB)",
		Type:        ModelTypeLLM,
		URL:         "https://huggingface.co/TheBloke/TinyLlama-1.1B-Chat-v1.0-GGUF/resolve/main/tinyllama-1.1b-chat-v1.0.Q4_K_M.gguf",
		Filename:    "tinyllama-1.1b-chat-v1.0.Q4_K_M.gguf",
		Size:        "637MB",
	},
	{
		Name:        "Qwen 1.5 1.8B",
		Description: "Higher quality (Requires newer backend)",
		Type:        ModelTypeLLM,
		URL:         "https://huggingface.co/Qwen/Qwen1.5-1.8B-Chat-GGUF/resolve/main/qwen1_5-1_8b-chat-q4_k_m.gguf",
		Filename:    "qwen1_5-1_8b-chat-q4_k_m.gguf",
		Size:        "1.2GB",
	},
}

type ModelManager struct {
	window    fyne.Window
	statusLabel *widget.Label
}

func NewModelManager(w fyne.Window) *ModelManager {
	return &ModelManager{
		window: w,
	}
}

func (m *ModelManager) GetContent() fyne.CanvasObject {
	title := widget.NewLabelWithStyle("Sussurro Models", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	
	// Separate models by type
	var asrModels, llmModels []ModelInfo
	for _, model := range AvailableModels {
		if model.Type == ModelTypeASR {
			asrModels = append(asrModels, model)
		} else {
			llmModels = append(llmModels, model)
		}
	}

	asrContainer := m.createModelSection("Speech Recognition (Whisper)", asrModels)
	llmContainer := m.createModelSection("Language Models (LLM)", llmModels)

	m.statusLabel = widget.NewLabel("Ready")
	m.statusLabel.Alignment = fyne.TextAlignCenter

	content := container.NewVBox(
		layout.NewSpacer(),
		title,
		layout.NewSpacer(),
		asrContainer,
		layout.NewSpacer(),
		llmContainer,
		layout.NewSpacer(),
		m.statusLabel,
		layout.NewSpacer(),
	)

	return container.NewPadded(content)
}

func (m *ModelManager) createModelSection(title string, models []ModelInfo) fyne.CanvasObject {
	header := widget.NewLabelWithStyle(title, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	
	list := container.NewVBox()
	for _, model := range models {
		list.Add(m.createModelRow(model))
		list.Add(layout.NewSpacer())
	}

	return container.NewVBox(header, container.NewPadded(list))
}

func (m *ModelManager) createModelRow(model ModelInfo) fyne.CanvasObject {
	name := widget.NewLabel(model.Name)
	name.TextStyle = fyne.TextStyle{Bold: true}
	
	desc := widget.NewLabel(model.Description)
	desc.TextStyle = fyne.TextStyle{Italic: true}
	
	info := container.NewVBox(name, desc)
	
	// Check if installed
	installed := false
	if _, err := os.Stat("models/" + model.Filename); err == nil {
		installed = true
	}

	var btn *CapsuleButton
	if installed {
		btn = NewCapsuleButton("Installed", func() {}, false)
		// Disable button visually or logic? 
		// For now just secondary style
	} else {
		btn = NewCapsuleButton("Download", func() {
			go m.downloadModel(model)
		}, true) // Primary
	}

	return container.NewBorder(nil, nil, info, btn)
}

func (m *ModelManager) downloadModel(model ModelInfo) {
	m.statusLabel.SetText(fmt.Sprintf("Downloading %s...", model.Name))
	
	err := os.MkdirAll("models", 0755)
	if err != nil {
		m.statusLabel.SetText("Error creating models directory")
		return
	}

	out, err := os.Create("models/" + model.Filename)
	if err != nil {
		m.statusLabel.SetText("Error creating file")
		return
	}
	defer out.Close()

	resp, err := http.Get(model.URL)
	if err != nil {
		m.statusLabel.SetText("Error downloading")
		return
	}
	defer resp.Body.Close()

	// Simple download without progress bar update for now (limitations of simple UI)
	// TODO: Add progress bar
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		m.statusLabel.SetText("Error saving file")
		return
	}

	m.statusLabel.SetText(fmt.Sprintf("Downloaded %s", model.Name))
	
	// Refresh UI to show "Installed"
	// Ideally we trigger a refresh of the list.
	// For this simple version, we might just update status.
	m.window.Content().Refresh()
}
