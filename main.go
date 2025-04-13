package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// url is the remote address we want to check.
const url = "https://charm.sh/"

// model represents the state of our application. It includes
// the HTTP status code (if any) and an error variable.
type model struct {
	status int   // HTTP status code returned from the server.
	err    error // Any error encountered during the HTTP request.
}

// statusMsg is a custom message type used to wrap an HTTP status code.
type statusMsg int

// errMsg is a custom message type used to wrap an error encountered during the HTTP request.
type errMsg struct{ err error }

// checkServer performs an HTTP GET request to the URL and returns a tea.Msg,
// which is either a statusMsg (with the HTTP status code) or an errMsg (on error).
func checkServer() tea.Msg {
	// Create an HTTP client with a timeout of 10 seconds.
	c := &http.Client{Timeout: 10 * time.Second}

	// Perform an HTTP GET request.
	res, err := c.Get(url)
	if err != nil {
		// If an error occurs, wrap and return it as an errMsg.
		return errMsg{err}
	}
	// It is best practice to close the response body to avoid resource leaks.
	res.Body.Close()

	// Return the HTTP status code wrapped as a statusMsg.
	return statusMsg(res.StatusCode)
}

// Init is the initialization function required by the Bubble Tea framework.
// It returns an initial command (tea.Cmd) to be executed, which in this case is the checkServer command.
func (m model) Init() tea.Cmd {
	return checkServer
}

// Update handles incoming messages (tea.Msg) and updates the model accordingly.
// It is invoked by the Bubble Tea runtime whenever an event (like a key press or the completion
// of a command) occurs.
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// When we receive a statusMsg, update the model with the HTTP status.
	case statusMsg:
		m.status = int(msg) // Cast our custom statusMsg to an int.
		// We have the desired information, so signal Bubble Tea to quit.
		return m, tea.Quit

	// When we receive an errMsg, update the model with the error.
	case errMsg:
		m.err = msg.err // Correctly assign the underlying error, not the whole struct.
		// Signal to quit the program.
		return m, tea.Quit

	// Handle key press messages.
	case tea.KeyMsg:
		// Allow the user to exit the program by pressing Ctrl+C.
		if msg.Type == tea.KeyCtrlC {
			return m, tea.Quit
		}
	}

	// If any other message types are received, do nothing.
	return m, nil
}

// View renders the output based on the current state of the model.
// It returns a string that is displayed in the terminal.
func (m model) View() string {
	// If there was an error during the HTTP request, display the error.
	if m.err != nil {
		return fmt.Sprintf("\nWe had some trouble: %v\n\n", m.err)
	}

	// Otherwise, build a string indicating that the program is checking the URL.
	s := fmt.Sprintf("Checking %s ... ", url)

	// If a status code is present, display it along with its standard text representation.
	if m.status > 0 {
		s += fmt.Sprintf("%d %s!", m.status, http.StatusText(m.status))
	}

	// Add some line breaks for nice formatting.
	return "\n" + s + "\n\n"
}

// main is the entry point of the program.
// It creates a new Bubble Tea program using the model, runs it, and handles any errors.
func main() {
	// Create a new Bubble Tea program with an empty model.
	p := tea.NewProgram(model{})

	// Run the program. If there is an error during runtime, print it and exit.
	if _, err := p.Run(); err != nil {
		fmt.Printf("Uh oh, there was an error: %v\n", err)
		os.Exit(1)
	}
}
