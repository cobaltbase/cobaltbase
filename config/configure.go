package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	ct "github.com/cobaltbase/cobaltbase/customTypes"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

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

	err = DB.AutoMigrate(&ct.Schema{})
	if err != nil {
		fmt.Printf("failed to create table: %v", err)
	}
	err = DB.AutoMigrate(&ct.SchemaField{})
	if err != nil {
		fmt.Printf("failed to create table: %v", err)
	}
}
