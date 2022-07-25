package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
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

//goland:noinspection RegExpRedundantEscape
func main() {
	// NGINX “combined” log format: http://nginx.org/en/docs/http/ngx_http_log_module.html#log_format
	var myNginxRegex = regexp.MustCompile(`^(?P<remote_addr>[^ ]+)\s-\s(?P<remote_user>[^ ]+)\s\[(?P<time_local>[^\]]+)\]\s"(?P<request>[^"]*)"\s(?P<status>\d{1,3})\s(?P<body_bytes_send>\d+)\s"(?P<http_referer>[^"]+)"\s"(?P<http_user_agent>[^"]+)"`)
	var myDateTimeRegex = regexp.MustCompile("^(?P<day>\\d{1,2})\\/(?P<month>\\w{1,3})\\/(?P<year>\\d{2,4}):(?P<hour>\\d{1,2}):(?P<minute>\\d{1,2}):(?P<second>\\d{1,2})")
	args := os.Args[1:]
	var logPath string
	// l := log.New(os.Stdout, fmt.Sprintf("[%s]", APP), log.Ldate|log.Ltime|log.Lshortfile)
	l := log.New(ioutil.Discard, fmt.Sprintf("[%s]", APP), log.Ldate|log.Ltime|log.Lshortfile)
	l.Printf("# INFO: 'Starting %s version:%s  num args:%d'\n", APP, VERSION, len(args))
	if len(args) == 1 {
		logPath = os.Args[1]
	} else {
		flag.StringVar(&logPath, "f", defaultLogFile, "Path to your log file")
		flag.Parse()
	}
	l.Printf("# INFO: 'about to open log file : %s'\n", logPath)
	file, err := os.Open(logPath)
	if err != nil {
		l.Fatalf("💥💥 ERROR: 'problem opening log at os.Open(*logPath:%s), got error: %v'\n", logPath, err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)
	l.Printf("# INFO: 'about to read log file : %s'\n", logPath)
	scanner := bufio.NewScanner(file)
	numLine, lines := 0, 0
	for scanner.Scan() {
		line := scanner.Text()
		numLine++
		// fmt.Printf("[%8d]\t%s\n", numLine, line)
		match := myNginxRegex.FindStringSubmatch(line)
		nginxCombinedFields := make(map[string]string)
		for i, name := range myNginxRegex.SubexpNames() {
			if i != 0 && name != "" {
				nginxCombinedFields[name] = match[i]
			}
		}
		matchDate := myDateTimeRegex.FindStringSubmatch(nginxCombinedFields["time_local"])
		nginxDateTimeFields := make(map[string]string)
		for j, name := range myDateTimeRegex.SubexpNames() {
			if j != 0 && name != "" {
				nginxDateTimeFields[name] = matchDate[j]
			}
		}

		// let's filter only 200 http status code
		if nginxCombinedFields["status"] == "200" {
			// verb, url, protocol := strings.Split(nginxCombinedFields["request"], " ")
			requestParts := strings.Split(nginxCombinedFields["request"], " ")
			// usually res wll be [HTTP_VERB URL PROTOCOL] like : "GET /index.html HTTP/1.1"
			if len(requestParts) > 1 {
				if requestParts[0] == "GET" {
					if strings.Contains(requestParts[1], "?") {
						urlParts := strings.Split(requestParts[1], "?")
						posLayers := strings.Index(urlParts[1], "LAYERS=")
						if posLayers > 0 {
							lines++
							if len(urlParts) > 1 {
								layersAndForward := urlParts[1][posLayers:]
								layersExtract := strings.Split(layersAndForward, "&")
								// layerList := strings.ReplaceAll(layersExtract[0][7:], "%2C", ", ")
								layers := layersExtract[0][7:]
								var layerList []string
								if strings.Contains(layers, "%2C") {
									layerList = strings.Split(layersExtract[0][7:], "%2C")
								} else {
									if strings.Contains(layers, ",") {
										layerList = strings.Split(layersExtract[0][7:], ",")
									}
								}
								// remove all default layers queries
								for _, layer := range layerList {
									if layer == "bdcad_projets_msgroup" || layer == "perimetres_lim_com_msgroup" {
										// do not print default layers we are not interested in what is always there
									} else {
										fmt.Printf("%s %s %s %s %s:%s %s %s\n", layer, nginxDateTimeFields["day"], nginxDateTimeFields["month"], nginxDateTimeFields["year"], nginxDateTimeFields["hour"], nginxDateTimeFields["minute"], nginxCombinedFields["remote_addr"], nginxCombinedFields["http_referer"])
									}
								}
							}
						}
					}
				}
			}
		}
	}
	l.Printf("# INFO: 'found %d \tlines with status code 200 and Http verb = GET in log file : %s'\n", lines, logPath)
}
