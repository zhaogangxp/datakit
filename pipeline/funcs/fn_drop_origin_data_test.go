package funcs

import "testing"

func TestDropOriginData(t *testing.T) {
	cases := []struct {
		name, pl, in string
		key          string
		fail         bool
	}{
		{
			name: "value type: string",
			in:   `162.62.81.1 - - [29/Nov/2021:07:30:50 +0000] "POST /?signature=b8d8ea&timestamp=1638171049 HTTP/1.1" 200 413 "-" "Mozilla/4.0"`,
			pl: `
	grok(_, "%{IPORHOST:client_ip} %{NOTSPACE} %{NOTSPACE} \\[%{HTTPDATE:time}\\] \"%{DATA} %{GREEDYDATA} HTTP/%{NUMBER}\" %{INT:status_code} %{INT:bytes}")
	drop_origin_data()
	`,
			key:  "message",
			fail: false,
		},
	}

	for idx, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			runner, err := NewTestingRunner(tc.pl)
			if err != nil {
				if tc.fail {
					t.Logf("[%d]expect error: %s", idx, err)
				} else {
					t.Errorf("[%d] failed: %s", idx, err)
				}
				return
			}

			if err := runner.Run(tc.in); err != nil {
				t.Error(err)
				return
			}
			t.Log(runner.Result())
			if v, err := runner.GetContentStr(tc.key); err == nil || v != "" {
				t.Errorf("[%d] failed: key `%s` value `%v`", idx, tc.key, v)
			} else {
				t.Logf("[%d] PASS", idx)
			}
		})
	}
}
