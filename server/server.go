package server

import (
	"net/http"

	"code.cloudfoundry.org/lager"
	"github.com/pivotal-cf/brokerapi"
)

var logger = lager.NewLogger("server")

func Run(broker brokerapi.ServiceBroker) error {
	httpHandler := brokerapi.New(broker, logger, brokerapi.BrokerCredentials{
		Username: "",
		Password: "",
	})

	return http.ListenAndServe(":8080", httpHandler)
}
