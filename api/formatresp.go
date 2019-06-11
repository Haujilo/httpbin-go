package api

import "net/http"

const robotTxt = `User-agent: *
Disallow: /deny
`

func RobotTxtHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	w.Write([]byte(robotTxt))
}

const jsonBody = `{
  "title": "Sample Slide Show",
  "date": "date of publication",
  "author": "Yours Truly",
  "slides": [
    {
      "type": "all",
      "title": "Wake up to WonderWidgets!"
    },
    {
      "type": "all",
      "title": "Overview",
      "items": [
        "Why <em>WonderWidgets</em> are great",
        "Who <em>buys</em> WonderWidgets"
      ]
    }
  ]
}
`

func JsonHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(jsonBody))
}
