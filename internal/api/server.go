package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/guicybercode/systui/internal/exports"
	"github.com/guicybercode/systui/internal/system"
)

type Server struct {
	port int
}

func NewServer(port int) *Server {
	return &Server{port: port}
}

func (s *Server) Start() error {
	http.HandleFunc("/metrics", s.handleMetrics)
	http.HandleFunc("/processes", s.handleProcesses)
	http.HandleFunc("/services", s.handleServices)
	http.HandleFunc("/network", s.handleNetwork)
	http.HandleFunc("/report", s.handleReport)

	addr := fmt.Sprintf(":%d", s.port)
	log.Printf("API server starting on %s", addr)
	return http.ListenAndServe(addr, nil)
}

func (s *Server) handleMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	cpu, _ := system.GetCPUMetrics()
	mem, _ := system.GetMemoryMetrics()
	disk, _ := system.GetDiskMetrics()

	metrics := map[string]interface{}{
		"timestamp": time.Now(),
		"cpu":       cpu,
		"memory":    mem,
		"disk":      disk,
	}

	json.NewEncoder(w).Encode(metrics)
}

func (s *Server) handleProcesses(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	processes, err := system.GetProcesses()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(processes)
}

func (s *Server) handleServices(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	services, err := system.GetServices()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(services)
}

func (s *Server) handleNetwork(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	stats, err := system.GetNetworkStats()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	connections, err := system.GetNetworkConnections()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	network := map[string]interface{}{
		"stats":       stats,
		"connections": connections,
	}

	json.NewEncoder(w).Encode(network)
}

func (s *Server) handleReport(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	report, err := exports.GenerateReport()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(report)
}
