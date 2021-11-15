package lib

import (
	"flag"
	"github.com/sdblg/hoa-auth/src/constants"
	"os"
	"regexp"

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

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func IsEmailValid(e string) bool {
	if len(e) < 3 && len(e) > 254 {
		return false
	}
	return emailRegex.MatchString(e)
}
