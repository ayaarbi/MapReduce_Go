package master

import (
	"encoding/json"
	"fmt"
	"net/http"
	"project/common"
	"sort"
	"strconv"
	"path/filepath"
)

func (m *Master) StartDashboard() {
	// Serve HTML dashboard
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "dashboard/index.html")
	})

	// Serve static assets like script.js
	http.Handle("/dashboard/", http.StripPrefix("/dashboard/", http.FileServer(http.Dir("dashboard"))))

	// JSON state of system
	http.HandleFunc("/data", func(w http.ResponseWriter, r *http.Request) {
		state := m.Snapshot()
		w.Header().Set("Content-Type", "application/json")
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		enc.Encode(state)
	})

	// Top 10 word result (after reduce)
	http.HandleFunc("/result", func(w http.ResponseWriter, r *http.Request) {
		if m.Done < len(m.Tasks) {
			w.WriteHeader(http.StatusAccepted)
			w.Write([]byte("Processing..."))
			return
		}

		wordCounts := make(map[string]int)
		for i := 0; i < 3; i++ {
			filename := fmt.Sprintf("mrtmp.job-res-%d", i)
			kvs, err := common.ReadKeyValuesFromFile(filename)
			if err != nil {
				fmt.Println("Error reading:", filename, err)
				continue
			}
			for _, kv := range kvs {
				n, _ := strconv.Atoi(kv.Value)
				wordCounts[kv.Key] += n
			}
		}

		type Pair struct {
			Word  string `json:"word"`
			Count int    `json:"count"`
		}
		var pairs []Pair
		for k, v := range wordCounts {
			pairs = append(pairs, Pair{k, v})
		}

		sort.Slice(pairs, func(i, j int) bool {
			return pairs[i].Count > pairs[j].Count
		})

		if len(pairs) > 10 {
			pairs = pairs[:10]
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(pairs)
	})

	http.HandleFunc("/debug/files", func(w http.ResponseWriter, r *http.Request) {
	files, _ := filepath.Glob("mrtmp.*")
	json.NewEncoder(w).Encode(files)
})

fmt.Println("Dashboard server listening on :8080")
http.ListenAndServe(":8080", nil)

}
