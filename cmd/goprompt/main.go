package main

import (
	"fmt"

	"github.com/lobster-bujiaban/goprompt/prompt"
)

func main() {
	pm, err := prompt.NewPromptManager("roles.json")
	if err != nil {
		panic(err)
	}

	// 示例：代码审查
	data := map[string]interface{}{
		"Language":         "Go",
		"Code":             "func GetUser(id int) *User { return db.Query(\"SELECT * FROM users WHERE id=\" + id) }",
		"CheckSecurity":    true,
		"CheckPerformance": false,
		"CheckStyle":       false,
		"CheckConcurrency": true,
		"OutputFormat":     "每个问题一行，格式：[严重程度] 文件名:行号 - 问题描述",
	}
	result, _ := pm.Render("code-review", data)
	fmt.Println(result)
}
