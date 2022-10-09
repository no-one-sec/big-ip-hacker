package big_ip

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/CC11001100/go-StringBuilder/pkg/string_builder"
	"github.com/fatih/color"
	"github.com/go-resty/resty/v2"
)

// FromUrl 从url中解析
func FromUrl(targetUrl string) {

	fmt.Printf("targetUrl = %s\n", targetUrl)

	response, err := request(targetUrl)
	if err != nil {
		_, _ = fmt.Fprint(color.Output, err.Error())
		return
	}
	bigIpCookies, err := parseBIGipCookieFromHttpResponse(response)
	if err != nil {
		_, _ = fmt.Fprint(color.Output, err.Error())
		return
	}

	for _, bigIpCookie := range bigIpCookies {
		FromCookie(bigIpCookie)
	}

}

// FromCookie 从Cookie中解析
func FromCookie(cookieString string) {

	fmt.Printf("\nbegin decode cookie...\n")

	decode, err := Decode(cookieString)
	if err != nil {
		_, _ = fmt.Fprint(color.Output, err.Error())
		return
	}

	color.HiGreen("decode success: \n")
	color.HiGreen("%15s : %s\n", "pool name", decode.PoolName)
	color.HiGreen("%15s : %s\n", "ip", decode.IP.String())
	color.HiGreen("%15s : %d\n", "port", decode.Port)
	color.HiGreen("%15s : %s\n", "end", decode.End)
}

func request(targetUrl string) (*resty.Response, error) {
	for tryTimes := 1; tryTimes <= 3; tryTimes++ {

		fmt.Printf("tryTimes = %d, begin request...\n", tryTimes)

		resp, err := resty.New().
			SetTransport(&http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			}).
			SetTimeout(time.Second * 30).
			R().
			Get(targetUrl)
		if err != nil {
			color.HiRed("request failed: %s", err.Error())
			continue
		}
		color.HiGreen("request success!")
		return resp, nil
	}
	return nil, errors.New(color.HiRedString("run out of retries, still unsuccessful\n"))
}

// 从响应中解析Cookie，如果没有解析到的话会返回错误
func parseBIGipCookieFromHttpResponse(response *resty.Response) ([]string, error) {
	fmt.Printf("\ntry to find big ip cookie from response headers...\n")
	var bigIpCookie []string
	index := 1
	for name, valueSlice := range response.Header() {
		for _, value := range valueSlice {
			if strings.HasPrefix(value, "BIGipServer~") {
				bigIpCookie = append(bigIpCookie, value)
				color.HiGreen("%3d: %s : %s\n", index, name, value)
			} else {
				fmt.Printf("%3d: %s : %s\n", index, name, value)
			}
			index++
		}
	}
	if len(bigIpCookie) == 0 {
		return bigIpCookie, errors.New(color.HiRedString("The cookie for BIG IP could not be found in the response header\n"))
	}
	return bigIpCookie, nil
}

type BigIpCookie struct {
	PoolName string
	IP       net.IP
	Port     uint16
	End      string
}

// Decode 将传入的Big IP Cookie解析为结构化的struct
// cookieString: BIGipServer~ag_web~webhosting80_pool=975239178.20480.0000; path=/; Httponly; Secure
func Decode(cookieString string) (*BigIpCookie, error) {
	split := strings.SplitN(cookieString, "=", 2)
	if len(split) != 2 {
		return nil, buildCookieFormatError(cookieString)
	}

	bigIpCookie := &BigIpCookie{}

	bigIpCookie.PoolName = split[0]
	cookieValue := strings.SplitN(split[1], ";", 2)[0]
	split = strings.Split(cookieValue, ".")
	if len(split) != 3 {
		return nil, buildCookieFormatError(cookieString)
	}

	hostUInt32 := split[0]
	port := split[1]
	bigIpCookie.End = split[2]

	// 把主机部分转换为ip类型
	var err error
	if bigIpCookie.IP, err = convertUInt32StringToIp(hostUInt32); err != nil {
		return nil, err
	}

	if bigIpCookie.Port, err = parsePort(port); err != nil {
		return nil, err
	}

	return bigIpCookie, nil
}

// 解析字符串类型的端口，会进行范围有效性校验
func parsePort(portString string) (uint16, error) {
	atoi, err := strconv.Atoi(portString)
	if err != nil {
		return 0, errors.New(color.HiRedString("parse port error: %s", err.Error()))
	}
	if atoi <= 0 || atoi > 65535 {
		return 0, errors.New(color.HiRedString("port range error: %s, must in [1,65535]", portString))
	}
	return uint16(atoi), nil
}

// uint32字符串格式的ip转为net.IP格式
func convertUInt32StringToIp(hostUInt32String string) (net.IP, error) {
	parseUint32, err := strconv.ParseUint(hostUInt32String, 10, 32)
	if err != nil {
		return nil, errors.New(color.HiRedString("parse ip error: %s", err.Error()))
	}

	return uint32ToIp(uint32(parseUint32))
}

// uint32格式的ip转为net.IP格式
func uint32ToIp(ipUint32 uint32) (net.IP, error) {
	ipString := fmt.Sprintf("%d.%d.%d.%d", (ipUint32>>24)&0xFF, (ipUint32>>16)&0xFF, (ipUint32>>8)&0xFF, ipUint32&0xFF)
	ip := net.ParseIP(ipString)
	if ip == nil {
		return nil, errors.New(color.HiRedString("parse ip error: %d --> %s", ipUint32, ipString))
	}
	return ip, nil
}

// Cookie格式错误的时候构造友好一些的提示信息
func buildCookieFormatError(cookieString string) error {
	errMsg := string_builder.New().AppendString(color.HiRedString("Cookie parsing failed, perhaps due to formatting error?\n")).
		AppendString("\n").
		AppendString("Expected format: \n").
		AppendString(color.HiGreenString("  BIGipServer~${pool_name}=${value}; path=/; Httponly; Secure\n")).
		AppendString("Example: \n").
		AppendString(color.HiGreenString("  BIGipServer~Banner~pool-pbanapp.banner-8038=1312824236.26143.0000; path=/; Httponly; Secure\n")).
		AppendString(color.HiGreenString("  BIGipServer~Banner~pool-pbanapp.banner-8038=1312824236.26143.0000\n")).
		AppendString("\n").
		AppendString("Your Input: \n").
		AppendString("  ").AppendString(color.HiRedString(cookieString)).
		AppendString("\n").
		String()
	return errors.New(errMsg)
}
