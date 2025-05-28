package main

import (
	"flag"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const (
	defaultUploadPath = "./uploads"
	defaultPort       = "8080"
)

var (
	uploadPath string
	addr       string // Changed from port to addr for consistency
)

func main() {
	flag.StringVar(&uploadPath, "upload-dir", defaultUploadPath, "Directory to store wallpapers")
	flag.StringVar(&addr, "addr", "", "Address to serve (format: [host]:port)")
	flag.Parse()

	if addr == "" {
		addr = ":" + os.Getenv("PORT")
		if addr == ":" {
			addr = ":" + defaultPort
		}
	}

	if !strings.HasPrefix(addr, ":") && !strings.Contains(addr, ":") {
		addr = ":" + addr
	}

	if err := os.MkdirAll(uploadPath, os.ModePerm); err != nil {
		log.Fatalf("Failed to create upload directory: %v", err)
	}

	http.HandleFunc("/", listWallpapers)
	http.HandleFunc("/upload", uploadHandler)
	http.HandleFunc("/delete/", deleteHandler)
	http.HandleFunc("/rename/", renameHandler)
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/download/", downloadHandler)
	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir(uploadPath))))

	log.Printf("Wallpaper manager running on http://%s\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

// New download handler function
func downloadHandler(w http.ResponseWriter, r *http.Request) {
	filename := strings.TrimPrefix(r.URL.Path, "/download/")
	if filename == "" {
		http.Error(w, "Filename required", http.StatusBadRequest)
		return
	}

	filePath := filepath.Join(uploadPath, filename)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.NotFound(w, r)
		return
	}

	// Set headers to force download
	w.Header().Set("Content-Disposition", "attachment; filename="+filename)
	w.Header().Set("Content-Type", "application/octet-stream")

	http.ServeFile(w, r, filePath)
}

// Updated listWallpapers template with multiple file upload
func listWallpapers(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	files, err := os.ReadDir(uploadPath)
	if err != nil {
		http.Error(w, "Failed to read wallpapers directory", http.StatusInternalServerError)
		return
	}

	type Wallpaper struct {
		Name string
		URL  string
	}

	var wallpapers []Wallpaper
	for _, file := range files {
		if !file.IsDir() {
			wallpapers = append(wallpapers, Wallpaper{
				Name: file.Name(),
				URL:  "/images/" + file.Name(),
			})
		}
	}

	tmpl := template.Must(template.New("list").Parse(`
<!DOCTYPE html>
<html>
<head>
    <title>Wallpaper Collection</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; }
        .wallpaper { margin: 10px; display: inline-block; text-align: center; }
        img { max-width: 300px; max-height: 200px; display: block; margin-bottom: 5px; }
        .actions { margin-top: 5px; }
        form { margin: 20px 0; }
        .upload-status { color: green; margin: 10px 0; }
    </style>
</head>
<body>
    <h1>Wallpaper Collection</h1>
    
    <form action="/upload" method="post" enctype="multipart/form-data">
        <h2>Upload Wallpapers</h2>
        <input type="file" name="images" multiple accept="image/*" required>
        <button type="submit">Upload</button>
        <div class="upload-status">You can select multiple files (Ctrl+Click or Shift+Click)</div>
    </form>
    
    <h2>Your Wallpapers</h2>
    {{if .}}
        {{range .}}
        <div class="wallpaper">
            <img src="{{.URL}}" alt="{{.Name}}">
            <div>{{.Name}}</div>
            <div class="actions">
                <a href="/view/{{.Name}}">View</a> | 
                <a href="/download/{{.Name}}">Download</a> | 
                <a href="/rename/{{.Name}}">Rename</a> | 
                <a href="/delete/{{.Name}}" onclick="return confirm('Are you sure?')">Delete</a>
            </div>
        </div>
        {{end}}
    {{else}}
        <p>No wallpapers found. Upload some above!</p>
    {{end}}
</body>
</html>
`))

	if err := tmpl.Execute(w, wallpapers); err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
	}
}

// Updated uploadHandler to support multiple files
func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse multipart form with 50MB limit (for multiple files)
	if err := r.ParseMultipartForm(50 << 20); err != nil {
		http.Error(w, "Total upload size too large (max 50MB)", http.StatusBadRequest)
		return
	}

	// Get all files from the "images" field (now supports multiple)
	files := r.MultipartForm.File["images"]
	if len(files) == 0 {
		http.Error(w, "No files uploaded", http.StatusBadRequest)
		return
	}

	var successCount int
	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			log.Printf("Failed to open uploaded file %s: %v", fileHeader.Filename, err)
			continue
		}

		// Create destination file
		dstPath := filepath.Join(uploadPath, filepath.Base(fileHeader.Filename))
		dst, err := os.Create(dstPath)
		if err != nil {
			log.Printf("Failed to create file %s: %v", dstPath, err)
			file.Close()
			continue
		}

		// Copy the file
		if _, err := io.Copy(dst, file); err != nil {
			log.Printf("Failed to save file %s: %v", dstPath, err)
			file.Close()
			dst.Close()
			continue
		}

		file.Close()
		dst.Close()
		successCount++
	}

	if successCount == 0 {
		http.Error(w, "Failed to save all files", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	filename := strings.TrimPrefix(r.URL.Path, "/delete/")
	if filename == "" {
		http.Error(w, "Filename required", http.StatusBadRequest)
		return
	}

	filePath := filepath.Join(uploadPath, filename)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.NotFound(w, r)
		return
	}

	if err := os.Remove(filePath); err != nil {
		http.Error(w, "Failed to delete file", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func renameHandler(w http.ResponseWriter, r *http.Request) {
	filename := strings.TrimPrefix(r.URL.Path, "/rename/")
	if filename == "" {
		http.Error(w, "Filename required", http.StatusBadRequest)
		return
	}

	filePath := filepath.Join(uploadPath, filename)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.NotFound(w, r)
		return
	}

	if r.Method == http.MethodGet {
		// Show rename form
		tmpl := template.Must(template.New("rename").Parse(`
<!DOCTYPE html>
<html>
<head>
    <title>Rename Wallpaper</title>
</head>
<body>
    <h1>Rename Wallpaper</h1>
    <form action="/rename/{{.}}" method="post">
        <input type="text" name="newname" value="{{.}}" required>
        <button type="submit">Rename</button>
        <a href="/">Cancel</a>
    </form>
</body>
</html>
`))
		if err := tmpl.Execute(w, filename); err != nil {
			http.Error(w, "Failed to render template", http.StatusInternalServerError)
		}
		return
	}

	if r.Method == http.MethodPost {
		// Process rename
		newName := r.FormValue("newname")
		if newName == "" {
			http.Error(w, "New name required", http.StatusBadRequest)
			return
		}

		newPath := filepath.Join(uploadPath, newName)
		if err := os.Rename(filePath, newPath); err != nil {
			http.Error(w, "Failed to rename file", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

// Updated viewHandler template with download button
func viewHandler(w http.ResponseWriter, r *http.Request) {
	filename := strings.TrimPrefix(r.URL.Path, "/view/")
	if filename == "" {
		http.Error(w, "Filename required", http.StatusBadRequest)
		return
	}

	filePath := filepath.Join(uploadPath, filename)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.NotFound(w, r)
		return
	}

	tmpl := template.Must(template.New("view").Parse(`
<!DOCTYPE html>
<html>
<head>
    <title>View {{.Name}}</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; }
        img { max-width: 90%; max-height: 80vh; display: block; margin: 20px auto; }
        .actions { margin: 20px; text-align: center; }
    </style>
</head>
<body>
    <h1>{{.Name}}</h1>
    <img src="/images/{{.Name}}" alt="{{.Name}}">
    <div class="actions">
        <a href="/download/{{.Name}}">Download</a> | 
        <a href="/rename/{{.Name}}">Rename</a> | 
        <a href="/delete/{{.Name}}" onclick="return confirm('Are you sure?')">Delete</a> | 
        <a href="/">Back to list</a>
    </div>
</body>
</html>
`))

	if err := tmpl.Execute(w, struct{ Name string }{filename}); err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
	}
}