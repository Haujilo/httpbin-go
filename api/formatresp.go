package api

import "net/http"

const robotTxt = `User-agent: *
Disallow: /deny
`

func RobotTxtHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(robotTxt))
}
