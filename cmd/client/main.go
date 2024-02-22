/*
sysverity
Copyright (C) 2024  mutablefigment

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/
package main

import (
	"bytes"
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mutablefigment/sysverity/v2/pkg/utils"
	"github.com/pelletier/go-toml/v2"
)

func main() {

	// FIXME: add -i flag for interactive
	// and one that does a simple setup automatically

	// set index to zero
	index = 0

	p := tea.NewProgram(initialModel())

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

type (
	errMsg error
)

type model struct {
	viewport    viewport.Model
	messages    []string
	commands    []string
	textarea    textarea.Model
	senderStyle lipgloss.Style
	err         error
}

// TODO: move to external file!
type ConfigFile struct {
	Path           string
	Pubkey         ed25519.PublicKey
	Privkey        ed25519.PrivateKey
	Hashchaintoken []string
	IsRegistered   bool
}

type Config struct {
	Host    string
	Mbrhash string
	File    ConfigFile
}

var (
	Globalconf Config
	index      int
)

func initialModel() model {
	ta := textarea.New()
	ta.Placeholder = "Type help for commands ..."
	ta.Focus()

	ta.Prompt = "┃ "
	ta.CharLimit = 280

	ta.SetWidth(30)
	ta.SetHeight(3)

	// Remove cursor line styling
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()

	ta.ShowLineNumbers = false

	vp := viewport.New(100, 15)
	vp.SetContent(`Type in command, help`)

	ta.KeyMap.InsertNewline.SetEnabled(false)

	return model{
		textarea:    ta,
		messages:    []string{},
		viewport:    vp,
		senderStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("5")),
		err:         nil,
	}
}

func (m model) Init() tea.Cmd {
	return textarea.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		tiCmd tea.Cmd
		vpCmd tea.Cmd
	)

	m.textarea, tiCmd = m.textarea.Update(msg)
	m.viewport, vpCmd = m.viewport.Update(msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			fmt.Println(m.textarea.Value())
			return m, tea.Quit

		case tea.KeyDown:
			if len(m.commands) < index+1 {
				return m, nil
			}
			if index < 0 {
				index = 0
			}

			// TODO: if the index is zero because the user pressed the up key to get it to zero
			// how do we skip having to press the down button twice?

			m.textarea.Reset()
			// Do the wraparound using module of the lenght of the array
			m.textarea.InsertString(m.commands[(index)%len(m.commands)])
			index += 1

		case tea.KeyUp, tea.KeyCtrlP:
			if index < 0 {
				return m, nil
			}
			m.textarea.Reset()
			m.textarea.InsertString(m.commands[(index)%len(m.commands)])
			index -= 1

		// Clear screen
		case tea.KeyCtrlL:
			m.viewport.SetContent("")
			m.messages = make([]string, 0)

		case tea.KeyEnter:
			m.commands = append(m.commands, m.textarea.Value())
			index = len(m.commands) - 1

			m.messages = append(m.messages, m.senderStyle.Render("λ> ")+m.textarea.Value())

			messages, err := rumCmd(m.textarea.Value())
			m.err = err
			for _, message := range messages {
				m.messages = append(m.messages, m.senderStyle.Render(message))
			}

			// FIXME: handle error messages better!
			if err != nil {
				panic(err)
			}

			// update the viewport
			m.viewport.SetContent(strings.Join(m.messages, "\n"))

			m.textarea.Reset()
			m.viewport.GotoBottom()
		}

	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	return m, tea.Batch(tiCmd, vpCmd)
}

func (m model) View() string {
	return fmt.Sprintf(
		"%s\n\n%s",
		m.viewport.View(),
		m.textarea.View(),
	) + "\n\n"
}

func dump() ([]string, error) {
	var output []string
	bytes, err := utils.UniversalMBRDump()
	if err != nil {
		panic(err)
	}

	hash_bytes := sha256.Sum256(bytes[:])
	hash_hex := hex.EncodeToString(hash_bytes[:])

	Globalconf.Mbrhash = hash_hex

	output = append(output, fmt.Sprintf("MBR SHA256 sum is >> %s", hash_hex))

	return output, nil
}

/*

Setup workflow:

set hostname 	-> setup the remote blockhost
ping 			-> test connection
init 			-> generate ed25519 keypair
register 		-> register with the blockhost and setup the linkchain tokens
				   TODO: save linkchain tokens to file with encrypte privatekey
dump 			-> create a hash of the mbr
genisis 		-> upload the initial hash of the mbr signed with keypair and use up first hashchain token
save			-> save a config file

Update workflow:

load 			-> load config file & keys
ping 			-> test connection
dump			-> dump new hash of mbr
update 			-> recalculate last next hashchain token, sign update with key and use up hashchain token
				   also update the config

check workflow:

load
ping
dump
check			-> check last hash against current calculated one to see if everything is correct
*/

func rumCmd(cmd string) ([]string, error) {

	cmdParts := strings.Split(strings.ToLower(cmd), " ")
	output := []string{}

	switch cmdParts[0] {

	case "i", "init":
		pubkey, privkey, err := ed25519.GenerateKey(nil)
		if err != nil {
			panic(err)
		}

		Globalconf.File.Path = "./config.ini"
		Globalconf.File.Privkey = privkey
		Globalconf.File.Pubkey = pubkey

		pubkey_hex := hex.EncodeToString(pubkey)

		output = append(output, fmt.Sprintf("Generated ed25519 pubkey is %s", pubkey_hex))

	case "w", "write":
		if Globalconf.File.Path == "" {
			return []string{"No config changes were made!", "Initialize the client with init first!"}, nil
		}

		b, err := toml.Marshal(Globalconf)
		if err != nil {
			return []string{"Failed to marshall config!"}, err
		}
		f, err := os.Create(Globalconf.File.Path)
		if err != nil {
			panic(err)
		}
		defer f.Close()

		_, err = f.Write(b)
		if err != nil {
			panic(err)
		}

		output = append(output, fmt.Sprintf("Wrote config file to path %s", Globalconf.File.Path))

	case "l", "load":
		if _, err := os.Stat("./config.ini"); err != nil {
			return []string{"Couldn't find a config file, please save first with <write>"}, err
		}

		data, err := os.ReadFile("./config.ini")
		if err != nil {
			panic(err)
		}

		err = toml.Unmarshal(data, &Globalconf)
		if err != nil {
			panic(err)
		}

		output = append(output, "Loaded config file! Check now!")

	// Dump the MBR hash/
	case "d", "dump":
		dstr, err := dump()
		if err != nil {
			panic(err)
		}
		output = dstr

	case "c", "check":
		output = append(output, " .......")
		output = append(output, fmt.Sprintf("Filepath : %s", Globalconf.File.Path))
		output = append(output, fmt.Sprintf("Hostname : %s", Globalconf.Host))
		output = append(output, fmt.Sprintf("MBR Hash : %s", Globalconf.Mbrhash))
		output = append(output, fmt.Sprintf("Pubkey   : %s", hex.EncodeToString(Globalconf.File.Pubkey)))
		output = append(output, " .......")

		if !check(&output) {
			output = append(output, "Check the following steps: !!!!", "")
		}

	case "r", "register":

		Globalconf.File.IsRegistered = true
		if !check(&output) {
			Globalconf.File.IsRegistered = false
			return output, nil
		}

		// TODO: generate hash chain
		magic := "magicpleasechangeme"

		var hashlist [][32]byte

		ihash := sha256.Sum256([]byte(magic))
		Globalconf.File.Hashchaintoken = append(Globalconf.File.Hashchaintoken, hex.EncodeToString(ihash[:]))
		hashlist = append(hashlist, ihash)

		count := 10
		for i := 0; i < count; i++ {
			chash := sha256.Sum256(hashlist[len(hashlist)-1][:])
			Globalconf.File.Hashchaintoken = append(Globalconf.File.Hashchaintoken, hex.EncodeToString(chash[:]))
			output = append(output, hex.EncodeToString(chash[:]))
			hashlist = append(hashlist, chash)
		}

		blhash := sha256.Sum256(hashlist[len(hashlist)-1][:])
		output = append(output, fmt.Sprintf("Hashchain token is: %s", hex.EncodeToString(blhash[:])))

		output = append(output, "Registering with server!")

		url := fmt.Sprintf("http://%s/register", Globalconf.Host)

		// FIXME: marshall to json from globalconf
		var jsonStr = []byte(`{"Pubkey": "` + hex.EncodeToString(Globalconf.File.Pubkey) + `"}`)
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
		if err != nil {
			return []string{"Failed to create POST request for registration"}, err

		}
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return []string{"Failed to Send POST Request for registration"}, err
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		output = append(output, string(body))

	case "g", "genisis":
		Globalconf.File.IsRegistered = true
		if !check(&output) {
			Globalconf.File.IsRegistered = false
			return output, nil
		}

	// Set remote host
	// FIXME: add subcommands to set configpath or hostname
	case "s", "set":
		if len(cmdParts) < 2 {
			return []string{"Please provide a host to set for server!"}, nil
		}

		rh := cmdParts[1]
		output = append(output, "Adding remote host: "+rh+"")
		Globalconf.Host = rh

	// Try to ping the set hostname
	case "p", "ping":
		if Globalconf.Host == "" {
			return []string{"Please provide a hostname with <<set>> first!"}, nil
		}

		conn, err := net.Dial("tcp", Globalconf.Host)
		if err != nil {
			output = append(output, fmt.Sprintf("Failed to ping %s", Globalconf.Host))
			return output, err
		}
		defer conn.Close()

		output = append(output, fmt.Sprintf("ping to %s worked!", Globalconf.Host))

	default:
		output = append(output, "Keybindings:")
		output = append(output, "CTRL-L  Clears the screen")
		output = append(output, "CTRL-P  Get last entered command")
		output = append(output, "                     *******")

		output = append(output, "(d)ump                  Dumps the hash of the MBR")
		output = append(output, "(s)et <hostname>:<port> Sets a remote host")
		output = append(output, "(p)ing                  Pings the host set with set")
		output = append(output, "(u)pdate                Updates the blockchain of boothashes, needs hostname set!")
		output = append(output, "(l)oad                  Load the private key to sign new bootloader hash")
		output = append(output, "(w)rite                 Save the current config to a config.ini file")
		output = append(output, "(i)nit                  Generate a ed25519 keypair")
		output = append(output, "(r)egister              Registers the boothash with the server from <set>")
		output = append(output, "(c)heck                 Check the current config and give some hints")
		output = append(output, "                     *******")

		output = append(output, "Hint: All commands can be shortened to just the first letter!")
	}

	return output, nil
}

func check(msgs *[]string) bool {

	var result bool = true
	result = result && (Globalconf.Host != "")
	if !result {
		*msgs = append(*msgs, "1. Setup a hostname with the <set> command, see <help> for more info!")
	}

	result = result && (Globalconf.File.Path != "")
	// if !result {
	// 	*msgs = append(*msgs, "Setup a hostname with the <set> command, see <help> for more info!")
	// }

	result = result && (Globalconf.File.Privkey != nil)
	result = result || (Globalconf.File.Pubkey != nil)
	if !result {
		*msgs = append(*msgs, "2. Generate a ed25519 keypair with <init>, see <help> for more info!")
	}

	result = result || (Globalconf.Mbrhash != "")
	if !result {
		*msgs = append(*msgs, "3. Dump the MBR hash with <dump>, see <help> for more info!")
	}

	result = result && (Globalconf.File.Hashchaintoken != nil)
	if !result || !Globalconf.File.IsRegistered {
		*msgs = append(*msgs, "4. Register with the host using <register>, see <help> for more info!")
		return result
	}

	return result

}
