package controller

import (
	"github.com/labstack/echo/v4"
	"github.com/spf13/afero"
	"io"
	MyOS "myagent/src/core/myos"
	MyWeb "myagent/src/core/web"
	"net/http"
	"os"
)

/**
Contrller 写法
1、定义一个结构体(类),其中包含(继承)BaseController
2、实现接口 func (this *FileController) ControllerInit(myWeb *MyWeb.EchoServer) ...
3、controller中的 对外服务方法写法如下  fun (this FileController) DoSomething() ...
   在初始化网络部分服务时会反射这种写法的方法进行执行,服务方法不需要形参
4、将执行web 接口服务在 func (this *FileController) ControllerInit(myWeb *MyWeb.EchoServer) ... 中注册
5、如果非网络服务方法：1、不建议在这里编写；2、可以使用 fun (this *FileController) DoSomething() ... 形式进行编写
*/

//url 常量
const (
	FileController__ListFile     = "/FileController/ListFiles"
	FileController__CopyFiles    = "/FileController/CopyFiles"
	FileController__UploadFile   = "/FileController/UploadFile"
	FileController__DownloadFile = "/FileController/DownloadFile"
	FileController__RenameFile   = "/FileController/RenameFile"
)

type FileController struct {
	BaseController
}

//实现接口
func (this *FileController) ControllerInit(myWeb *MyWeb.EchoServer) {
	this.BaseController.ControllerInit(myWeb)

	//将所需服务在这里自行维护
	this.ListFiles()
	this.CopyFiles()
	this.UploadFile()
	this.DownloadFile()
	this.RenameFile()
	//this.MoveFile()
}

//获取文件列表
func (this FileController) ListFiles() {
	this.echoServer.GetEcho().GET(FileController__ListFile, func(context echo.Context) error {
		result := MyOS.SystemFileResult{}
		param := this.GetParamFormEchoContext(context, "filePath", 0)
		result.ListFiles(param)
		return context.JSON(http.StatusOK, this.Success(result))
	})
}

/*
复制文件、层级
形参：
	可以是文件层级
	可以是文件，如果是文件需要指定被复制文件的文件名+文件类型如：
						srcFile: /root/Hello.txt
						dstFile: /home/hello.txt
*/
func (this FileController) CopyFiles() {
	this.echoServer.GetEcho().GET(FileController__CopyFiles, func(context echo.Context) error {
		srcFile := this.GetParamFormEchoContext(context, "srcFile", 0)
		dstFile := this.GetParamFormEchoContext(context, "dstFile", 0)
		nullParam, b := this.NullParam(srcFile, dstFile)
		if !b {
			return context.JSON(http.StatusOK, this.Failed(nullParam, nil))
		}

		result := MyOS.Copy(afero.NewOsFs(), srcFile, dstFile)
		if nil != result {
			return context.JSON(http.StatusOK, this.Failed("操作失败", nil))
		} else {
			return context.JSON(http.StatusOK, this.Success("操作成功"))
		}
	})
}

func (this FileController) UploadFile() {
	this.echoServer.GetEcho().POST(FileController__UploadFile, func(context echo.Context) error {
		//配置文件中的上传地址
		file, err := context.FormFile("file")

		if err != nil {
			return err
		}
		src, err := file.Open()
		if err != nil {
			return err
		}
		defer src.Close()

		// Destination
		dst, err := os.Create(this.GetConfigUploadFilePath() + file.Filename)
		if err != nil {
			return err
		}
		defer dst.Close()

		// Copy
		if _, err = io.Copy(dst, src); err != nil {
			return err
		}

		return context.JSON(http.StatusOK, this.Success("操作成功"))
	})
}

func (this FileController) DownloadFile() {
	this.echoServer.GetEcho().POST(FileController__DownloadFile, func(context echo.Context) error {
		filePath := this.GetParamFormEchoContext(context, "filePath", 0)
		fileInfo, err := os.Stat(filePath)
		if fileInfo == nil {
			return context.JSON(http.StatusOK, this.Failed("文件不存在", nil))
		}
		if err != nil {
			return context.JSON(http.StatusOK, this.Failed("出错", nil))
		}
		return context.File(filePath)
	})
}

//文件重命名或文件移动
func (this FileController) RenameFile() {
	this.echoServer.GetEcho().GET(FileController__RenameFile, func(context echo.Context) error {
		srcFile := this.GetParamFormEchoContext(context, "srcFile", 0)
		dstFile := this.GetParamFormEchoContext(context, "dstFile", 0)

		fileInfo, err := os.Stat(srcFile)
		if err != nil {
			return context.JSON(http.StatusOK, this.Failed("出错", nil))
		}
		if fileInfo == nil {
			return context.JSON(http.StatusOK, this.Failed("文件不存在", nil))
		}

		err = os.Rename(srcFile, dstFile)
		if err != nil {
			return context.JSON(http.StatusOK, this.Failed("操作失败", nil))
		}

		return context.JSON(http.StatusOK, this.Success("操作成功"))
	})
}
