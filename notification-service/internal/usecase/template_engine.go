package usecase

import (
    "bytes"
    "fmt"
    "html/template"
    "path/filepath"
    "sync"

    "github.com/Temutjin2k/Tyndau/notification-service/internal/model"
    "github.com/Temutjin2k/Tyndau/notification-service/pkg/logger"
)

// GoTemplateEngine реализует интерфейс TemplateEngine с использованием Go templates
type GoTemplateEngine struct {
    templatesDir string
    templates    map[string]*template.Template
    mu           sync.RWMutex
    logger       *logger.Logger
}

// NewGoTemplateEngine создает новый template engine
func NewGoTemplateEngine(templatesDir string, logger *logger.Logger) *GoTemplateEngine {
    return &GoTemplateEngine{
        templatesDir: templatesDir,
        templates:    make(map[string]*template.Template),
        logger:       logger,
    }
}

// RenderTemplate рендерит шаблон с указанными данными
func (e *GoTemplateEngine) RenderTemplate(templateName string, data interface{}) (string, error) {
    // Проверяем, загружен ли шаблон
    e.mu.RLock()
    tmpl, ok := e.templates[templateName]
    e.mu.RUnlock()
    
    if !ok {
        // Шаблон не загружен, загружаем его
        templatePath := filepath.Join(e.templatesDir, templateName+".html")
        var err error
        tmpl, err = template.ParseFiles(templatePath)
        if err != nil {
            return "", fmt.Errorf("failed to parse template %s: %w", templateName, err)
        }
        
        // Сохраняем шаблон
        e.mu.Lock()
        e.templates[templateName] = tmpl
        e.mu.Unlock()
    }
    
    // Рендерим шаблон
    var buf bytes.Buffer
    if err := tmpl.Execute(&buf, data); err != nil {
        return "", fmt.Errorf("failed to execute template %s: %w", templateName, err)
    }
    
    return buf.String(), nil
}