package main

import "strings"

func getTopDomain(host string) string {
	domain, err := getTopDomainByExt(host, ".com")
	if err == nil {
		return domain
	}

	domain, err = getTopDomainByExt(host, ".net")
	if err == nil {
		return domain
	}

	return ""
}

func getTopDomainByExt(host, ext string) (string, error) {
	if strings.HasSuffix(host, ext) {
		return host[strings.LastIndex(strings.TrimSuffix(host, ext), ".")+1:], nil
	}
	return host, &Error{"top domain not found"}
}
