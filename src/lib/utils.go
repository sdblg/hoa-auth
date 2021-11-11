package lib

import (
	"flag"
	"github.com/sdblg/hoa-auth/src/constants"
	"os"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

func InitEnv() {
	env := flag.String("env", "local", "-env=local")
	flag.Parse()
	environment := *env
	if _, ok := os.LookupEnv(constants.CurrentEnv); !ok {
		_ = os.Setenv(constants.CurrentEnv, environment)
	}

	var configFile = "config/" + environment + "/application.env"
	if err := godotenv.Load(configFile); err != nil {
		log.Panicf("Loading Config file: %v %v", configFile, err)
	}

	log.Infof("Current environment :: %v", environment)
}
