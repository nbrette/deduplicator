package models

type Request struct {
	Message struct {
		Attributes map[string]string
		Data       []byte
	}
}
