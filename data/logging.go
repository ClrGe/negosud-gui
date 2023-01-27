package data

import (
	"log"
	"os"
	"time"
)

// Logger logs error message to a file under data/logs (created if folder doesn't exist)
// "source" defines the origin of the error
func Logger(isError bool, source string, message string) {
	var filename string
	// create folder if it doesn't exist
	if _, err := os.Stat("data/logs"); os.IsNotExist(err) {
		os.Mkdir("data/logs", 0755)
	}

	if isError {
		filename = "data/logs/errors-" + time.Now().Format("2006-01-02") + ".log"
	} else {
		filename = "data/logs/info-" + time.Now().Format("2006-01-02") + ".log"
	}

	file, err := os.OpenFile(
		filename,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0664,
	)
	if err != nil {
		Logger(true, "LOGGING", err.Error())
	}
	defer file.Close()

	logger := log.New(
		file,
		"SOURCE : "+source,
		log.Ldate|log.Ltime|log.Lmicroseconds|log.LUTC,
	)

	logger.Println(message)
}
