package system

import (
	"github.com/joho/godotenv"
)

func LoadEnv() error {
    // 1. 加载 .env 到环境变量
    if err := godotenv.Load(); err != nil {
        return err
    }
    return nil
}