package api

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

var Response = ApiResponses{
	/* 100 */

	/* 200 */
	Ok: func(w http.ResponseWriter, data ...any) {
		Send(w, 200, data)
	},

	Created: func(w http.ResponseWriter, data ...any) {
		Send(w, 201, data)
	},

	NoContent: func(w http.ResponseWriter, data ...any) {
		Send(w, 204, nil)
	},

	/* 300 */

	/* 400 */
	BadRequest: func(w http.ResponseWriter, data ...any) {
		Send(w, http.StatusBadRequest, data)
	},

	Unauthorized: func(w http.ResponseWriter, data ...any) {
		Send(w, http.StatusUnauthorized, data)
	},

	Forbidden: func(w http.ResponseWriter, data ...any) {
		Send(w, http.StatusForbidden, data)
	},

	NotFound: func(w http.ResponseWriter, data ...any) {
		Send(w, http.StatusNotFound, data)
	},

	TooManyRequests: func(w http.ResponseWriter, data ...any) {
		Send(w, http.StatusTooManyRequests, data)
	},

	/* 500 */
	ServerError: func(w http.ResponseWriter, data ...any) {
		Send(w, http.StatusInternalServerError, data)
	},

	/* Misc */
	Send: func(w http.ResponseWriter, status int, data ...any) {
		Send(w, status, data)
	},

	Gzip: func(w http.ResponseWriter, status int, data ...any) {
		zip(w, status, data...)
	},

	StreamBytes: func(w http.ResponseWriter, status int, bytes []byte, name string) {
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%v\"", name))
		w.WriteHeader(status)
		w.Write(bytes)
	},

	StreamFile: func(w http.ResponseWriter, status int, binPath string, name string) {

		f, err := os.Open(binPath)
		if err != nil {
			//Send server error
			Send(w, http.StatusInternalServerError, nil)
			return
		}
		defer f.Close()

		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%v\"", name))
		w.WriteHeader(status)

		size, err := io.Copy(w, f)
		if err != nil {
			log.Println("Streaming failed:", err)
		}

		log.Printf("Copied %v bytes", size)

		err = os.RemoveAll(filepath.Dir(binPath))
		if err != nil {
			log.Println("Failed to clean up temp dir:", err)
		}
	},
}

func Send(w http.ResponseWriter, status int, data any) {
	if status == http.StatusNoContent {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	switch v := data.(type) {
	// If only one argument passed, unwrap it from slice
	case []any:
		if len(v) == 1 {
			data = v[0]
		}
	}

	var buf bytes.Buffer
	if data != nil {
		err := json.NewEncoder(&buf).Encode(data)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if buf.Len() > 0 {
		w.Write(buf.Bytes())
	}
}

func zip(w http.ResponseWriter, status int, data ...any) {
	var payload any
	if len(data) == 1 {
		payload = data[0]
	} else {
		payload = data
	}

	// Set headers
	w.Header().Set("Content-Encoding", "gzip")
	w.Header().Set("Content-Type", "application/json")

	gz := gzip.NewWriter(w)
	defer gz.Close()

	//json -> gz -> res
	w.WriteHeader(status)
	err := json.NewEncoder(gz).Encode(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
