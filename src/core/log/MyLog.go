package log

import (
	"io"
	"log"
	"os"
)

var (
	InfoLogger *log.Logger
	WarningLogger *log.Logger
	ErrorLogger * log.Logger
)

func init(){
	logFile :="errors.log"
	errFile,err:=os.OpenFile(logFile,os.O_CREATE|os.O_WRONLY|os.O_APPEND,0666)
	if err!=nil{
		log.Fatalln("打开日志文件失败：",err)
	}

	//Info = log.New(errFile,"Info:",log.Ldate | log.Ltime | log.Lshortfile)
	InfoLogger = log.New(os.Stdout,"Info:",log.Ldate | log.Ltime | log.Lshortfile)
	WarningLogger = log.New(os.Stdout,"Warning:",log.Ldate | log.Ltime | log.Lshortfile)
	ErrorLogger = log.New(io.MultiWriter(os.Stderr,errFile),"Error:",log.Ldate | log.Ltime | log.Lshortfile)
}

func Info(v ...interface{}) {
	InfoLogger.Println(v)
}
func Warning(v ...interface{}) {
	WarningLogger.Println(v)
}
func Error(v ...interface{}) {
	ErrorLogger.Println(v)
}