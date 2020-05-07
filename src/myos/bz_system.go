package myos

import "time"

//------------------------------------------------------------
//OS -- 后期考虑 系统不同
type OSType = uint8

const (
	//windows 0 ~ 49
	WIN_OS OSType = 0

	//linux 50 ~ 200
	LINUX_OS OSType = 50
)

//------------------------------------------------------------

//============================================================
//OS Cmd Tool 
type OSTool = string

const (
	//windows 0 ~ 49
	WIN_TOOL_BAT OSTool = "bat"

	//linux 50 ~ 200
	LINUX_TOOL_BASH OSTool = "bash"
)

//============================================================

//------------------------------------------------------------
//OS System Result Tool Type
type SYSTEM_DISPOSE_TYPE uint8

const (
	DISPOSE_FILE_TYPE SYSTEM_DISPOSE_TYPE = 10
	DISPOSE_CMD_TYPE  SYSTEM_DISPOSE_TYPE = 20
)

//OS System Result
type BaseSystemResult struct {
	Type      SYSTEM_DISPOSE_TYPE
	Status    bool
	StartTime time.Time
	EndTime   time.Time
	Sync      bool
	Mesg      string
	cmd       string
}

/**
设置工作类型
*/
func (baseSystemResult *BaseSystemResult) SetSystemToolType(systemDisposeType SYSTEM_DISPOSE_TYPE) {
	baseSystemResult.Type = systemDisposeType
}

/**
设置工作开始时间
*/
func (baseSystemResult *BaseSystemResult) SetSystemToolStartTime(startTime time.Time) {
	baseSystemResult.StartTime = startTime
}

/**
设置工作开始时间
*/
func (baseSystemResult *BaseSystemResult) SetSystemToolEndTime(endTime time.Time) {
	baseSystemResult.EndTime = endTime
}

/**
设置工作sync
*/
func (baseSystemResult *BaseSystemResult) SetSystemToolSync(sync bool) {
	baseSystemResult.Sync = sync
}

/**
失败设置
*/
func (baseSystemResult *BaseSystemResult) Fault(reason string) {
	baseSystemResult.Status = false
	if len(reason) == 0 {
		baseSystemResult.Mesg = "命令为空"
	} else {
		baseSystemResult.Mesg = reason
	}
}

/**
成功设置
*/
func (baseSystemResult *BaseSystemResult) OK(reason string) {
	baseSystemResult.Status = true
	baseSystemResult.EndTime = time.Now()
	if len(reason) == 0 {
		baseSystemResult.Mesg = "执行成功"
	} else {
		baseSystemResult.Mesg = reason
	}
}

//------------------------------------------------------------
