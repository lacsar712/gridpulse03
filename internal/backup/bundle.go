package backup

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"example.com/gridpulse/internal/model"
)

type Manifest struct {
	App     string    `json:"app"`
	Created time.Time `json:"created"`
	Count   int       `json:"count"`
}

func Bundle(recs []model.Record) ([]byte, error) {
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	man := Manifest{App: "GridPulse Relay", Created: time.Now().UTC(), Count: len(recs)}
	mb, err := json.MarshalIndent(man, "", "  ")
	if err != nil {
		return nil, err
	}
	if err := writeZip(zw, "manifest.json", mb); err != nil {
		return nil, err
	}
	body, err := json.MarshalIndent(recs, "", "  ")
	if err != nil {
		return nil, err
	}
	if err := writeZip(zw, "records.json", body); err != nil {
		return nil, err
	}
	if err := zw.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func writeZip(zw *zip.Writer, name string, body []byte) error {
	w, err := zw.Create(name)
	if err != nil {
		return err
	}
	_, err = w.Write(body)
	return err
}

func ReadManifest(r io.ReaderAt, size int64) (Manifest, error) {
	zr, err := zip.NewReader(r, size)
	if err != nil {
		return Manifest{}, err
	}
	for _, f := range zr.File {
		if f.Name != "manifest.json" {
			continue
		}
		rc, err := f.Open()
		if err != nil {
			return Manifest{}, err
		}
		b, err := io.ReadAll(rc)
		rc.Close()
		if err != nil {
			return Manifest{}, err
		}
		var m Manifest
		if err := json.Unmarshal(b, &m); err != nil {
			return Manifest{}, err
		}
		return m, nil
	}
	return Manifest{}, fmt.Errorf("manifest.json missing")
}
