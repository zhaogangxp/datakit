package sqlserver

import (
	"fmt"
	"testing"

	_ "github.com/denisenkom/go-mssqldb"
)

func TestCon(t *testing.T) {
	n := Input{
		Host:     "10.100.64.109:1433",
		User:     "_",
		Password: "_",
	}
	if err := n.initDB(); err != nil {
		l.Error(err.Error())
		return
	}

	n.getMetric()
	for _, v := range collectCache {
		fmt.Println(v.String())
	}
}
