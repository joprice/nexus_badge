package main

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"text/template"
)

const (
	height        = 20
	width         = 130
	ratio         = .60 // proportion of badge taken up by "nexus" label
	labelWidth    = int(ratio * float64(width))
	versionWidth  = width - labelWidth
	labelCenter   = int(labelWidth / 2)
	versionCenter = labelWidth + (versionWidth / 2)
	badgeSvg      = `<?xml version="1.0"?>
  <svg xmlns="http://www.w3.org/2000/svg" width="{{.Width}}" height="{{.Height}}">
  <linearGradient id="b" x2="0" y2="100%">
    <stop offset="0" stop-color="#bbb" stop-opacity=".1"/>
    <stop offset="1" stop-opacity=".1"/>
  </linearGradient>
  <mask id="a">
    <rect width="{{.Width}}" height="{{.Height}}" rx="3" fill="#fff"/>
  </mask>
  <g mask="url(#a)">
    <path fill="#555" d="M0 0h{{.LabelWidth}}v{{.Height}}H0z"/>
    <path fill="#4c1" d="M{{.LabelWidth}} 0h{{.VersionWidth}}v{{.Height}}H{{.LabelWidth}}z"/>
    <path fill="url(#b)" d="M0 0h{{.Width}}v{{.Height}}H0z"/>
  </g>
  <g fill="#fff" text-anchor="middle" font-family="DejaVu Sans,Verdana,Geneva,sans-serif" font-size="11">
    <text x="{{.LabelCenter}}" y="15" fill="#010101" fill-opacity=".3">nexus</text>
    <text x="{{.LabelCenter}}" y="14">nexus</text>
    <text x="{{.VersionCenter}}" y="15" fill="#010101" fill-opacity=".3">{{.Version}}</text>
    <text x="{{.VersionCenter}}" y="14">{{.Version}}</text>
  </g>
  </svg>`
)

var badgeTmpl *template.Template

func init() {
	var err error
	badgeTmpl, err = template.New("badge").Parse(badgeSvg)
	if err != nil {
		log.Fatal(err)
	}
}

func renderBadge(version string) (string, error) {
	data := map[string]interface{}{
		"Height":        height,
		"Width":         width,
		"LabelWidth":    labelWidth,
		"VersionWidth":  versionWidth,
		"Version":       version,
		"VersionCenter": versionCenter,
		"LabelCenter":   labelCenter,
	}
	var buf bytes.Buffer
	if err := badgeTmpl.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func badgeHandler(nexusURL string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//TODO: return cache headers
		request, err := parseArtifactRequest(r)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}

		artifact, err := latest(nexusURL, request)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		//TODO: return not found error from latest instead?
		if artifact == nil {
			http.Error(w, "artifact not found", 404)
			return
		}

		b, err := renderBadge(artifact.Version)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		w.Header().Set("Content-Type", "image/svg+xml")
		io.WriteString(w, b)
	}
}
