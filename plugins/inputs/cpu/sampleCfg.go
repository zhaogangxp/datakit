package cpu

const sampleCfg = `
[[inputs.cpu]]
  ## Collect interval, default is 10 seconds. (optional)
  interval = '10s'
  ##
  ## Collect CPU usage per core, default is false. (optional)
  percpu = false
  ##
  ## Setting disable_temperature_collect to false will collect cpu temperature stats for linux.
  ##
  # disable_temperature_collect = false
  enable_temperature = true

  [inputs.cpu.tags]
    # some_tag = "some_value"
    # more_tag = "some_other_value"
`
