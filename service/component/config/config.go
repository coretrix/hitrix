package config

import (
	"os"
	"regexp"
	"strings"

	"github.com/gookit/config"
	"github.com/gookit/config/yaml"
	"github.com/joho/godotenv"
)

type IConfig interface {
	Get(key string, findByPath ...bool) (value interface{}, ok bool)
	String(key string) (value string, ok bool)
	DefString(key string, defVal ...string) string
	MustString(key string) string

	Int(key string) (value int, ok bool)
	DefInt(key string, defVal ...int) int
	MustInt(key string) int
	Int64(key string) (value int64, ok bool)
	DefInt64(key string, defVal ...int64) int64
	MustInt64(key string) int64
	// Bool looks up a value for a key in this section and attempts to parse that value as a boolean,
	// along with a boolean result similar to a map lookup.
	// of following(case insensitive):
	//  - true
	//  - yes
	//  - false
	//  - no
	//  - 1
	//  - 0
	// The `ok` boolean will be false in the event that the value could not be parsed as a bool
	Bool(key string) (value bool, ok bool)
	DefBool(key string, defVal ...bool) bool
	MustBool(key string) bool

	Float(key string) (value float64, ok bool)
	DefFloat(key string, defVal ...float64) float64

	Ints(key string) (arr []int, ok bool)
	IntMap(key string) (mp map[string]int, ok bool)

	Strings(key string) (arr []string, ok bool)
	StringMap(key string) (mp map[string]string, ok bool)

	MapStruct(key string, v interface{}) (err error)
	MapStructure(key string, v interface{}) (err error)

	// Structure get config data and map to a structure.
	// usage:
	// 	dbInfo := Db{}
	// 	config.Structure("db", &dbInfo)
	Structure(key string, v interface{}) (err error)

	// IsEmpty of the config
	IsEmpty() bool
	GetFolderPath() string
}

type Config struct {
	*config.Config
	folderPath string
}

func NewConfig(appName, mode, localConfigFolder string) (*Config, error) {
	var configFolder string

	appFolder, hasConfigFolder := os.LookupEnv("APP_FOLDER")
	if !hasConfigFolder {
		configFolder = localConfigFolder
	} else {
		configFolder = appFolder + "/config"
	}

	err := loadEnvConfig(mode, configFolder)
	if err != nil {
		return nil, err
	}

	c := config.NewEmpty(appName)
	c.WithOptions(config.ParseEnv)
	c.AddDriver(yaml.Driver)

	yamlFileAppConfig, err := os.ReadFile(configFolder + "/" + appName + "/config.yaml")
	if err != nil {
		return nil, err
	}

	yamlFileHitrixConfig, err := os.ReadFile(configFolder + "/hitrix.yaml")
	if err != nil {
		return nil, err
	}

	err = c.LoadSources(config.Yaml, parseEnvVariables(yamlFileAppConfig), parseEnvVariables(yamlFileHitrixConfig))
	if err != nil {
		return nil, err
	}

	configService := &Config{
		c,
		configFolder,
	}

	return configService, nil
}

func parseEnvVariables(content []byte) []byte {
	var newContent string
	newContent = string(content)

	re := regexp.MustCompile(`ENV\[(.*?)\]`)

	subMatchAll := re.FindAllString(string(content), -1)
	for _, element := range subMatchAll {
		element = strings.TrimPrefix(element, "ENV[")
		element = strings.TrimSuffix(element, "]")

		values := strings.Split(os.Getenv(element), ";")
		if len(values) > 1 {
			newContent = strings.Replace(newContent, "ENV["+element+"]", "["+strings.Join(values, ",")+"]", -1)
		} else {
			newContent = strings.Replace(newContent, "ENV["+element+"]", os.Getenv(element), -1)
		}
	}

	return []byte(newContent)
}

func loadEnvConfig(mode, configFolder string) error {
	if _, err := os.Stat(configFolder + "/.env." + mode); !os.IsNotExist(err) {
		err := godotenv.Load(configFolder + "/.env." + mode)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Config) GetFolderPath() string {
	return c.folderPath
}
