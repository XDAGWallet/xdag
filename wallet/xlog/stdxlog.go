package xlog

/*
   全局默认提供一个Log对外句柄，可以直接使用API系列调用
   全局日志对象 StdXdagLog
*/

import "os"

var StdXdagLog = NewXdagLog(os.Stderr, "", BitDefault)

//获取StdXdagLog 标记位
func Flags() int {
	return StdXdagLog.Flags()
}

//设置StdXdagLog标记位
func ResetFlags(flag int) {
	StdXdagLog.ResetFlags(flag)
}

//添加flag标记
func AddFlag(flag int) {
	StdXdagLog.AddFlag(flag)
}

//设置StdXdagLog 日志头前缀
func SetPrefix(prefix string) {
	StdXdagLog.SetPrefix(prefix)
}

//设置StdXdagLog绑定的日志文件
func SetLogFile(fileDir string, fileName string) {
	StdXdagLog.SetLogFile(fileDir, fileName)
}

//设置关闭debug
func CloseDebug() {
	StdXdagLog.CloseDebug()
}

//设置打开debug
func OpenDebug() {
	StdXdagLog.OpenDebug()
}

// ====> Debug <====
func Debugf(format string, v ...interface{}) {
	StdXdagLog.Debugf(format, v...)
}

func Debug(v ...interface{}) {
	StdXdagLog.Debug(v...)
}

// ====> Trace <====
func Trace(v ...interface{}) {
	StdXdagLog.Trace(v...)
}

// ====> Info <====
func Infof(format string, v ...interface{}) {
	StdXdagLog.Infof(format, v...)
}

func Info(v ...interface{}) {
	StdXdagLog.Info(v...)
}

// ====> Warn <====
func Warnf(format string, v ...interface{}) {
	StdXdagLog.Warnf(format, v...)
}

func Warn(v ...interface{}) {
	StdXdagLog.Warn(v...)
}

// ====> Error <====
func Errorf(format string, v ...interface{}) {
	StdXdagLog.Errorf(format, v...)
}

func Error(v ...interface{}) {
	StdXdagLog.Error(v...)
}

// ====> Fatal 需要终止程序 <====
func Fatalf(format string, v ...interface{}) {
	StdXdagLog.Fatalf(format, v...)
}

func Fatal(v ...interface{}) {
	StdXdagLog.Fatal(v...)
}

// ====> Panic  <====
func Panicf(format string, v ...interface{}) {
	StdXdagLog.Panicf(format, v...)
}

func Panic(v ...interface{}) {
	StdXdagLog.Panic(v...)
}

// ====> Stack  <====
func Stack(v ...interface{}) {
	StdXdagLog.Stack(v...)
}

func init() {
	//因为StdXdagLog对象 对所有输出方法做了一层包裹，所以在打印调用函数的时候，比正常的logger对象多一层调用
	//一般的zinxLogger对象 calldDepth=2, StdXdagLog的calldDepth=3
	StdXdagLog.calldDepth = 3
}
