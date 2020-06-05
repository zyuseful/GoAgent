package myos

import (
	"errors"
	"github.com/spf13/afero"
	"io"
	"io/ioutil"
	MyLog "myagent/src/core/log"
	"os"
	"path"
	"path/filepath"
	"time"
)

//继承自BaseSystemResult 的 文件处理返回结构体
type SystemFileResult struct {
	BaseSystemResult
	//Result  []os.FileInfo
	Result []SystemFile
}

type SystemFile struct {
	IsDir      bool
	FileName   string
	FileSize   int64
	FilePath   string
	UpdateTime time.Time
}

/**
根据文件路径查看文件列表
*/
func (systemFile *SystemFileResult) ListFiles(filePath string) {
	systemFile.defaultSystemFileResult()

	dir, err := ioutil.ReadDir(filePath)
	if err != nil {
		systemFile.Fault(err.Error())
		MyLog.Error(err)
		return
	}
	MyLog.Info(len(dir))

	//以下两种方法均可
	//files := *new([]SystemFile)
	files := make([]SystemFile, 0, len(dir))

	for _, df := range dir {
		if df != nil || len(df.Name()) > 0 || df.Name() != "." || df.Name() != ".." {
			var systemFileJsonObj SystemFile = SystemFile{
				IsDir:      df.IsDir(),
				FileName:   df.Name(),
				FileSize:   df.Size(),
				UpdateTime: df.ModTime(),
				FilePath:   path.Join(filePath, df.Name()),
			}
			files = append(files, systemFileJsonObj)
		} else {
			continue
		}
	}
	systemFile.Result = files
	systemFile.BaseSystemResult.SetSystemToolEndTime(time.Now())
	systemFile.OK("")

}

func (result *SystemFileResult)defaultSystemFileResult() {
	result.BaseSystemResult.SetSystemToolSync(true)
	result.BaseSystemResult.SetSystemToolType(DISPOSE_FILE_TYPE)
	result.BaseSystemResult.SetSystemToolStartTime(time.Now())
}

// Copy copies a file or folder from one place to another.
func Copy(fs afero.Fs, src, dst string) error {
	if src = path.Clean("/" + src); src == "" {
		return os.ErrNotExist
	}

	if dst = path.Clean("/" + dst); dst == "" {
		return os.ErrNotExist
	}

	if src == "/" || dst == "/" {
		// Prohibit copying from or to the virtual root directory.
		return os.ErrInvalid
	}

	if dst == src {
		return os.ErrInvalid
	}

	info, err := fs.Stat(src)
	if err != nil {
		return err
	}

	if info.IsDir() {
		return CopyDir(fs, src, dst)
	}

	return CopyFile(fs, src, dst)
}

func CopyFile(fs afero.Fs, source string, dest string) error {
	// Open the source file.
	src, err := fs.Open(source)
	if err != nil {
		return err
	}
	defer src.Close()

	// Makes the directory needed to create the dst
	// file.
	err = fs.MkdirAll(filepath.Dir(dest), 0666)
	if err != nil {
		return err
	}

	// Create the destination file.
	dst, err := fs.Create(dest)
	if err != nil {
		return err
	}
	defer dst.Close()

	// Copy the contents of the file.
	_, err = io.Copy(dst, src)
	if err != nil {
		return err
	}

	// Copy the mode if the user can't
	// open the file.
	info, err := fs.Stat(source)
	if err != nil {
		err = fs.Chmod(dest, info.Mode())
		if err != nil {
			return err
		}
	}

	return nil
}

func CopyDir(fs afero.Fs, source string, dest string) error {
	// Get properties of source.
	srcinfo, err := fs.Stat(source)
	if err != nil {
		return err
	}

	// Create the destination directory.
	err = fs.MkdirAll(dest, srcinfo.Mode())
	if err != nil {
		return err
	}

	dir, _ := fs.Open(source)
	obs, err := dir.Readdir(-1)
	if err != nil {
		return err
	}

	var errs []error

	for _, obj := range obs {
		fsource := source + "/" + obj.Name()
		fdest := dest + "/" + obj.Name()

		if obj.IsDir() {
			// Create sub-directories, recursively.
			err = CopyDir(fs, fsource, fdest)
			if err != nil {
				errs = append(errs, err)
			}
		} else {
			// Perform the file copy.
			err = CopyFile(fs, fsource, fdest)
			if err != nil {
				errs = append(errs, err)
			}
		}
	}

	var errString string
	for _, err := range errs {
		errString += err.Error() + "\n"
	}

	if errString != "" {
		return errors.New(errString)
	}

	return nil
}