package main

import (
	"io"
	"log"
	"os"
)

func SetLogFile(f string) error {
	//link the log file to the console
	file, err := os.OpenFile(f, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}

	m := io.MultiWriter(os.Stdout, file) //Create a multi-writer to write to both the console and the file thx to copilot!!
	log.SetOutput(m)

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	return nil
}
