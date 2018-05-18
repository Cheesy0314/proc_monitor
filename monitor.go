package main

import (
	"os/exec"
	"strings"
	"regexp"
)

func report(valid processList) string {
	out, err := exec.Command("sh","-c","ps -ef | grep 'jar' | grep -v \"grep\" | awk {'print $10'}  ").Output()

	if err != nil {
		panic(err)
	}

	re := regexp.MustCompile(`(.:)`)
	re2 := regexp.MustCompile(`(.jar)`)

	data := strings.Split(string(out),"\n")
	available := map[string]bool{}
	for _, value := range data {
		line := re2.ReplaceAllString(re.ReplaceAllString(value, ""),"")
		available[line] = true
	}

	bad := []string{}
	for _,p := range valid {

		if !available[p.name] {
			bad = append(bad,p.name)
		}
	}
	output := strings.Join(bad, ", ")
	return output
}

