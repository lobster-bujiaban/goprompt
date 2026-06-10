package prompt

// Message LLM 对话消息
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// Config 配置
type Config struct {
	Roles map[string]RoleConfig `json:"roles"`
}

// RoleConfig 角色配置
type RoleConfig struct {
	Temperature float64 `json:"temperature"`
	MaxTokens   int     `json:"max_tokens"`
}

// PromptManager 管理 Prompt 模板和角色配置
type PromptManager struct {
	engine *PromptEngine
	config *Config
}

func NewPromptManager(configPath string) (*PromptManager, error) {
	engine := New()

	// 注册基础模板
	engine.AddTemplate("base-system", `{{block "role" .}}你是 AI 编程助手。{{end}}

通用规则：
1. 不确定的 API 直接说「不确定」，不要编造
2. 代码完整可编译，包含必要的 import
3. 中文回答，代码注释也用中文

{{block "specialty" .}}{{end}}

输出约束：
- 代码用 {{block "lang" .}}Go{{end}} 代码块
- 不要输出额外的解释文字{{block "extra" .}}{{end}}`)

	// 注册代码审查模板
	engine.AddTemplate("code-review", `你是代码审查专家。审查以下 {{.Language}} 代码：

{{codeblock .Language .Code}}

请重点关注以下方面：
{{if .CheckSecurity}}
- SQL 注入、XSS、命令注入等安全漏洞
{{end}}
{{if .CheckPerformance}}
- 性能瓶颈：不必要的内存分配、低效的循环、缺少并发
{{end}}
{{if .CheckStyle}}
- 代码风格：命名规范、注释完整性、函数长度
{{end}}
{{if .CheckConcurrency}}
- 并发问题：数据竞争、死锁、goroutine 泄漏
{{end}}
{{if not .CheckSecurity and not .CheckPerformance and not .CheckStyle and not .CheckConcurrency}}
- 通用代码质量问题
{{end}}

输出格式：{{.OutputFormat}}`)

	// 注册代码生成模板
	engine.AddTemplate("code-gen", `你是 {{.Role}}，擅长 {{.Language}} 开发。

任务：{{.Task}}

要求：
- 语言：{{.Language}}
- 框架：{{.Framework}}
- 输出：{{if .WithTests}}包含单元测试{{else}}只输出业务代码{{end}}`)

	// 注册数据分析模板
	engine.AddTemplate("data-analysis", `你是一个数据分析师。分析以下数据：

数据源：
{{range .DataSources}}
- {{.Name}}：{{.Description}}（时间范围：{{.TimeRange}}）
{{end}}

分析指标（共 {{len .Metrics}} 个）：
{{range $i, $m := .Metrics}}
指标 {{add $i 1}}：{{$m.Name}}
  定义：{{$m.Definition}}
  目标值：{{$m.Target}}
  {{if $m.Alert}}⚠️ 如果低于 {{$m.Alert}} 需要重点关注{{end}}
{{end}}

请给出：
1. 每个指标的趋势判断（上升/下降/持平）
2. 异常指标的根因分析
3. 改进建议`)

	// 注册审查者模板（继承基类）
	engine.AddTemplate("reviewer-system", `{{template "base-system" .}}

{{define "role"}}你是代码审查员，专门审查 Go 代码的安全和性能问题。{{end}}

{{define "specialty"}}
审查规则：
- 高危问题 > 中危问题 > 低危问题
- 每个问题标注行号
- 提供具体的修复代码
{{end}}`)

	// 注册文档写手模板
	engine.AddTemplate("doc-writer-system", `{{template "base-system" .}}

{{define "role"}}你是技术文档写手，擅长写 Go godoc 风格注释。{{end}}

{{define "specialty"}}
写作规则：
- 每行注释不超过 80 字符
- 参数说明用「参数名: 说明」格式
- 返回值说明用「返回: 说明」格式
{{end}}`)

	return &PromptManager{engine: engine}, nil
}

// BuildPrompt 构建 System + User Prompt
func (pm *PromptManager) BuildPrompt(roleKey, templateName string, data interface{}) ([]Message, error) {
	systemPrompt, err := pm.engine.Render("system-"+templateName, data)
	if err != nil {
		return nil, err
	}

	userPrompt, err := pm.engine.Render("user-"+templateName, data)
	if err != nil {
		return nil, err
	}

	return []Message{
		{Role: "system", Content: systemPrompt},
		{Role: "user", Content: userPrompt},
	}, nil
}

// Render 直接渲染指定模板
func (pm *PromptManager) Render(name string, data interface{}) (string, error) {
	return pm.engine.Render(name, data)
}
