package myos

/**
系统操作接口
 */
type SysDispose interface {
	//设置具体处理的类型 如文件类型处理、命令类型处理
	DisposeSysToolType(systemDisposeType SYSTEM_DISPOSE_TYPE)
}
