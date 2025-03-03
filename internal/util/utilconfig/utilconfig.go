// Package utilconfig ...
package utilconfig

import (
	"encoding/json"
	"fmt"
	"go-infra/internal/util/utilhttp"
	xlog "go-infra/internal/util/utillog"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func LoadConfig(cfgPtr any, dir string, fileName string) error {

	xlog.Info("loading config from: %v", dir)

	isHTTP := strings.HasPrefix(dir, "http")

	// TODO from s3://

	if isHTTP {

		err := fromURL(cfgPtr, dir, fileName)
		if err != nil {
			return err
		}

	} else {
		err := fromFile(cfgPtr, dir, fileName)
		if err != nil {
			return err
		}
	}

	return nil
}

// fromFile errIfNotExists argument soft binding, no error if file not exists
func fromFile(cfgPtr any, dir string, file string) error {

	if file == "" {
		return nil
	}

	if !strings.HasSuffix(file, ".json") {
		return fmt.Errorf("error file not match  *.json: %v", file)
	}

	fullPath, err := filepath.Abs(filepath.Join(dir, file))

	if err != nil {
		return err
	}

	fullPath = filepath.Clean(fullPath)

	data, err := os.ReadFile(fullPath)

	if err != nil {
		return fmt.Errorf("error with file %v: %v", fullPath, err)
	}

	xlog.Info("loading config from file: %v", fullPath)

	err = fromJSON(cfgPtr, string(data))

	if err != nil {
		return err
	}

	return nil
}

// FromURL errIfNotExists argument soft binding, no error if file not exists
func fromURL(cfgPtr any, dir string, file string) error {

	if file == "" {
		return nil
	}

	if !strings.HasSuffix(file, ".json") {
		return fmt.Errorf("error file not match  *.json: %v", file)
	}

	fullPath := dir + "/" + file

	_, err := url.Parse(fullPath)
	if err != nil {
		return fmt.Errorf("invalid URL: %v", err)
	}

	// fmt.Println("reading config from file: ", file)

	data, err := utilhttp.GetBytes(fullPath, nil, nil)

	if err != nil {
		return fmt.Errorf("error with file %v: %v", fullPath, err)
	}

	xlog.Info("loading config from file: %v", fullPath)

	err = fromJSON(cfgPtr, string(data))
	if err != nil {
		return err
	}

	return nil
}

func expandEnv(data string) string {

	re := regexp.MustCompile(`\$\{([A-Z_][0-9_A-Z]*)\}`)
	data = re.ReplaceAllStringFunc(data, func(match string) string {
		name := match[2 : len(match)-1] // Remove '${' and '}'
		val := os.Getenv(name)
		if val == "" {
			xlog.Warn("missing env value for: %v", match)
		}
		return val // Return the original match if not found
	})

	return data

	// data = os.Expand(data, func(s string) string {

	// 	// TODO chek if var name [A-Z_0-9]+

	// 	parts := strings.SplitN(s, ":", 2)
	// 	name := parts[0]
	// 	// tail:=parts[1]
	// 	val := os.Getenv(name)

	// 	if val == "" {
	// 		//
	// 		xlog.Warn("missing env value for: %v", s)
	// 	}

	// 	return val

	// })

}

func fromJSON(cfgPtr any, data string) error {

	if data == "" {
		return nil
	}

	data = expandEnv(data)

	err := json.Unmarshal([]byte(data), cfgPtr)

	if err != nil {
		return err
	}

	return nil
}
