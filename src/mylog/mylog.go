package mylog

import (
	"log"
	"os"
)

var (
	Debug *log.Logger
	Info  *log.Logger
	Error *log.Logger
	Trace *log.Logger
	Warn  *log.Logger
)

func init() {
	log.SetFlags(2)
	log.Println("init ...")
	Debug = log.New(os.Stdout, "[DEBUG] ", log.Ldate|log.Ltime|log.Lshortfile)
	Info = log.New(os.Stdout, "[INFO] ", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(os.Stderr, "[ERROR] ", log.Ldate|log.Ltime|log.Lshortfile)
	Trace = log.New(os.Stderr, "[Trace] ", log.Ldate|log.Ltime|log.Lshortfile)
	Warn = log.New(os.Stderr, "[Warn] ", log.Ldate|log.Ltime|log.Lshortfile)
}
