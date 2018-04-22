package httputil

import (
	"net"
	"net/http"
	"strings"
)

// GetClientIPAddress return IP address
func GetClientIPAddress(r *http.Request) string {
	ipAddress := r.Header.Get("X-Forwarded-For")
	ipAddress = strings.Replace(ipAddress, " ", "", -1)
	//get from r.RemoteAddr
	if ipAddress == "" {
		ipFirst := strings.Replace(r.RemoteAddr, " ", "", -1)

		ipAddressArr := strings.Split(ipFirst, ",")
		if len(ipAddressArr) > 0 {
			ipAddress = ipAddressArr[0]
		}

		ipAddressArr = strings.Split(ipAddress, ":")
		if len(ipAddressArr) > 0 {
			ipAddress = ipAddressArr[0]
		}
	} else {
		//if no empty, looking for sign , and :
		ipAddressArr := strings.Split(ipAddress, ",")
		if len(ipAddressArr) > 0 {
			ipAddress = ipAddressArr[0]
		}

		ipAddressArr = strings.Split(ipAddress, ":")
		if len(ipAddressArr) > 0 {
			ipAddress = ipAddressArr[0]
		}
	}

	// to make sure that this IP is in right format, so we can input in database.
	addr := net.ParseIP(ipAddress)
	if addr == nil {
		ipAddress = ""
	}
	return ipAddress
}
