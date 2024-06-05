package editor

import (
	"fmt"
	"net/http"
	"os/exec"
	"runtime"

    "os"
)

type Server struct {
    Addr string
    Server *http.Server
    Fs http.Handler
    File string
}

func (s *Server) Start() Server {
    mux := http.NewServeMux()

    s.Fs = http.FileServer(http.Dir("./src/editor/public"))

    mux.HandleFunc("/css/", func(w http.ResponseWriter, r *http.Request) {
        s.Fs.ServeHTTP(w, r)
    })

    mux.HandleFunc("/js/", func(w http.ResponseWriter, r *http.Request) {
        s.Fs.ServeHTTP(w, r)
    })

	mux.HandleFunc("/editor", func(w http.ResponseWriter, r *http.Request) {
        
        query := r.URL.Query()
        if query.Get("file") == "" {
            http.Redirect(w, r, "/?file=" + s.File, http.StatusSeeOther)
            return
        }

        // serves /public/editor.html
        http.ServeFile(w, r, "./src/editor/public/editor.html")
    })

    // get file content
    mux.HandleFunc("/api/file", sendFileContent)
    // heartbeat
    mux.HandleFunc("/api/heartbeat", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "ok")
    })


	go func() {
        openBrowser("http://localhost:" + s.Addr + "/editor?file=" + s.File)
	}()

    s.Server = &http.Server{Addr: ":" + s.Addr, Handler: mux}

    go func() {
        if err := s.Server.ListenAndServe(); err != nil {
            // removes if the server was closed
            if err == http.ErrServerClosed {
                return
            }
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


func sendFileContent(w http.ResponseWriter, r *http.Request) {
    if r.Method == "GET" {
        r.ParseForm()
        file := r.FormValue("file")



        if file == "" {
            fmt.Fprintf(w, "No file specified")
            return
        }

        content, err := os.ReadFile("./src/macros/" + file)
        if err != nil {
            fmt.Fprintf(w, "Error reading file: %v", err)
            return
        }
        fmt.Fprintf(w, string(content))
    }
}
