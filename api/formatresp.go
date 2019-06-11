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

const xmlBody = `<?xml version='1.0' encoding='us-ascii'?>

<!--  A SAMPLE set of slides  -->

<slideshow
    title="Sample Slide Show"
    date="Date of publication"
    author="Yours Truly"
    >

    <!-- TITLE SLIDE -->
    <slide type="all">
      <title>Wake up to WonderWidgets!</title>
    </slide>

    <!-- OVERVIEW -->
    <slide type="all">
        <title>Overview</title>
        <item>Why <em>WonderWidgets</em> are great</item>
        <item/>
        <item>Who <em>buys</em> WonderWidgets</item>
    </slide>

</slideshow>
`

func XMLHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/xml")
	w.Write([]byte(xmlBody))
}
