package gelf_models

type ShortMessage struct {
	Level string `json:"level"`
	Msg   string `json:"msg"`
}
