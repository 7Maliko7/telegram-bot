package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	OpsChat = promauto.NewCounter(prometheus.CounterOpts{
		Name: "chat_count",
		Help: "The total number of chats",
	})
)

var (
	OpsScenario = promauto.NewCounter(prometheus.CounterOpts{
		Name: "scenario_count",
		Help: "The total number of finished scenarios",
	})
)
