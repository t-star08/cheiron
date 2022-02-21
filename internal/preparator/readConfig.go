package preparator

import (
	"github.com/t-star08/cheiron/internal/config"
	"github.com/t-star08/hand/pkg/io/ioJson"
)

func readConf() (string, *config.Cheiron, error) {
	var (
		conf *config.Cheiron
		pathToConfDir string
	)

	if path, err := config.SearchConfigDirPath(); err != nil {
		return pathToConfDir, conf, err
	} else {
		pathToConfDir = path
	}

	var pathToConfFile = pathToConfDir + "/" + config.CONF_FILE_NAME
	if err := ioJson.Gets(pathToConfFile, &conf); err != nil {
		return pathToConfDir, conf, err
	}

	return pathToConfDir, conf, nil
}
