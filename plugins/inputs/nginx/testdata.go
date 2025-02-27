package nginx

var vtsModelHandleData = `
{
    "hostName": "tan-thinkpad-e450",
    "moduleVersion": "0.1.19.dev.91bdb14",
    "nginxVersion": "1.9.2",
    "loadMsec": 1618888188619,
    "nowMsec": 1618888193244,
    "connections": {
        "active": 1,
        "reading": 0,
        "writing": 1,
        "waiting": 0,
        "accepted": 1,
        "handled": 1,
        "requests": 1
    },
    "sharedZones": {
        "name": "ngx_http_vhost_traffic_status",
        "maxSize": 1048575,
        "usedSize": 0,
        "usedNode": 0
    },
    "serverZones": {
        "*": {
            "requestCounter": 0,
            "inBytes": 0,
            "outBytes": 0,
            "responses": {
                "1xx": 0,
                "2xx": 0,
                "3xx": 0,
                "4xx": 0,
                "5xx": 0,
                "miss": 0,
                "bypass": 0,
                "expired": 0,
                "stale": 0,
                "updating": 0,
                "revalidated": 0,
                "hit": 0,
                "scarce": 0
            },
            "requestMsecCounter": 0,
            "requestMsec": 0,
            "requestMsecs": {
                "times": [
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0
                ],
                "msecs": [
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0
                ]
            },
            "requestBuckets": {
                "msecs": [],
                "counters": []
            },
            "overCounts": {
                "maxIntegerSize": 18446744073709551615,
                "requestCounter": 0,
                "inBytes": 0,
                "outBytes": 0,
                "1xx": 0,
                "2xx": 0,
                "3xx": 0,
                "4xx": 0,
                "5xx": 0,
                "miss": 0,
                "bypass": 0,
                "expired": 0,
                "stale": 0,
                "updating": 0,
                "revalidated": 0,
                "hit": 0,
                "scarce": 0,
                "requestMsecCounter": 0
            }
        }
    },
    "upstreamZones": {
        "test": [
            {
                "server": "10.100.64.215:8888",
                "requestCounter": 0,
                "inBytes": 0,
                "outBytes": 0,
                "responses": {
                    "1xx": 0,
                    "2xx": 0,
                    "3xx": 0,
                    "4xx": 0,
                    "5xx": 0
                },
                "requestMsecCounter": 0,
                "requestMsec": 0,
                "requestMsecs": {
                    "times": [],
                    "msecs": []
                },
                "requestBuckets": {
                    "msecs": [],
                    "counters": []
                },
                "responseMsecCounter": 0,
                "responseMsec": 0,
                "responseMsecs": {
                    "times": [],
                    "msecs": []
                },
                "responseBuckets": {
                    "msecs": [],
                    "counters": []
                },
                "weight": 1,
                "maxFails": 1,
                "failTimeout": 10,
                "backup": false,
                "down": false,
                "overCounts": {
                    "maxIntegerSize": 18446744073709551615,
                    "requestCounter": 0,
                    "inBytes": 0,
                    "outBytes": 0,
                    "1xx": 0,
                    "2xx": 0,
                    "3xx": 0,
                    "4xx": 0,
                    "5xx": 0,
                    "requestMsecCounter": 0,
                    "responseMsecCounter": 0
                }
            }
        ]
    }
}`
