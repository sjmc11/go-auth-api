package env

import (
	"auth-api/helpers"
	"bufio"
	"fmt"
	"github.com/logrusorgru/aurora/v4"
	"log"
	"os"
	"strings"
)

var keys = make(map[string]string)

type Environment struct {
	EnvPath string
}

func (e *Environment) LoadEnv() {

	if e.EnvPath == "" {
		e.EnvPath = ".env"
	}
	_, err := os.Stat(e.EnvPath)
	if err != nil {
		_, err = os.Stat(helpers.BasePath() + "/" + e.EnvPath)
		if err == nil {
			e.EnvPath = helpers.BasePath() + "/" + e.EnvPath
		} else {
			log.Fatal(err)
		}
	}

	file, err := os.Open(e.EnvPath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		if line != "" && string(line[0]) != "#" {
			envVar := strings.Split(line, "=")
			if os.Getenv(strings.TrimSpace(envVar[0])) != "" {
				keys[strings.TrimSpace(envVar[0])] = os.Getenv(strings.TrimSpace(envVar[0]))
			} else {
				keys[strings.TrimSpace(envVar[0])] = strings.TrimSpace(envVar[1])
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Println(aurora.Green("✓ ENV LOADED"))
	fmt.Println(aurora.Gray(12, "----------"))
	for key, value := range keys {
		fmt.Println(aurora.Gray(12, key+"  :  "+value))
	}
	fmt.Println(aurora.Gray(12, "----------"))

}

func Get(key string) string {
	return keys[key]
}
