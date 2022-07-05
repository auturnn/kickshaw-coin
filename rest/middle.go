package rest

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/auturnn/kickshaw-coin/utils"
	"github.com/gorilla/handlers"
)

func jsonContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Add("Content-type", "application/json")
		next.ServeHTTP(rw, r)
	})
}

const windowLogName = "2006_01_02"
const defaultLogName = "2006-01-02"

func loggerMiddleware(next http.Handler) http.Handler {
	var f *os.File
	t := time.Now().Local()

	switch runtime.GOOS {
	case "windows":
		f = loggingFileOpen(t.Format(windowLogName))
	default:
		f = loggingFileOpen(t.Format(defaultLogName))
	}

	h := handlers.LoggingHandler(os.Stdout, next)
	return handlers.LoggingHandler(f, h)
}

func loggingFileOpen(fileName string) *os.File {
	logPath := "./log"
	if _, err := os.Stat(logPath); err != nil {
		if err := os.Mkdir(logPath, 0755); err != nil {
			utils.HandleError(errors.New("failed logging path create"))
		}
	}

	f, err := os.OpenFile(fmt.Sprintf("%s/%s.log", logPath, fileName), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0655)
	utils.HandleError(err)

	return f
}
