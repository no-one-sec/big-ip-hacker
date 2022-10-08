package big_ip

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/stretchr/testify/assert"
	"net"
	"testing"
)

func TestDecode(t *testing.T) {
	decode, err := Decode("BIGipServer~Banner~pool-pbanapp.banner-8038=1312824236.26143.0000; path=/; Httponly; Secure")
	if err != nil {
		_, _ = fmt.Fprint(color.Output, err.Error())
		return
	}
	assert.Nil(t, err)
	assert.NotNil(t, decode)
	assert.Equal(t, decode.PoolName, "BIGipServer~Banner~pool-pbanapp.banner-8038")
	assert.Equal(t, decode.IP, net.ParseIP("78.64.27.172"))
	assert.Equal(t, decode.Port, uint16(26143))
	assert.Equal(t, decode.End, "0000")
}
