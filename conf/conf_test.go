package conf

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

const CONF = "./conf.json"

func TestReadConfig(t *testing.T) {
	conf, err := ReadConfig(CONF)
	assert.NoError(t, err)
	assert.NotNil(t, conf)
	assert.Equal(t, conf.DataSourceName, "root:linnaname@123456@tcp(localhost:3306)/sequence_test")
}
