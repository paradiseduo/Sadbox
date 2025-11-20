package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	// 解析命令行参数
	deleteFileName := flag.String("delete", "", "要删除的文件名（会删除包含该文件的容器目录），可以指定多个文件名，用空格分隔")
	showSystem := flag.Bool("system", false, "显示以 com.apple. 开头的系统文件")
	flag.Parse()

	// 获取用户主目录
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("无法获取用户主目录: %v\n", err)
		return
	}

	// 构建目标目录路径
	containersDir := filepath.Join(homeDir, "Library", "Containers")

	// 检查目录是否存在
	if _, err := os.Stat(containersDir); os.IsNotExist(err) {
		fmt.Printf("目录不存在: %s\n", containersDir)
		return
	}

	// 目标子路径
	targetSubPath := filepath.Join("Data", "Library", "Application Scripts")

	// 如果有 delete 参数，执行删除操作
	if *deleteFileName != "" {
		// 按空格分割多个文件名
		fileNames := strings.Fields(*deleteFileName)
		err = deleteContainersByFileNames(containersDir, targetSubPath, fileNames)
		if err != nil {
			fmt.Printf("删除操作失败: %v\n", err)
		}
		return
	}

	// 否则，正常遍历并输出所有路径
	err = filepath.Walk(containersDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			// 如果无法访问某个目录，继续遍历其他目录
			return nil
		}

		// 检查是否是目标文件夹
		if info.IsDir() {
			// 检查路径是否以目标子路径结尾
			normalizedPath := filepath.ToSlash(path)
			normalizedTarget := filepath.ToSlash(targetSubPath)
			if strings.HasSuffix(normalizedPath, normalizedTarget) {
				// 读取该文件夹中的唯一文件
				printSingleFile(path, *showSystem)
			}
		}

		return nil
	})

	if err != nil {
		fmt.Printf("遍历目录时出错: %v\n", err)
	}
}

// printSingleFile 打印指定目录中唯一文件的完整路径
func printSingleFile(dirPath string, showSystem bool) {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return
	}

	// 只处理第一个文件/目录（因为目录下只有一个文件）
	if len(entries) > 0 {
		fileName := entries[0].Name()

		// 如果文件名以 com.apple. 开头，且未指定 -system 参数，则跳过
		if strings.HasPrefix(fileName, "com.apple.") && !showSystem {
			return
		}

		filePath := filepath.Join(dirPath, fileName)
		fmt.Println(filePath)
	}
}

// deleteContainersByFileNames 根据多个文件名查找并删除包含这些文件的容器目录
func deleteContainersByFileNames(containersDir, targetSubPath string, fileNames []string) error {
	var errors []string
	successCount := 0

	for _, fileName := range fileNames {
		fileName = strings.TrimSpace(fileName)
		if fileName == "" {
			continue
		}

		fmt.Printf("\n处理文件名: %s\n", fileName)
		err := deleteContainerByFileName(containersDir, targetSubPath, fileName)
		if err != nil {
			errors = append(errors, fmt.Sprintf("  %s: %v", fileName, err))
		} else {
			successCount++
		}
	}

	// 输出总结
	fmt.Printf("\n删除操作完成: 成功 %d 个", successCount)
	if len(errors) > 0 {
		fmt.Printf(", 失败 %d 个\n", len(errors))
		fmt.Println("\n失败详情:")
		for _, errMsg := range errors {
			fmt.Println(errMsg)
		}
		return fmt.Errorf("部分删除操作失败")
	}
	fmt.Println()
	return nil
}

// deleteContainerByFileName 根据文件名查找并删除包含该文件的容器目录
func deleteContainerByFileName(containersDir, targetSubPath, fileName string) error {
	// 读取 Containers 目录下的所有子目录
	entries, err := os.ReadDir(containersDir)
	if err != nil {
		return fmt.Errorf("无法读取 Containers 目录: %v", err)
	}

	// 遍历每个容器目录
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		// 构建 Application Scripts 文件夹的完整路径
		appScriptsPath := filepath.Join(containersDir, entry.Name(), targetSubPath)

		// 检查该路径是否存在
		if _, err := os.Stat(appScriptsPath); os.IsNotExist(err) {
			continue
		}

		// 读取 Application Scripts 文件夹中的文件
		appScriptsEntries, err := os.ReadDir(appScriptsPath)
		if err != nil {
			continue
		}

		// 检查文件名是否匹配
		if len(appScriptsEntries) > 0 && appScriptsEntries[0].Name() == fileName {
			// 找到匹配的文件，获取容器目录路径
			containerDirToDelete := filepath.Join(containersDir, entry.Name())

			// 执行删除
			fmt.Printf("找到匹配的容器目录: %s\n", containerDirToDelete)
			fmt.Printf("正在删除...\n")
			err = os.RemoveAll(containerDirToDelete)
			if err != nil {
				return fmt.Errorf("删除目录失败: %v", err)
			}
			fmt.Printf("成功删除: %s\n", containerDirToDelete)
			return nil
		}
	}

	return fmt.Errorf("未找到包含文件 '%s' 的容器目录", fileName)
}
