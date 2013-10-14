package main

import (
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"flag"
	"fmt"
	"hash"
	"log"
	"net/http"
	"regexp"
	"text/template"
)

const (
	INFO_URL  = "http://day-to-day-stuff.blogspot.com/2007/10/announcement-version-99-does-not-exist.html"
	EMPTY_JAR = "UEsDBAoAAAAAAME+SDiyfwLuGQAAABkAAAAUAAQATUVUQS1JTkYvTUFOSUZFU1QuTUb+ygAATWFuaWZlc3QtVmVyc2lvbjogMS4wDQoNClBLAQIKAAoAAAAAAME+SDiyfwLuGQAAABkAAAAUAAQAAAAAAAAAAAAAAAAAAABNRVRBLUlORi9NQU5JRkVTVC5NRv7KAABQSwUGAAAAAAEAAQBGAAAATwAAAAAA"

	POM_XML = `<?xml version="1.0" encoding="iso-8859-1"?>
<project>
	<modelVersion>4.0.0</modelVersion>
	<groupId>{{.GroupId}}</groupId>
	<artifactId>{{.ArtifactId}}</artifactId>
	<name>{{.Name}}</name>
	<version>99.0-does-not-exist</version>
	<description>
		This is a generated pom. Version 99.0-does-not-exist is a dummy implementation with actually does nothing and has no dependencies. 
		VERSION 99.0-does-not-exist IS NOT IN ANY WAY AFFILIATED WITH THE ORIGINAL DEVELOPERS of {{.GroupId}}.
	</description>
	<url>{{.InfoURL}}</url>
</project>`

	INDEX_HTML = `<html>
<head>
<title>Version 99 Does Not Exist</title>
</head>
<body>
<h1>Version 99 Does Not Exist</h1>
<p>Please see <a href="http://day-to-day-stuff.blogspot.com/2007/10/announcement-version-99-does-not-exist.html">my blog</a> to read why I created Version 99 Does Not Exist and its predecessor no-commons-logging.</p>
<p>Version 99 Does Not Exist emulates a Maven 2 repository and serves empty jars for any valid package that has version number <i>99.0-does-not-exist</i>. It also generates poms, <span style="text-decoration: line-through">metadata files</span> (removed since 2.0) and of course the appropriate hashes.</p>
<p>For example the following links will give an <a href="http://no-commons-logging.zapto.org/mvn2/commons-logging/commons-logging/99.0-does-not-exist/commons-logging-99.0-does-not-exist.jar">empty jar</a>, its <a href="http://no-commons-logging.zapto.org/mvn2/commons-logging/commons-logging/99.0-does-not-exist/commons-logging-99.0-does-not-exist.pom">pom</a> and the <a href="http://no-commons-logging.zapto.org/mvn2/commons-logging/commons-logging/maven-metadata.xml"><span style="text-decoration: line-through">maven metadata</span></a> for commons-logging.</p>
<p><a href="version99.rb">Vesion 99 Does Not Exist source code</a> (rb file, 4Kb, MIT license). See <a href="http://day-to-day-stuff.blogspot.com/2007/10/announcement-version-99-does-not-exist.html">my blog</a> for installation instructions.</p>
</body>
</html>`

	NOT_FOUND_HTML = `<html>
<body>
<h1>Version 99 Does Not Exist (Error 404)</h1>
<h2>Not Found: {{.URL}}
<p>
<a href="{{.InfoURL}}">Version 99 Does Not Exist</a> is a virtual Maven2 repository. 
It generates jars and poms for for any artifact with version '99.0-does-not-exist'.
</p>
</body>
</html>
`
)

// --------------------------------------------------------------------

type notFound struct {
	URL     string
	InfoURL string
}

// --------------------------------------------------------------------

type match struct {
	GroupId    string
	ArtifactId string
	Name       string
	Ext        string
	Digest     string
	InfoURL    string
}

// --------------------------------------------------------------------

var (
	addr             = flag.String("addr", ":8080", "Host:port on which to listen")
	notFoundTemplate = template.Must(template.New("notFound").Parse(NOT_FOUND_HTML))
	pomTemplate      = template.Must(template.New("pom").Parse(POM_XML))
	emptyJar         []byte
)

// --------------------------------------------------------------------

func init() {
	var err error

	emptyJar, err = base64.StdEncoding.DecodeString(EMPTY_JAR)
	if err != nil {
		log.Fatalf("Error decoding empty jar. %v", err)
	}
}

// --------------------------------------------------------------------

func writeNotFound(w http.ResponseWriter, r *http.Request) {
	var buf bytes.Buffer
	err := notFoundTemplate.Execute(&buf, &notFound{URL: r.RequestURI, InfoURL: INFO_URL})
	if err != nil {
		writeError(err, w, r)
		return
	}

	w.WriteHeader(http.StatusNotFound)
	w.Header().Set("Content-Type", "text/html")
	w.Write(buf.Bytes())
}

// --------------------------------------------------------------------

func writeData(data []byte, contentType string, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", contentType)
	w.Write(data)
	logRequest(http.StatusOK, r)
}

// --------------------------------------------------------------------

func writeDigest(data []byte, digest string, w http.ResponseWriter, r *http.Request) {
	d, err := makeDigest(data, digest)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logRequest(http.StatusInternalServerError, r)
		return
	}

	writeData([]byte(d), "text/plain", w, r)
}

// --------------------------------------------------------------------

func writeError(err error, w http.ResponseWriter, r *http.Request) {
	writeStatus(http.StatusInternalServerError, w, r)
	log.Print(err)
}

// --------------------------------------------------------------------

func writeStatus(status int, w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(status)
	logRequest(status, r)
}

// --------------------------------------------------------------------

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		writeStatus(http.StatusBadRequest, w, r)
		return
	}

	if r.RequestURI == "/" || r.RequestURI == "/index.html" {
		writeData([]byte(INDEX_HTML), "text/html", w, r)
		return
	}

	m := matchURI(r.RequestURI)
	if m == nil {
		writeNotFound(w, r)
		return
	}

	switch m.Ext {
	case "jar":
		switch m.Digest {
		case "":
			writeData(emptyJar, "application/jar", w, r)
			return
		case "sha1", "md5":
			writeDigest(emptyJar, m.Digest, w, r)
			return
		}
	case "pom":
		var buf bytes.Buffer
		err := pomTemplate.Execute(&buf, m)
		if err != nil {
			writeError(err, w, r)
			return
		}

		switch m.Digest {
		case "":
			writeData(buf.Bytes(), "text/xml", w, r)
			return
		case "sha1", "md5":
			writeDigest(buf.Bytes(), m.Digest, w, r)
			return
		}
	}

	writeNotFound(w, r)
}

// --------------------------------------------------------------------

func main() {
	flag.Parse()
	http.HandleFunc("/", handler)
	log.Printf("Listening on %s", *addr)
	log.Fatal(http.ListenAndServe(*addr, nil))
}

// --------------------------------------------------------------------

func logRequest(status int, r *http.Request) {
	log.Printf("%d %s", status, r.RequestURI)
}

// --------------------------------------------------------------------

func makeDigest(data []byte, digest string) (string, error) {
	var h hash.Hash

	switch digest {
	case "md5":
		h = md5.New()
	case "sha1":
		h = sha1.New()
	}

	if h == nil {
		return "", fmt.Errorf("Invalid hash extension")
	}

	_, err := h.Write(data)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

// --------------------------------------------------------------------

func matchURI(uri string) *match {
	re := regexp.MustCompile("/mvn2/(.+)/([^/]+)/99\\.0-does-not-exist/(.+)-99\\.0-does-not-exist\\.(jar|pom)(\\.(sha1|md5))?")
	m := re.FindStringSubmatch(uri)

	// no match -> return
	if m == nil {
		return nil
	}

	// artifactId != package name -> return
	if m[2] != m[3] {
		return nil
	}

	return &match{GroupId: m[1], ArtifactId: m[2], Name: m[3], Ext: m[4], Digest: m[6], InfoURL: INFO_URL}
}
