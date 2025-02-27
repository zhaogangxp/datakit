package kubernetes

const sampleCfg = `
[inputs.kubernetes]
  ## URL for the Kubernetes API
  ## daemonset at: "https://kubernetes.default:443"
  url = "https://kubernetes.default:443"

  ## metrics interval
  interval = "60s"

  ## Authorization level:
  ##   bearer_token -> bearer_token_string -> TLS
  ## Use bearer token for authorization. ('bearer_token' takes priority)
  ## linux at:   /run/secrets/kubernetes.io/serviceaccount/token
  ## windows at: C:\var\run\secrets\kubernetes.io\serviceaccount\token
  bearer_token = "/run/secrets/kubernetes.io/serviceaccount/token"
  # bearer_token_string = "<your-token-string>"

  ## TLS Config
  # tls_ca = "/path/to/ca.pem"
  # tls_cert = "/path/to/cert.pem"
  # tls_key = "/path/to/key.pem"
  ## Use TLS but skip chain & host verification
  # insecure_skip_verify = false

  [inputs.kubernetes.tags]
  # some_tag = "some_value"
  # more_tag = "some_other_value"
`

const (
	defaultStringValue    string = ""
	defaultBoolerValue    bool   = false
	defaultIntegerValue   int    = 0
	defaultInteger32Value int32  = 0
	defaultInteger64Value int64  = 0
)
