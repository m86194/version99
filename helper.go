package main

import (
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"fmt"
	"hash"
	"log"
	"net/http"
	"regexp"
)

// --------------------------------------------------------------------

func logRequest(status int, r *http.Request) {
	log.Printf("%d %s", status, r.RequestURI)
}

// --------------------------------------------------------------------

func matchMavenURI(uri string) *maven {
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

	return &maven{GroupId: m[1], ArtifactId: m[2], Name: m[3], Ext: m[4], Digest: m[6], InfoURL: INFO_URL}
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

func writeStatus(status int, w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(status)
	logRequest(status, r)
}
