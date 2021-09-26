package send

import (
	"errors"
	"math/rand"
	"net"
	"time"

	"github.com/haiyiyun/cache"
	"github.com/haiyiyun/log"
)

var ipsCache *cache.Cache = cache.New(12*time.Hour, 3*time.Hour)

func IPs(mxDomain string) (map[string][]string, error) {
	ips := map[string][]string{}
	mxs, err := net.LookupMX(mxDomain)
	if err != nil {
		log.Error("<Addr>", " LookupMX Error:", err)
		return nil, err
	}

	log.Debug("<IPs> mxs:", mxs)
	if len(mxs) > 0 {
		for _, mx := range mxs {
			host := mx.Host
			if is, err := net.LookupIP(host); err == nil {
				log.Debug("<IPs> is:", is)
				for _, ip := range is {
					//暂不支持IPv6,只提取IPv4的地址
					ipv4 := ip.To4()
					if len(ipv4) == net.IPv4len {
						ips[host] = append(ips[host], ipv4.String())
					}
				}
			} else {
				log.Error("<IPs>", "LookupIP IP:", host, " Error:", err)
				return nil, err
			}
		}

	}

	log.Debugf("<IPs> IPs:%+v", ips)
	return ips, nil
}

func RandomIP(mxDomain string) (string, error) {
	var ips map[string][]string
	if x, found := ipsCache.Get(mxDomain); found {
		ips = x.(map[string][]string)
	} else {
		ipst, err := IPs(mxDomain)
		if err != nil || ipst == nil {
			return "", err
		}

		if len(ipst) == 0 {
			return "", errors.New("Empty Ip")
		}

		ipsCache.Add(mxDomain, ipst, 0)
		ips = ipst
	}

	selectOneIps := []string{}
	for _, is := range ips {
		if len(is) > 0 {
			for _, i := range is {
				if !findIp(i, selectOneIps) {
					selectOneIps = append(selectOneIps, i)
				}
			}
		}
	}

	log.Debug("<RandomIP> selectOneIps:", selectOneIps)
	if len(selectOneIps) == 0 {
		log.Debug("<RandomIP> Select Ip Error")
		return "", errors.New("Select Ip Error")
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	ip := selectOneIps[r.Intn(len(selectOneIps))]
	log.Debug("<RandomIP> selectSecondIps:", ip)
	return ip, nil
}

func findIp(ip string, ips []string) bool {
	for _, i := range ips {
		if ip == i {
			return true
		}
	}

	return false
}
