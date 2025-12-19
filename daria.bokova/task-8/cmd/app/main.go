package main

import (
    "fmt"
    "os"
    "task-8/internal/settings"
)

func main() {
    // Проверяем, инициализирована ли функция
    if settings.GetConfig == nil {
        fmt.Fprintln(os.Stderr, "ERROR: GetConfig not initialized")
        os.Exit(1)
    }
    
    configData := settings.GetConfig()
    
    if len(configData) == 0 {
        fmt.Fprintln(os.Stderr, "ERROR: Empty configuration data")
        os.Exit(1)
    }
    
    cfg, err := settings.ParseConfig(configData)
    if err != nil {
        fmt.Fprintf(os.Stderr, "ERROR parsing config: %v\n", err)
        os.Exit(1)
    }
    
    fmt.Printf("%s %s\n", cfg.Env, cfg.Logging)
}