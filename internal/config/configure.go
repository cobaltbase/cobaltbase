package config

import (
	"errors"
	"fmt"
	"github.com/cobaltbase/cobaltbase/internal/ct"
	"github.com/cobaltbase/cobaltbase/internal/utils"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
)

var SMTPConfig ct.SMTPConfig

var DB *gorm.DB

var Validate = validator.New()

func loadConfig() (Config, error) {
	var config Config

	// Load environment variables from .env file, if it exists
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, falling back to environment variables")
	}

	// Read from environment variables and return an error if not found
	config.Host, err = getEnv("DB_HOST")
	if err != nil {
		return config, err
	}
	config.Port, err = getEnv("DB_PORT")
	if err != nil {
		return config, err
	}
	config.User, err = getEnv("DB_USER")
	if err != nil {
		return config, err
	}
	config.Password, err = getEnv("DB_PASSWORD")
	if err != nil {
		return config, err
	}
	config.DBName, err = getEnv("DB_NAME")
	if err != nil {
		return config, err
	}
	config.SSLMode, err = getEnv("DB_SSLMODE")
	if err != nil {
		return config, err
	}

	return config, nil
}

type Config struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	DBName   string `json:"dbname"`
	SSLMode  string `json:"sslmode"`
}

func Configure() {
	SetupDatabaseConnections()
	ApplyAllStaticSchemaMigrations()
	FetchAllSchemas()
	FetchAllOauthConfigs()
	ApplyAllDynamicSchemaMigrations()
	SetupSMTPConfig()
}

func FetchAllSchemas() {
	var schemas []ct.Schema

	err := DB.Model(&ct.Schema{}).Preload("Fields").Find(&schemas).Error
	if err == nil {
		log.Println("Fetched All Schemas")
	}

	for _, s := range schemas {
		utils.Schemas[s.Table] = s
	}
}

func ApplyAllDynamicSchemaMigrations() {
	for _, s := range utils.Schemas {
		err := utils.CreateSchema(DB, s)
		if err != nil {
			log.Fatalf("Could not apply migration for '%v'", s.Table)
		}
	}
}

func ApplyAllStaticSchemaMigrations() {
	err := DB.AutoMigrate(&ct.Schema{})
	if err != nil {
		log.Fatalf("failed to create table: %v", err)
	}
	err = DB.AutoMigrate(&ct.SchemaField{})
	if err != nil {
		log.Fatalf("failed to create table: %v", err)
	}
	err = DB.AutoMigrate(&ct.Auth{})
	if err != nil {
		log.Fatalf("failed to create table: %v", err)
	}
	err = DB.AutoMigrate(&ct.Session{})
	if err != nil {
		log.Fatalf("failed to create table: %v", err)
	}
	err = DB.AutoMigrate(&ct.OTP{})
	if err != nil {
		log.Fatalf("failed to create table: %v", err)
	}
	err = DB.AutoMigrate(&ct.OauthConfig{})
	if err != nil {
		log.Fatalf("failed to create table: %v", err)
	}
	err = DB.AutoMigrate(&ct.SMTPConfig{})
	if err != nil {
		log.Fatalf("failed to create table: %v", err)
	}
}

func SetupDatabaseConnections() {
	config, err := loadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v\n", err)
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", config.Host, config.Port, config.User, config.Password, config.DBName, config.SSLMode)

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect database")
	} else {
		log.Println("Connected to database")
	}
}

func UpdateAndMigrateSchemas() {
	FetchAllSchemas()
	ApplyAllDynamicSchemaMigrations()
}

func SetupSMTPConfig() {
	err := DB.First(&SMTPConfig).Error
	if err != nil {
		log.Printf("Failed to get SMTP config from database: %v\n", err)
	}
	log.Println("SMTP config loaded")
}

func getEnv(key string) (string, error) {
	value := os.Getenv(key)
	if value == "" {
		return "", errors.New("required environment variable not set: " + key)
	}
	return value, nil
}
