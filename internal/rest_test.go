package internal

import (
	"log"
	"regexp"
	"testing"
)

func TestName(t *testing.T) {
	macRegex := regexp.MustCompile(`^([0-9A-Fa-f]{2}[:-]){5}([0-9A-Fa-f]{2})$`)
	ipRegex := regexp.MustCompile(`^([0-9]{1,3}\.){3}[0-9]{1,3}$`)

	mac := "14-07-2D-BB-BB-BF"
	ip := "192.168.10.214:9"

	if !macRegex.MatchString(mac) {
		log.Println("error mac:", mac)
		return
	}
	if !ipRegex.MatchString(ip) {
		log.Println("error ip:", ip)
	}
}
