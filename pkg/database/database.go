package database

import (
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// PostgreSQLConfig used to set base config for all database types.
type PostgreSQLConfig struct {
	Host     string `json:"host" mapstructure:"host" yaml:"host"`
	Database string `json:"database" mapstructure:"database" yaml:"database"`
	Port     int    `json:"port" mapstructure:"port" yaml:"port"`
	Username string `json:"username" mapstructure:"username" yaml:"username"`
	Password string `json:"password" mapstructure:"password" yaml:"password"`
	Options  string `json:"options" mapstructure:"options" yaml:"options"`

	// ConnMaxIdleTime reformat type time.Duration (minute)
	ConnMaxIdleTime int `json:"conn_max_idle_time" yaml:"conn_max_idle_time"`
	MaxIdleConn     int `json:"max_idle_conn" yaml:"max_idle_conn"`
	MaxOpenConn     int `json:"max_open_conn" yaml:"max_open_conn"`
}

func (c PostgreSQLConfig) DSN() string {
	return fmt.Sprintf("host=%v port=%v user=%v sslmode=disable dbname=%v password=%v ", c.Host, c.Port, c.Username, c.Database, c.Password)
}

func (c PostgreSQLConfig) String() string {
	return fmt.Sprintf("postgresql://%s", c.DSN())
}

// PostgresSQLDefaultConfig returns default config for mysql, usually use on development.
func PostgresSQLDefaultConfig() PostgreSQLConfig {
	return PostgreSQLConfig{
		Host:            "127.0.0.1",
		Database:        "test",
		Port:            5432,
		Username:        "default",
		Password:        "secret",
		Options:         "",
		ConnMaxIdleTime: 10,
		MaxIdleConn:     16,
		MaxOpenConn:     64,
	}
}

func (c PostgreSQLConfig) ConnectDatabase() *gorm.DB {
	db, err := gorm.Open(postgres.Open(c.DSN()), &gorm.Config{
		PrepareStmt: true,
	})
	if err != nil {
		panic(err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}

	sqlDB.SetConnMaxIdleTime(5 * time.Minute)
	sqlDB.SetMaxIdleConns(16)
	sqlDB.SetMaxOpenConns(64)

	if err = sqlDB.Ping(); err != nil {
		panic(err)
	}
	fmt.Printf("Connected with database server host [%v], port [%v], db [%v] \n", c.Host, c.Port, c.Database)
	return db
}
