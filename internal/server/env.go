package server

import (
	"os"
	"strconv"

	"github.com/sirupsen/logrus"
)

type Env struct {
	HOST                       string
	POSTGRES_DB                string
	POSTGRES_USER              string
	POSTGRES_PASSWORD          string
	PGADMIN_DEFAULT_EMAIL      string
	PGADMIN_DEFAULT_PASSWORD   string
	PGADMIN_CONFIG_SERVER_MODE string
	JWT_SECRET                 string
	JWT_DURATION_MIN           int
	REFRESH_DURATION_DAYS      int
}

func NewEnv() *Env {
	JWT_DURATION_MIN, err := strconv.Atoi(os.Getenv("JWT_DURATION_MIN"))
	if err != nil {
		logrus.Fatal("Could not get JWT_DURATION_MIN env variable (should be a number of minutes token expiration time)")
	}

	REFRESH_DURATION_DAYS, err := strconv.Atoi(os.Getenv("REFRESH_DURATION_DAYS"))
	if err != nil {
		logrus.Fatal("Could not get REFRESH_DURATION_DAYS env variable (should be a number of days refresh token expiration time)")
	}

	return &Env{
		HOST:                       os.Getenv("HOST"),
		POSTGRES_DB:                os.Getenv("POSTGRES_DB"),
		POSTGRES_USER:              os.Getenv("POSTGRES_USER"),
		POSTGRES_PASSWORD:          os.Getenv("POSTGRES_PASSWORD"),
		PGADMIN_DEFAULT_EMAIL:      os.Getenv("PGADMIN_DEFAULT_EMAIL"),
		PGADMIN_DEFAULT_PASSWORD:   os.Getenv("PGADMIN_DEFAULT_PASSWORD"),
		PGADMIN_CONFIG_SERVER_MODE: os.Getenv("PGADMIN_CONFIG_SERVER_MODE"),
		JWT_SECRET:                 os.Getenv("JWT_SECRET"),
		JWT_DURATION_MIN:           JWT_DURATION_MIN,
		REFRESH_DURATION_DAYS:      REFRESH_DURATION_DAYS,
	}
}
