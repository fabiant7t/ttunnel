package main

import (
	"fmt"
	tt "github.com/johnnylee/ttunnel"
	"github.com/johnnylee/tui"
	"io/ioutil"
	"time"
)

func main() {
	err := checkCreateConfig()
	if err != nil {
		fmt.Printf("Error generating configuration file: %v\n", err)
		return
	}
	mainMenu()
}

// checkCreateConfig will prompt the user to create a new
// configuration file if the file doesn't exist.
func checkCreateConfig() error {
	var err error

	// File exists.
	if tt.FileExists(tt.ConfigPath) {
		return nil
	}

	tui.Clear()

	// Prompt for creating a configuration file.
	if !tui.YesNo("Config file not found. Create a new one?") {
		return fmt.Errorf("Config file not created.")
	}

	for {
		tui.Line()

		sc := tt.ServerConfig{}

		fmt.Printf("\n" +
			"ListenAddr is the address that the server will listen on.\n" +
			"It has the format <address>:<port>.\n" +
			"The address can be omitted to listen on all interfaces.\n\n")

		sc.ListenAddr = tui.String("ListenAddr (:1044)")
		if len(sc.ListenAddr) == 0 {
			sc.ListenAddr = ":1044"
		}

		fmt.Printf("\n" +
			"ConnectAddr is the address that clients will connect on.\n" +
			"It has the same format as ListenAddr.\n\n")

		sc.ConnectAddr = tui.StringNotEmpty("ConnectAddr")
		if sc.EncKey, err = tt.RandomBytes(32); err != nil {
			return err
		}

		fmt.Printf("\n" +
			"ListenAddr : " + sc.ListenAddr + "\n" +
			"ConnectAddr: " + sc.ConnectAddr + "\n\n")
		if tui.YesNo("Is this correct?") {
			return tt.WriteServerConfig(sc)
		}
	}
}

func mainMenu() {
	for {
		fmt.Println("\n")
		s := tui.Menu("ttunnel-config",
			nil,
			"l", "List clients",
			"a", "Add client",
			"r", "Remove client",
			"q", "Quit")
		switch s {
		case "l":
			listClients()
		case "a":
			addClient()
		case "r":
			removeClient()
		case "q":
			return
		}
	}
}

func pressEnter() {
	tui.String("Press enter to continue")
}

func showError(err error) {
	fmt.Printf("\nError: %v\n\n", err)
	pressEnter()
}

func listClients() {
	fmt.Printf("\n")
	tui.Line()

	// Get the file list.
	files, err := ioutil.ReadDir(tt.TunnelDir)
	if err != nil {
		showError(err)
		return
	}

	fmt.Println("")
	for _, f := range files {
		fmt.Println(f.Name())
	}
	fmt.Println("")

	pressEnter()
}

func addClient() {
	// Get name and port.
	for {
		fmt.Printf("\n")
		tui.Line()

		fmt.Printf("\n" +
			"The name given for a client must be unique. It is specific\n" +
			"to the client and connection address. A good name might be\n" +
			"something like machine1-machine2-myqsl.\n\n")

		name := tui.StringNotEmpty("Name")

		fmt.Printf("\n" +
			"The client port is the port on the client's local machine\n" +
			"that will be forwarded.\n\n")

		clientPort := tui.Int("Client port")

		fmt.Printf("\n" +
			"This is the address that the client's local port will be\n" +
			"forwarded to.\n\n")

		connectTo := tui.StringNotEmpty("Connect to")

		fmt.Printf("\n" +
			"The number of days until the client's token expires.\n\n")

		expires := tui.Int("Expires")

		fmt.Printf("\n"+
			"Name       : %v\n"+
			"ClientPort : %v\n"+
			"ConnectTo  : %v\n"+
			"Expires    : %v days\n\n",
			name, clientPort, connectTo, expires)

		if tui.YesNo("Is this correct?") {

			// Load server configuration.
			sc, err := tt.ReadServerConfig()
			if err != nil {
				showError(err)
				return
			}

			// Create token.
			te, err := tt.NewTokenHandler(sc.EncKey)
			if err != nil {
				showError(err)
				return
			}

			token := tt.Token{}
			token.Name = name
			token.ConnectAddr = connectTo
			token.Expires = time.Now().Unix() + (expires * 24 * 3600)

			encToken, err := te.Encode(token)
			if err != nil {
				showError(err)
				return
			}

			cc := tt.ClientConfig{}
			cc.Host = sc.ConnectAddr
			cc.Port = int32(clientPort)
			cc.Token = encToken

			if err = tt.WriteClientConfig(name, cc); err != nil {
				showError(err)
				return
			}

			tui.Line()
			fmt.Printf("\n"+
				"Client added.\n"+
				"Configuration file: %v\n\n",
				tt.ClientConfigPath(name))

			pressEnter()
			return
		}
	}

}

func removeClient() {
	fmt.Printf("\n")
	tui.Line()
	name := tui.StringNotEmpty("Name")

	if tui.YesNo("Remove " + name + "?") {
		if err := tt.RemoveClientConfig(name); err != nil {
			showError(err)
		} else {
			tui.Line()
			fmt.Printf("\n"+
				"Client removed: %v\n"+
				"Configuration saved: %v\n\n",
				name, tt.ClientRemovedPath(name))
			pressEnter()
		}
	}
}
