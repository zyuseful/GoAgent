package myos

import (
	"fmt"
	MyLog "myagent/src/core/log"
	"os/exec"
	"strings"
	"time"
)

//继承自BaseSystemResult 的 命令处理返回结构体
type SystemCommandResult struct {
	BaseSystemResult
	Result interface{}
}

/**
对外调用
*/
func (sysDispose *SystemCommandResult) DoCommandWithCmdStr(osTool string, cmd string) {
	defaultSystemCommand(sysDispose)
	sysDispose.Dispose(osTool, cmd)
}

func (sysDispose *SystemCommandResult) Dispose(osTool string, cmd string) {
	if len(cmd) <= 0 {
		sysDispose.Fault("")
		return
	}

	cmd = strings.Trim(cmd, " ")
	//mac 测试 ： 复合命令测试通过 ls ~ && echo 'asdf' > ~/asdf.lg
	cmdDispose(sysDispose, osTool, cmd)

	//复合命令无法通过 单命令OK
	//dispose2(dispose, osTool, cmd)

}

//Private methods
func defaultSystemCommand(result *SystemCommandResult) {
	result.BaseSystemResult.SetSystemToolSync(true)
	result.BaseSystemResult.SetSystemToolType(DISPOSE_CMD_TYPE)
	result.BaseSystemResult.SetSystemToolStartTime(time.Now())
}

func successCmd(dispose *SystemCommandResult, reason interface{}) *SystemCommandResult {
	dispose.OK("")
	dispose.Result = reason
	return dispose
}

func cmdDispose(dispose *SystemCommandResult, osTool string, cmd string) *SystemCommandResult {
	out, err := exec.Command(osTool, "-c", cmd).Output()
	if err != nil {
		dispose.Fault(err.Error())
		MyLog.Error(err.Error())
		return dispose
	}
	fmt.Println(string(out))
	return successCmd(dispose, string(out))
}

/*
func cmdDispose(dispose *SystemCommandResult, osTool string, cmd string) *SystemCommandResult{
	var outInfo bytes.Buffer
	c := exec.Command(osTool, "-c", cmd)
	// 设置接收
	c.Stdout = &outInfo
	// 执行
	c.Run()

	fmt.Println(outInfo.String())
	return successCmd(dispose, outInfo.String())
}
*/

func dispose2(dispose *SystemCommandResult, tool string, cmd string) {
	splitIndex := strings.Index(cmd, " ")
	if splitIndex <= 1 {
		dispose.Fault("命令有误")
		return
	}

	MyLog.Info("cmd1 ", string(cmd[0:splitIndex]))
	MyLog.Info("cmd2 ", string(cmd[splitIndex+1:]))

	var c1 string = strings.Trim(string(cmd[0:splitIndex]), " ")
	var c2 string = strings.Trim(string(cmd[splitIndex:]), " ")

	out, err := exec.Command(c1, c2).Output()
	if err != nil {
		dispose.Fault(err.Error())
		MyLog.Error(err.Error())
		return
	}

	successCmd(dispose, string(out))
}
