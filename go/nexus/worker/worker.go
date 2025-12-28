package main

import (
	"log"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/contrib/envconfig"
	"go.temporal.io/sdk/worker"

	"temporal-order-management/activities"
	"temporal-order-management/app"
	"temporal-order-management/nexus/handler"
	"temporal-order-management/workflows"

	"github.com/nexus-rpc/sdk-go/nexus"
)

func main() {
	co, err := envconfig.LoadClientOptions(envconfig.LoadClientOptionsRequest{
		EnvLookup: EnvLookupMap{
			"TEMPORAL_ADDRESS":              app.GetEnv("TEMPORAL_NEXUS_ADDRESS", ""),
			"TEMPORAL_NAMESPACE":            app.GetEnv("TEMPORAL_NEXUS_NAMESPACE", ""),
			"TEMPORAL_API_KEY":              app.GetEnv("TEMPORAL_NEXUS_API_KEY", ""),
			"TEMPORAL_TLS_CLIENT_CERT_PATH": app.GetEnv("TEMPORAL_NEXUS_TLS_CLIENT_CERT_PATH", ""),
			"TEMPORAL_TLS_CLIENT_KEY_PATH":  app.GetEnv("TEMPORAL_NEXUS_TLS_CLIENT_KEY_PATH", ""),
		},
	})
	if err != nil {
		log.Fatalln("error loading default client options", err)
	}

	c, err := client.Dial(co)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	log.Printf("✅ Client connected to %v in namespace '%v'", co.HostPort, co.Namespace)
	defer c.Close()

	w := worker.New(c, app.GetEnv("TEMPORAL_NEXUS_TASK_QUEUE", "shipping"), worker.Options{})
	service := nexus.NewService(app.ShippingServiceName)
	err = service.Register(handler.ShippingOperation)
	if err != nil {
		log.Fatalln("Unable to register operations", err)
	}
	w.RegisterNexusService(service)
	w.RegisterWorkflow(workflows.ShippingWorkflow)
	w.RegisterActivity(activities.ShipOrder)

	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker", err)
	}
}

type EnvLookupMap map[string]string

func (e EnvLookupMap) Environ() []string {
	ret := make([]string, 0, len(e))
	for k, v := range e {
		ret = append(ret, k+"="+v)
	}
	return ret
}

func (e EnvLookupMap) LookupEnv(key string) (string, bool) {
	v, ok := e[key]
	if v == "" {
		ok = false
	}
	return v, ok
}
