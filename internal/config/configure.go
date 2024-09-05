package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/cobaltbase/cobaltbase/internal/ct"
	"github.com/cobaltbase/cobaltbase/internal/utils"
	"github.com/go-playground/validator/v10"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

var Validate = validator.New()

const configFileName = "dbconfig.json"

func loadConfig() (Config, error) {
	var config Config

	// Get the current working directory
	dir, err := os.Getwd()
	if err != nil {
		return config, fmt.Errorf("failed to get current directory: %v", err)
	}

	// Construct the full path to the config file
	configPath := filepath.Join(dir, configFileName)

	// Read the config file
	file, err := os.ReadFile(configPath)
	if err != nil {
		return config, fmt.Errorf("failed to read config file: %v", err)
	}

	// Unmarshal the JSON data into our Config struct
	err = json.Unmarshal(file, &config)
	if err != nil {
		return config, fmt.Errorf("failed to parse config file: %v", err)
	}

	return config, nil
}

func saveConfig(config Config) error {
	// Get the current working directory
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %v", err)
	}

	// Construct the full path to the config file
	configPath := filepath.Join(dir, configFileName)

	// Marshal the config struct to JSON
	data, err := json.MarshalIndent(config, "", "    ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %v", err)
	}

	// Write the JSON data to the file
	err = os.WriteFile(configPath, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write config file: %v", err)
	}

	return nil
}

type Config struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	DBName   string `json:"dbname"`
}

func Configure() {
	SetupDatabaseConnections()
	ApplyAllStatiSchemaMigrations()
	FetchAllSchemas()
	ApplyAllDynamicSchemaMigrations()
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

func ApplyAllStatiSchemaMigrations() {
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
}

func SetupDatabaseConnections() {
	config, err := loadConfig()
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		config = Config{
			Host:     "localhost",
			Port:     "5432",
			User:     "cb",
			Password: "cobaltbase",
			DBName:   "cobaltbasedb",
		}
		if err := saveConfig(config); err != nil {
			fmt.Printf("Failed to save default config: %v\n", err)
		}
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", config.Host, config.Port, config.User, config.Password, config.DBName)

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
