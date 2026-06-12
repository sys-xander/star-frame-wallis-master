package sls

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/aliyun/aliyun-log-go-sdk/producer"
	sls "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/zeromicro/go-zero/core/logx"
)

// Config 阿里云 SLS 配置
type Config struct {
	Enabled         bool   `json:"enabled"`
	Endpoint        string `json:"endpoint"`
	Project         string `json:"project"`
	LogStore        string `json:"logstore"`
	AccessKeyID     string `json:"accessKeyID"`
	AccessKeySecret string `json:"accessKeySecret"`
	Topic           string `json:"topic"`
}

// Writer 实现 go-zero 的 logx.Writer 接口，异步写入 SLS
type Writer struct {
	cfg          Config
	producer     *producer.Producer
	fallback     *os.File
	fallbackPath string
	mu           sync.Mutex
	closed       bool
}

// NewWriter 创建 SLS 日志写入器
func NewWriter(cfg Config) (*Writer, error) {
	if !cfg.Enabled {
		return nil, nil
	}

	producerConfig := producer.GetDefaultProducerConfig()
	producerConfig.Endpoint = cfg.Endpoint
	producerConfig.AccessKeyID = cfg.AccessKeyID
	producerConfig.AccessKeySecret = cfg.AccessKeySecret

	p, err := producer.NewProducer(producerConfig)
	if err != nil {
		return nil, fmt.Errorf("init producer failed: %w", err)
	}
	p.Start()

	// 创建本地降级日志文件
	fallbackPath := fmt.Sprintf("logs/sls-fallback-%s.log", time.Now().Format("2006-01-02"))
	_ = os.MkdirAll("logs", 0755)
	fallback, err := os.OpenFile(fallbackPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		p.Close(1000)
		return nil, fmt.Errorf("open fallback log file failed: %w", err)
	}

	return &Writer{
		cfg:          cfg,
		producer:     p,
		fallback:     fallback,
		fallbackPath: fallbackPath,
	}, nil
}

// Write 实现 io.Writer 接口，接收 logx 输出的 JSON 日志
func (w *Writer) Write(p []byte) (n int, err error) {
	if w == nil || w.closed {
		return len(p), nil
	}

	var logEntry map[string]any
	if err := json.Unmarshal(p, &logEntry); err != nil {
		w.writeFallback(p)
		return len(p), nil
	}

	level := "info"
	if l, ok := logEntry["level"].(string); ok {
		level = l
	}
	content := ""
	if c, ok := logEntry["content"].(string); ok {
		content = c
	}

	now := uint32(time.Now().Unix())
	log := &sls.Log{
		Time:     &now,
		Contents: make([]*sls.LogContent, 0),
	}

	log.Contents = append(log.Contents, newLogContent("level", level))
	log.Contents = append(log.Contents, newLogContent("content", content))
	log.Contents = append(log.Contents, newLogContent("raw", string(p)))

	for k, v := range logEntry {
		if k == "level" || k == "content" || k == "timestamp" {
			continue
		}
		strVal := fmt.Sprintf("%v", v)
		log.Contents = append(log.Contents, newLogContent(k, strVal))
	}

	w.sendLog(log)

	return len(p), nil
}

func newLogContent(k, v string) *sls.LogContent {
	return &sls.LogContent{Key: &k, Value: &v}
}

func (w *Writer) sendLog(log *sls.Log) {
	defer func() {
		if r := recover(); r != nil {
		}
	}()

	err := w.producer.SendLog(w.cfg.Project, w.cfg.LogStore, w.cfg.Topic, "", log)
	if err != nil {
		data, _ := json.Marshal(log)
		w.writeFallback(data)
	}
}

func (w *Writer) writeFallback(p []byte) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.fallback != nil {
		w.fallback.Write(p)
		w.fallback.Write([]byte("\n"))
	}
}

// Alert 实现 logx.Writer 接口
func (w *Writer) Alert(v any) {
	w.writeLevelLog("alert", v)
}

// Close 实现 logx.Writer 接口
func (w *Writer) Close() error {
	if w == nil || w.closed {
		return nil
	}
	w.mu.Lock()
	w.closed = true
	w.mu.Unlock()

	if w.producer != nil {
		w.producer.Close(3000)
	}
	if w.fallback != nil {
		w.fallback.Close()
	}
	return nil
}

// Debug 实现 logx.Writer 接口
func (w *Writer) Debug(v any, fields ...logx.LogField) {
	w.writeLevelLogWithFields("debug", v, fields...)
}

// Error 实现 logx.Writer 接口
func (w *Writer) Error(v any, fields ...logx.LogField) {
	w.writeLevelLogWithFields("error", v, fields...)
}

// Info 实现 logx.Writer 接口
func (w *Writer) Info(v any, fields ...logx.LogField) {
	w.writeLevelLogWithFields("info", v, fields...)
}

// Severe 实现 logx.Writer 接口
func (w *Writer) Severe(v any) {
	w.writeLevelLog("severe", v)
}

// Slow 实现 logx.Writer 接口
func (w *Writer) Slow(v any, fields ...logx.LogField) {
	w.writeLevelLogWithFields("slow", v, fields...)
}

// Stack 实现 logx.Writer 接口
func (w *Writer) Stack(v any) {
	w.writeLevelLog("stack", v)
}

// Stat 实现 logx.Writer 接口
func (w *Writer) Stat(v any, fields ...logx.LogField) {
	w.writeLevelLogWithFields("stat", v, fields...)
}

func (w *Writer) writeLevelLog(level string, content any) {
	w.writeLevelLogWithFields(level, content)
}

func (w *Writer) writeLevelLogWithFields(level string, content any, fields ...logx.LogField) {
	logEntry := map[string]any{
		"level":   level,
		"content": content,
		"ts":      time.Now().Format(time.RFC3339),
	}
	for _, f := range fields {
		logEntry[f.Key] = f.Value
	}
	data, _ := json.Marshal(logEntry)
	w.Write(data)
}