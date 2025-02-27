package nsq

const sampleCfg = `
[[inputs.nsq]]
  ## NSQ Lookupd HTTP API endpoint
  lookupd = "http://localhost:4161"

  ## NSQD HTTP API endpoint
  ## example:
  ##   ["http://localhost:4151"]
  nsqd = []
  
  ## time units are "ms", "s", "m", "h"
  interval = "10s"
  
  ## Optional TLS Config
  # tls_ca = "/etc/telegraf/ca.pem"
  # tls_cert = "/etc/telegraf/cert.pem"
  # tls_key = "/etc/telegraf/key.pem"
  ## Use TLS but skip chain & host verification
  # insecure_skip_verify = false
  
  [inputs.nsq.tags]
  # some_tag = "some_value"
  # more_tag = "some_other_value"
`
