package agent

import "log"

// monitorThreats aggregates threat detection tasks.
// It can be invoked periodically or as needed.
func monitorThreats() {
	monitorUnauthorizedFileModifications()
	monitorNewProcesses()
}

func monitorUnauthorizedFileModifications() {
	// TODO: Implement OS-specific file system monitoring.
	log.Println("Monitoring unauthorized file modifications (Not yet implemented)")
}

func monitorNewProcesses() {
	// TODO: Implement OS-specific process monitoring.
	log.Println("Monitoring new running processes (Not yet implemented)")
}
