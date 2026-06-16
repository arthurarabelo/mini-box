package files

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Diretório raiz do storage — em produção viria de configuração
const storageRoot = "./storage"

// FileItem representa um arquivo ou pasta na listagem
type FileItem struct {
	Name     string    `json:"name"`
	Type     string    `json:"type"` // "file" ou "folder"
	Size     int64     `json:"size,omitempty"`
	Modified time.Time `json:"modified"`
	MimeType string    `json:"mimeType,omitempty"`
}

// sanitizePath garante que o path não saia do storageRoot (path traversal)
// Segurança: impede chamadas como ?path=../../etc/passwd
func sanitizePath(rawPath string) (string, error) {
	// Remove barras iniciais e pontos duplos
	clean := filepath.Clean("/" + rawPath)
	full := filepath.Join(storageRoot, clean)

	// Confirma que o path está dentro do storageRoot
	if !strings.HasPrefix(full, filepath.Clean(storageRoot)) {
		return "", os.ErrPermission
	}
	return full, nil
}

// ListHandler processa GET /api/files?path=/pasta
func ListHandler(w http.ResponseWriter, r *http.Request) {
	rawPath := r.URL.Query().Get("path")
	dirPath, err := sanitizePath(rawPath)
	if err != nil {
		http.Error(w, "Path inválido", http.StatusBadRequest)
		return
	}

	entries, err := os.ReadDir(dirPath)
	if err != nil {
		if os.IsNotExist(err) {
			http.Error(w, "Pasta não encontrada", http.StatusNotFound)
			return
		}
		http.Error(w, "Erro ao ler diretório", http.StatusInternalServerError)
		return
	}

	items := make([]FileItem, 0, len(entries))
	for _, entry := range entries {
		info, _ := entry.Info()
		item := FileItem{
			Name:     entry.Name(),
			Modified: info.ModTime(),
		}
		if entry.IsDir() {
			item.Type = "folder"
		} else {
			item.Type = "file"
			item.Size = info.Size()
			item.MimeType = detectMime(entry.Name())
		}
		items = append(items, item)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

// UploadHandler processa POST /api/files/upload?path=/pasta
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	rawPath := r.URL.Query().Get("path")
	dirPath, err := sanitizePath(rawPath)
	if err != nil {
		http.Error(w, "Path inválido", http.StatusBadRequest)
		return
	}

	// Cria o diretório se não existir
	os.MkdirAll(dirPath, 0755)

	// Lê o arquivo do multipart form (limite de 50MB)
	r.ParseMultipartForm(50 << 20)
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Arquivo não encontrado no form", http.StatusBadRequest)
		return
	}
	defer file.Close()

	destPath := filepath.Join(dirPath, filepath.Base(header.Filename))
	dest, err := os.Create(destPath)
	if err != nil {
		http.Error(w, "Erro ao salvar arquivo", http.StatusInternalServerError)
		return
	}
	defer dest.Close()

	io.Copy(dest, file)
	w.WriteHeader(http.StatusCreated)
}

// DownloadHandler processa GET /api/files/download?path=/pasta/arquivo.txt
func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	rawPath := r.URL.Query().Get("path")
	filePath, err := sanitizePath(rawPath)
	if err != nil {
		http.Error(w, "Path inválido", http.StatusBadRequest)
		return
	}

	http.ServeFile(w, r, filePath)
}

// detectMime retorna um MIME type básico pela extensão
func detectMime(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	mimes := map[string]string{
		".md":   "text/markdown",
		".txt":  "text/plain",
		".csv":  "text/csv",
		".pdf":  "application/pdf",
		".jpg":  "image/jpeg",
		".jpeg": "image/jpeg",
		".png":  "image/png",
		".docx": "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
	}
	if m, ok := mimes[ext]; ok {
		return m
	}
	return "application/octet-stream"
}
