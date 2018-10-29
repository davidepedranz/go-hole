package main

func main() {
	go runPrometheusServer()
	runDNSServer()
}
