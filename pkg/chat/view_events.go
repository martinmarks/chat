package chat

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"os"

	"crypto/x509"
	"io/ioutil"

	"strings"

	"github.com/jroimartin/gocui"
)

var (
	connection net.Conn
	reader     *bufio.Reader
	writer     *bufio.Writer
)

// RunApplication runs the main client chat app. It auths with the server, connects user, and shows messages.
func RunApplication(g *gocui.Gui, v *gocui.View) error {
	authenticate(g, v)
	setUserConnected(g, v)
	go readNewMessages(g, v)
	return nil
}

// Quit closes out the connection and exits the app
func Quit(g *gocui.Gui, v *gocui.View) error {
	if connection != nil {
		err := connection.Close()
		if err != nil {
			log.Fatal(err)
		}
	}

	return gocui.ErrQuit
}

// SendMessage writes a message to the hub
func SendMessage(g *gocui.Gui, v *gocui.View) error {
	msg := UserMessage(strings.TrimSuffix(v.Buffer(), "\n"))

	_, err := writer.WriteString(msg.ToString())
	if err != nil {
		return err
	}
	err = writer.Flush()
	if err != nil {
		return err
	}
	g.Update(func(g *gocui.Gui) error {
		v.Clear()
		v.SetCursor(0, 0)
		v.SetOrigin(0, 0)
		return nil
	})
	return nil
}

// hard close the app with a display message
func closeWithError(g *gocui.Gui, message string) {
	g.Close()
	if connection != nil {
		err := connection.Close()
		if err != nil {
			log.Println(err.Error())
		}
	}
	log.Println(message)
	os.Exit(1)
}

// authenticate and connect to the server hub over TLS
func authenticate(g *gocui.Gui, v *gocui.View) {
	pool := x509.NewCertPool()
	severCert, err := ioutil.ReadFile("./selfsigned.crt")
	if err != nil {
		closeWithError(g, "Could not load server certificate!")
	}
	pool.AppendCertsFromPEM(severCert)

	conn, err := tls.Dial("tcp", "127.0.0.1:5000", &tls.Config{RootCAs: pool})
	if err != nil {
		closeWithError(g, fmt.Sprintf("Could not connection to chat Hub with error: %s", err.Error()))

	}
	connection = conn
	reader = bufio.NewReader(connection)
	writer = bufio.NewWriter(connection)
}

func setUserConnected(g *gocui.Gui, v *gocui.View) {
	msg := ConnectMessage(strings.TrimSuffix(v.Buffer(), "\n"))

	_, err := writer.WriteString(msg.ToString())
	if err != nil {
		fmt.Printf("error writing string %s\n", err.Error())
	}

	err = writer.Flush()
	if err != nil {
		fmt.Printf("error flushing %s\n", err.Error())
	}

	g.SetViewOnTop("messages")
	g.SetViewOnTop("users")
	g.SetViewOnTop("input")
	g.SetCurrentView("input")
}

func readNewMessages(g *gocui.Gui, v *gocui.View) {
	messagesView, _ := g.View("messages")
	usersView, _ := g.View("clientIDList")

	for {
		data, err := reader.ReadString('\n')
		if err != nil {
			closeWithError(g, fmt.Sprintf("Hub server may have shut down. Ensure the Hub server is running and then reconnect: %s", err.Error()))
		}
		message := ParseMessage(strings.TrimSpace(data))
		switch {
		case message.Type == TypeClientListMessage:
			clientsSlice := strings.Split(message.Text, " ")
			clientsCount := len(clientsSlice)
			var clients string
			for _, client := range clientsSlice {
				clients += client + "\n"
			}
			g.Update(func(g *gocui.Gui) error {
				usersView.Title = fmt.Sprintf(" %d Clients:", clientsCount)
				usersView.Clear()
				fmt.Fprintln(usersView, clients)
				return nil
			})
		default:
			g.Update(func(g *gocui.Gui) error {
				fmt.Fprintln(messagesView, message.DisplayString())
				return nil
			})
		}
	}

}
