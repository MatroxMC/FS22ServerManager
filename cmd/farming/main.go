package farming

import (
	"bufio"
	"log"
	"os"
	"time"
)

const (
	GameDir = "D:\\Jeux\\Farming Simulator 22"
)

func Logi(message string) {
	h := time.Now()
	file, err := os.OpenFile("logs.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}

	datawriter := bufio.NewWriter(file)

	_, _ = datawriter.WriteString(h.Format("3:4:5") + ": " + message + "\n")

	datawriter.Flush()
}
