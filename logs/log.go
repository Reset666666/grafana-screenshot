package logs

import (
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"log"
	"os"
)

// 记录错误日志
var ErrorLogger *log.Logger

// 记录信息日志
var InfoLogger *log.Logger

// Logs
// @Description: 定义日志的输出格式
// @param logFile *os.File "字段说明"
func Logs(logFile *os.File) {
	// MaxSize 定义了日志文件的最大大小，超过这个大小后会被滚动
	// MaxAge 定义了日志文件的最长存活时间，超过这个时间后也会被滚动
	// MaxBackups 定义了最大备份数量，超过这个数量后最旧的备份会被删除
	// LocalTime 表示使用本地时间
	// Compress 表示启用压缩
	logWriter := &lumberjack.Logger{
		Filename:   logFile.Name(),
		MaxSize:    20,  // MB
		MaxAge:     180, // days
		LocalTime:  true,
		Compress:   true,
		MaxBackups: 5,
	}
	ErrorLogger = log.New(io.MultiWriter(logWriter, os.Stderr), "ERROR: ", log.Ldate|log.Ltime|log.Lmsgprefix|log.Lshortfile) //设置日志输出格式
	InfoLogger = log.New(io.MultiWriter(logWriter, os.Stdout), "INFO: ", log.Ldate|log.Ltime|log.Lmsgprefix|log.Lshortfile)
}

// CreateLog
// @Description: 创建日志文件
func CreateLog() {
	err := os.MkdirAll("./log", 0755)
	if err != nil {
		// 处理创建文件夹失败的错误
		log.Panic("创建日志文件夹失败:", err)
	}
	// 记录日志
	logFile, err := os.OpenFile("./log/screenshot.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Panic("Failed to create log file:", err)
	}
	Logs(logFile)
}
