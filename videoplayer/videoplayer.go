package videoplayer

import (
	"encoding/json"
	"net/http"
	"time"
)

type VideoPlayerBackendConfig struct {
	RootDir string
}

type VideoPlayerBackend struct {
	conf VideoPlayerBackendConfig
	fm   FileManager
}

type FileManager interface {
	GetFiles(cam string, day time.Time) []string
}

func New(conf VideoPlayerBackendConfig, fm FileManager) *VideoPlayerBackend {
	var ret VideoPlayerBackend
	ret.conf = conf
	ret.fm = fm
	return &ret
}

func (v *VideoPlayerBackend) Run(addr string) {
	// Start the video player backend
	http.Handle("/content/", http.FileServer(http.Dir(v.conf.RootDir)))
	http.HandleFunc("/get/date", func(w http.ResponseWriter, r *http.Request) {
		vars := r.URL.Query()
		t, err := time.Parse("2006-01-02", vars.Get("date"))
		if err != nil {
			http.Error(w, "Invalid date", http.StatusBadRequest)
			return
		}
		data := v.fm.GetFiles(vars.Get("cam"), t)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(data)
	})
	http.ListenAndServe(addr, nil)
}
