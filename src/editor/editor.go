package editor

import (
	"fmt"
	"net/http"
    "os/exec"
	"runtime"
)

type Server struct {
    Addr string
    Server *http.Server
}

func (s *Server) Start() Server {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello, World!")
	})

	go func() {
        openBrowser("http://localhost:" + s.Addr)
	}()

    s.Server = &http.Server{Addr: ":" + s.Addr}

    go func() {
        if err := s.Server.ListenAndServe(); err != nil {
            fmt.Printf("Failed to start server: %v\n", err)
        }
    }()

    return *s

}

func openBrowser(url string) error  { 
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	case "linux":
		cmd = exec.Command("xdg-open", url)
	default:
		fmt.Printf("Please open the browser manually to this URL: %s\n", url)
	}

    return cmd.Start()
}
