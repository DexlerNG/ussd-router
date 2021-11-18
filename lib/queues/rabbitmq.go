package queues

import "log"
import 	"github.com/blackhades/go-amqp-lib/rabbitmq"

func Init() {
	channel := rabbitmq.Connect("")
	if channel == nil {
		log.Panic("Unable to create rabbitMQ Connection")
	}
}