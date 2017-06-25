package api

import (
	"encoding/json"
	"github.com/gomqtt/transport"
	"github.com/labstack/echo"
	"github.com/zengfu/web/broker"
	"net/http"
	"fmt"
	"net/url"
	"crypto/tls"
)

var (
	engine = broker.NewEngine()
)

type Client struct {
	Username string
	Clientid string
	Subtopic []string
}

func LaunchMqtt(protocol string) error {
	server, err := transport.Launch(protocol)
	if err != nil {
		fmt.Println(err)
		return (err)
	}
	engine.Accept(server)
	//defer engine.Close()
	return nil
}
func LaunchMqtts(protocol string) error {
	urlParts, err := url.ParseRequestURI(protocol)
	if err != nil {
		return err
	}
	cer, err := tls.LoadX509KeyPair("/usr/local/apache/conf/2_www.whyengineer.com.crt", "/usr/local/apache/conf/3_www.whyengineer.com.key")
	if err!=nil{
		return err
	}
	config := &tls.Config{Certificates: []tls.Certificate{cer}}
        server, err := transport.NewSecureWebSocketServer(urlParts.Host,config)
        if err != nil {
                fmt.Println(err)
                return (err)
        }
        engine.Accept(server)
        //defer engine.Close()
        return nil
}
func CloseMqtt() {
	engine.Close()
}

//get
func GetAllClients(c echo.Context) error {
	clis := engine.Clients()
	var all string
	for _, cli := range clis {
		var t Client
		t.Clientid = cli.ClientID()
		t.Username = cli.Username
		subs, err := cli.Session().AllSubscriptions()
		if err != nil {
			return err
		}
		for _, sub := range subs {
			t.Subtopic = append(t.Subtopic, sub.Topic)
		}
		b, _ := json.Marshal(t)
		all = all + string(b) + "\n"
	}

	return c.String(http.StatusOK, all)
}
