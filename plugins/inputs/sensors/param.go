// Package sensors collect hardware sensor metrics.
package sensors

import (
	"gitlab.jiagouyun.com/cloudcare-tools/cliutils/logger"
)

var (
	inputName = "sensors"

	sampleConfig = `
[[inputs.sensors]]
  ## Command path of 'senssor' usually under /usr/bin/sensors
  # path = "/usr/bin/senssors"

  ## Gathering interval
  # interval = "10s"

  ## Command timeout
  # timeout = "3s"

  ## Customer tags, if set will be seen with every metric.
  [inputs.sensors.tags]
    # "key1" = "value1"
    # "key2" = "value2"
`
	l = logger.DefaultSLogger(inputName)
)
