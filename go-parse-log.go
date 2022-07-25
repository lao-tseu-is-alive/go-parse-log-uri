package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

// const defaultLogFile = "data/nginx_ssl_access.log"
const (
	VERSION        = "0.1.0"
	APP            = "go-parse-log"
	defaultLogFile = "data/sample.log"
)

func main() {
	//NGINX “combined” log format: http://nginx.org/en/docs/http/ngx_http_log_module.html#log_format
	var myNginxRegex = regexp.MustCompile(`^(?P<remote_addr>[^ ]+)\s-\s(?P<remote_user>[^ ]+)\s\[(?P<time_local>[^\]]+)\]\s"(?P<request>[^"]+)"\s(?P<status>\d{1,3})\s(?P<body_bytes_send>\d+)\s"(?P<http_referer>[^"]+)"\s"(?P<http_user_agent>[^"]+)"`)
	args := os.Args[1:]
	var logPath string
	l := log.New(os.Stdout, fmt.Sprintf("[%s]", APP), log.Ldate|log.Ltime|log.Lshortfile)
	l.Printf("INFO: 'Starting %s version:%s  num args:%d'\n", APP, VERSION, len(args))
	if len(args) == 1 {
		logPath = os.Args[1]
	} else {
		flag.StringVar(&logPath, "f", defaultLogFile, "Path to your log file")
		flag.Parse()
	}
	l.Printf("INFO: 'about to open log file : %s'\n", logPath)
	file, err := os.Open(logPath)
	if err != nil {
		l.Fatalf("💥💥 ERROR: 'problem opening log at os.Open(*logPath:%s), got error: %v'\n", logPath, err)
	}
	defer file.Close()
	l.Printf("INFO: 'about to read log file : %s'\n", logPath)
	scanner := bufio.NewScanner(file)
	lines, sizeInBytes := 0, 0
	for scanner.Scan() {

		line := scanner.Text()
		// fmt.Println(line)
		sizeInBytes += len(line)
		match := myNginxRegex.FindStringSubmatch(line)
		result := make(map[string]string)
		for i, name := range myNginxRegex.SubexpNames() {
			if i != 0 && name != "" {
				result[name] = match[i]
			}
		}
		// let's filter only 200 http status code
		if result["status"] == "200" {
			// verb, url, protocol := strings.Split(result["request"], " ")
			requestParts := strings.Split(result["request"], " ")
			// usually res wll be [HTTP_VERB URL PROTOCOL] like : "GET /index.html HTTP/1.1"
			if len(requestParts) > 1 {
				if requestParts[0] == "GET" {
					if strings.Contains(requestParts[1], "?") {
						urlParts := strings.Split(requestParts[1], "?")
						posLayers := strings.Index(urlParts[1], "LAYERS=")
						if posLayers > 0 {
							lines++
							layersAndForward := urlParts[1][posLayers:]
							layersExtract := strings.Split(layersAndForward, "&")
							fmt.Printf("%s \t %v\n", layersExtract[0], layersAndForward)
						}
					}
				}
			}
		}
	}
	l.Printf("INFO: 'found %d \tlines with status code 200 and Http verb = GET in log file : %s'\n", lines, logPath)
	l.Printf("INFO: 'found %d \tbytes in log file : %s'\n", sizeInBytes, logPath)
	fmt.Printf("%d\t%s\n", lines, logPath)
}
