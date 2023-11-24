package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	// 定义命令行参数
	pathPtr := flag.String("p", "", "目标文件夹路径，多个文件夹使用逗号分割")
	outputPtr := flag.String("o", "release.txt", "输出文件路径和文件名")
	extensionsPtr := flag.String("t", "txt,md", "需要合并的文件后缀，多个后缀使用逗号分割")
	flag.Parse()

	// 检查目标文件夹路径是否为空
	if *pathPtr == "" {
		fmt.Println("请提供目标文件夹路径")
		return
	}

	// 解析文件夹路径
	paths := strings.Split(*pathPtr, ",")

	// 解析文件后缀
	extensions := strings.Split(*extensionsPtr, ",")

	// 获取目标文件夹下的所有文件
	fileList := []string{}
	for _, path := range paths {
		err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() && hasExtension(info.Name(), extensions) {
				fileList = append(fileList, path)
			}
			return nil
		})
		if err != nil {
			fmt.Printf("遍历文件夹 %s 失败: %s\n", path, err)
			continue
		}
	}

	if len(fileList) == 0 {
		fmt.Println("没有找到符合条件的文件")
		return
	}

	// 合并文件内容
	var content string
	for i, file := range fileList {
		fileContent, err := ioutil.ReadFile(file)
		if err != nil {
			fmt.Println("读取文件失败:", err)
			return
		}
		content += string(fileContent)

		// 在不同文件间添加两行空行
		if i < len(fileList)-1 {
			content += "\n\n"
		}
	}

	// 写入输出文件
	err := ioutil.WriteFile(*outputPtr, []byte(content), 0644)
	if err != nil {
		fmt.Println("写入文件失败:", err)
		return
	}

	fmt.Println("文件合并成功")
}

// 判断文件名是否具有指定的后缀
func hasExtension(filename string, extensions []string) bool {
	for _, ext := range extensions {
		if strings.HasSuffix(filename, "."+ext) {
			return true
		}
	}
	return false
}

