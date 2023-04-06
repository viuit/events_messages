package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	config "messageDelivery/config"
	"os"
	"time"

	"github.com/fernet/fernet-go"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

type rabbitConfig struct {
	Protocol     string
	Host         string
	Port         int
	User         string
	Password     string
	VHost        string
	exchangeName string
	heartbeat    int
	name_queue   string
	// routingKey   string
	key []*fernet.Key
}

type rabbitSession struct {
	conn            *amqp.Connection
	ch              *amqp.Channel
	messages        <-chan amqp.Delivery
	NotifyConnClose chan *amqp.Error
	NotifyChanClose chan *amqp.Error
	NotifyConfirm   chan amqp.Confirmation
	// headers         map[string]interface{}
	isReady bool
}

var mqConf *rabbitConfig
var confEnv *config.Environment
var confjs map[string]interface{}
var errNotConnected = errors.New("not connected to a server")

func init() {
	opts := config.NewArgs()
	configFile, err := os.Open(opts.ConfigFile)

	confEnv = config.NewEnv()

	if err != nil {
		log.Fatalf("config file load err   #%v ", err)
	}
	log.Printf("Successfully Opened")

	byteValue, _ := ioutil.ReadAll(configFile)
	if err := json.Unmarshal([]byte(byteValue), &confjs); err != nil {
		log.Fatalf("config file decode err   #%v ", err)
	}

}

func main() {
	mqConf = initRabbit(confjs["rabbit"].(map[string]interface{}))
	mqConf.Host = confEnv.Rabbit.Host
	mqConf.Port = confEnv.Rabbit.Port
	mqConf.VHost = confEnv.Rabbit.VHost

	mqSession := &rabbitSession{isReady: false}

	go func() {
		for {
			if !mqSession.isReady {
				mqSession.Create()
			}
			time.Sleep(100 * time.Millisecond)
		}
	}()

	// var forever chan struct{}

	// var messages <-chan amqp.Delivery

	var err error
	// for !mqSession.isReady {
	// 	time.Sleep(100 * time.Millisecond)
	// 	continue
	// }
	mqSession.messages, err = mqSession.Stream()
	if err != nil {
		log.Errorf("Stream error: %v", err)
	}

	log.Infoln("Consume success to RabbitMQ!")

	for d := range mqSession.messages {
		log.Printf("Received a message: %s", d.Body)
		d.Ack(true)
	}

	// <-forever

	// log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	// <-forever
}

func initRabbit(confjs map[string]interface{}) *rabbitConfig {
	conf := &rabbitConfig{
		Protocol:     "amqp",
		User:         confjs["username"].(string),
		exchangeName: confEnv.Rabbit.Exchange,
		heartbeat:    int(confjs["heartbeat"].(float64)),
		name_queue:   confjs["name_queue"].(string),
	}
	conf.key = fernet.MustDecodeKeys(confEnv.CipherKey)
	conf.Password = string(
		fernet.VerifyAndDecrypt([]byte(confjs["password"].(string)),
			time.Duration(0)*time.Hour, conf.key))
	// conf.VHost = fmt.Sprintf("%s.%s", config.Env.Project, config.Rabbit.VHost)
	return conf
}

func (conf *rabbitConfig) URL() string {
	return fmt.Sprintf("%s://%s:%s@%s:%d/%s?heartbeat=%d}",
		conf.Protocol,
		conf.User,
		conf.Password,
		conf.Host,
		conf.Port,
		conf.VHost,
		conf.heartbeat,
	)
}

func (conf *rabbitConfig) URLINFO() string {
	return fmt.Sprintf("%s://%s:%s@%s:%d/%s",
		conf.Protocol,
		conf.User,
		"***",
		conf.Host,
		conf.Port,
		conf.VHost,
	)
}

func (session *rabbitSession) Create() {
	session.Connection(mqConf.URL())

	// defer session.conn.Close()
	// log.Infof("Connected to the RabbitMQ: %s", )

	chOpen := make(chan error)
	go func() {
		var err error
		session.ch, err = session.conn.Channel()
		chOpen <- err
	}()

	select {
	case err := <-chOpen:
		if err != nil {
			log.Errorf("Failed to open a channel RabbitMQ: %v", err)
			return
		}
	case <-time.After(5 * time.Second):
		log.Errorf("Failed to open a channel RabbitMQ: error time After 5 second")
		return
	}
	close(chOpen)

	log.Infoln("Channel success open to RabbitMQ!")

	if err := session.ch.Qos(1, 0, false); err != nil {
		log.Errorf("Failed Qos RabbitMQ: %v", err)
		return
	}
	// session.MessagesChan = make(<-chan amqp.Delivery)

	// session.NotifyConfirm = make(chan amqp.Confirmation, 1)
	session.NotifyChanClose = make(chan *amqp.Error)
	session.ch.NotifyClose(session.NotifyChanClose)
	// session.ch.NotifyPublish(session.NotifyConfirm)
	session.isReady = true
	// return session
}

func (session *rabbitSession) Connection(url string) {
	log.Info("Connecting to RabbitMQ server...")
	for {
		var err error
		session.conn, err = amqp.Dial(url)
		if err != nil {
			log.Errorf("Failed connection to RabbitMQ: %v", err)
			time.Sleep(100 * time.Millisecond)
			continue
		}
		break
	}
	session.NotifyConnClose = make(chan *amqp.Error)
	session.conn.NotifyClose(session.NotifyConnClose)
	log.Infof("Connect to RabbitMQ server is completed, %s!", mqConf.URLINFO())
}

// Stream will continuously put queue items on the channel.
// It is required to call delivery.Ack when it has been
// successfully processed, or delivery.Nack when it fails.
// Ignoring this will cause data to build up on the server.
func (session *rabbitSession) Stream() (<-chan amqp.Delivery, error) {
	for !session.isReady {
		log.Errorf("Consume error %v", errNotConnected)
		time.Sleep(100 * time.Millisecond)
		continue
	}
	return session.ch.Consume(
		mqConf.name_queue,
		"",                                    // Consumer
		false,                                 // Auto-Ack
		false,                                 // Exclusive
		false,                                 // No-local
		false,                                 // No-Wait
		amqp.Table{"x-stream-offset": "last"}, // Args
	)
}

// Close will cleanly shutdown the channel and connection.
func (session *rabbitSession) Close() error {
	err := session.ch.Close()
	if err != nil {
		return err
	}
	err = session.conn.Close()
	if err != nil {
		return err
	}
	return nil
}
