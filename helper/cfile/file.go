package cfile

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// GetFileContentType 获取文件内容的类型os.open()后调用
func GetFileContentType(out *os.File) (string, error) {
	// 只需前512 个字节即可
	buffer := make([]byte, 512)
	_, err := out.Read(buffer)
	if err != nil {
		return "", err
	}
	t := http.DetectContentType(buffer)
	return t, nil
}

type dirEntry struct {
	Label    string     `json:"label"`
	Value    string     `json:"value"`
	Size     int64      `json:"size"`
	ModTime  time.Time  `json:"modtime"`
	IsDir    bool       `json:"isDir"`
	Children []dirEntry `json:"children,omitempty"`
}

// ReadDirTree 递归获取目录树
func ReadDirTree(dirPath string) ([]dirEntry, error) {
	var treeData = make([]dirEntry, 0)
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return treeData, err
	}
	for _, f := range files {
		fileInfo, err := f.Info()
		if err != nil {
			if errors.Is(err, fs.ErrNotExist) { // 访问过程中目录被删除，则跳过此次
				continue
			} else {
				return treeData, err
			}
		}
		tmp := dirEntry{
			Label:    f.Name(),
			Value:    filepath.Join(dirPath, f.Name()),
			Size:     fileInfo.Size(),
			ModTime:  fileInfo.ModTime(),
			IsDir:    f.IsDir(),
			Children: nil,
		}
		if f.IsDir() {
			child, err := ReadDirTree(filepath.Join(dirPath, f.Name()))
			if err != nil {
				return treeData, err
			}
			if len(child) > 0 {
				tmp.Children = child
			}
		}
		treeData = append(treeData, tmp)
	}
	return treeData, nil
}

// ReadDirFiles 获取目录下所有文件
func ReadDirFiles(dirPath string) ([]dirEntry, error) {
	var treeData = make([]dirEntry, 0)
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return treeData, err
	}
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		fileInfo, err := f.Info()
		if err != nil {
			if errors.Is(err, fs.ErrNotExist) { // 访问过程中目录被删除，则跳过此次
				continue
			} else {
				return treeData, err
			}
		}
		treeData = append(treeData, dirEntry{
			Label:    f.Name(),
			Value:    filepath.Join(dirPath, f.Name()),
			Size:     fileInfo.Size(),
			ModTime:  fileInfo.ModTime(),
			IsDir:    f.IsDir(),
			Children: nil,
		})
	}
	return treeData, nil
}

// FormatFileSize 字节的单位转换 保留两位小数
func FormatFileSize(fileSize int64) (size string) {
	if fileSize < 1024 {
		// return strconv.FormatInt(fileSize, 10) + "B"
		return fmt.Sprintf("%.2fB", float64(fileSize)/float64(1))
	} else if fileSize < (1024 * 1024) {
		return fmt.Sprintf("%.2fKB", float64(fileSize)/float64(1024))
	} else if fileSize < (1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2fMB", float64(fileSize)/float64(1024*1024))
	} else if fileSize < (1024 * 1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2fGB", float64(fileSize)/float64(1024*1024*1024))
	} else if fileSize < (1024 * 1024 * 1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2fTB", float64(fileSize)/float64(1024*1024*1024*1024))
	} else { // if fileSize < (1024 * 1024 * 1024 * 1024 * 1024 * 1024)
		return fmt.Sprintf("%.2fEB", float64(fileSize)/float64(1024*1024*1024*1024*1024))
	}
}

// Copy 复制文件
func Copy(src, dst string) (int64, error) {
	// dst 要过虑非正常路径，如&之类的字符
	filterStr := []string{"..", "&", ":", ";", "|", "$", "%", "?", "\r", "\n", "`", ","}
	for _, str := range filterStr {
		if strings.Contains(dst, str) {
			return 0, errors.New("目标路径存在特殊字符: " + str)
		}
	}
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}
	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}
	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()
	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}

// IsDir 判断文件夹是否存在
func IsDir(path string) bool {
	f, err := os.Stat(path)
	if err != nil {
		return false
	}
	return f.IsDir()
}

// IsFile 判断文件是否存在
func IsFile(path string) bool {
	f, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !f.IsDir()
}

// FileExist 判断文件是否存在
func FileExist(fileName string) bool {
	_, err := os.Stat(fileName)
	return err == nil || os.IsExist(err)
}

// Mkdir 创建目录
func Mkdir(path string, perm os.FileMode) error {
	if path == "" {
		return errors.New("path empty")
	}
	// recordID 要过虑非正常路径，如&之类的字符
	filterStr := []string{"..", "&", ":", ";", "|", "$", "%", "?", "\r", "\n", "`", ",", " "}
	for _, str := range filterStr {
		if strings.Contains(path, str) {
			return errors.New("目标路径存在特殊字符: " + str)
		}
	}
	err := os.MkdirAll(path, perm)
	return err
}

// WriteFile 写文件
func WriteFile(fileName string, content string) error {
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	// 写入字符串到文件中
	_, err = io.WriteString(file, content)
	return err
}
