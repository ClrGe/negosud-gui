package data

import (
	"log"
	"os"
	"time"
)

func ErrorLogger(source string, message string) {
	// create folder if it doesn't exist
	if _, err := os.Stat("data/logs"); os.IsNotExist(err) {
		os.Mkdir("data/logs", 0755)
	}

	file, err := os.OpenFile(
		"data/logs/errors-"+time.Now().Format("2006-01-02")+".log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0664,
	)
	if err != nil {
		ErrorLogger("LOGGING", err.Error())
	}

	defer file.Close()

	logger := log.New(
		file,
		"SOURCE : "+source,
		log.Ldate|log.Ltime|log.Lmicroseconds|log.LUTC,
	)
	logger.Println(message)
}
