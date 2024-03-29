package yuwiki

import (
	"bytes"
	"github.com/dchest/captcha"
	"net/http"
	"path"
	"strings"
	"time"
)

func Serve(w http.ResponseWriter, r *http.Request, id, ext, lang string, download bool, width, height int) error {
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	var content bytes.Buffer
	var err error
	switch ext {
	case ".png":
		w.Header().Set("Content-Type", "image/png")
		err = captcha.WriteImage(&content, id, width, height)
	case ".wav":
		w.Header().Set("Content-Type", "audio/x-wav")
		err = captcha.WriteAudio(&content, id, lang)
	default:
		return captcha.ErrNotFound
	}
	if err != nil {
		return err
	}

	if download {
		w.Header().Set("Content-Type", "application/octet-stream")
	}
	http.ServeContent(w, r, id+ext, time.Time{}, bytes.NewReader(content.Bytes()))
	return nil
}

func ServeHTTP(w http.ResponseWriter, r *http.Request) {
	dir, file := path.Split(r.URL.Path)
	ext := path.Ext(file)
	id := file[:len(file)-len(ext)]
	log.Debug("file : " + file)
	log.Debug("ext : " + ext)
	log.Debug("id : " + id)
	if ext == "" || id == "" {
		http.NotFound(w, r)
		return
	}
	log.Debug("reload : " + r.FormValue("reload"))
	if r.FormValue("reload") != "" {
		captcha.Reload(id)
	}
	lang := strings.ToLower(r.FormValue("lang"))
	download := path.Base(dir) == "download"
	if Serve(w, r, id, ext, lang, download, captcha.StdWidth, captcha.StdHeight) == captcha.ErrNotFound {
		http.NotFound(w, r)
	}
}
