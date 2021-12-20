package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type saver struct {
	file    *os.File
	encoder *json.Encoder
}

func NewSaver(fileName string) (*saver, error) {
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		return nil, err
	}
	return &saver{
		file:    file,
		encoder: json.NewEncoder(file),
	}, nil
}
func (p *saver) WriteEvent(event *[]MyURL) error {
	return p.encoder.Encode(&event)
}
func (p *saver) Close() error {
	return p.file.Close()
}

type loader struct {
	file    *os.File
	decoder *json.Decoder
}

func NewLoader(fileName string) (*loader, error) {
	file, err := os.OpenFile(fileName, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		return nil, err
	}
	return &loader{
		file:    file,
		decoder: json.NewDecoder(file),
	}, nil
}
func (c *loader) ReadEvent() (*[]MyURL, error) {
	event := &[]MyURL{}
	if err := c.decoder.Decode(&event); err != nil {
		return nil, err
	}
	return event, nil
}
func (c *loader) Close() error {
	return c.file.Close()
}

func LoadDate(fileName string) {
	p, err := NewLoader(fileName)
	if err != nil {
		panic(err)
	}

	defer p.Close()
	reder, err := p.ReadEvent()
	if err != nil {
		return
	}
	myurl = *reder
	fmt.Println(reder)
}

func SaveDate(fileName string) {
	p, err := NewSaver(fileName)
	if err != nil {
		panic(err)
	}
	defer p.Close()
	p.WriteEvent(&myurl)
}
