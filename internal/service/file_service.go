package service

import (
	"fmt"
	"os"
	"path/filepath"
)

// FileService fornece operações de sistema de arquivos.
type FileService interface {
	// EnsureDir garante que o diretório existe.
	EnsureDir(path string) error
	// CreateFile cria um arquivo e retorna o writer.
	CreateFile(path string) (*os.File, error)
	// RemoveFile remove um arquivo.
	RemoveFile(path string) error
	// RenameFile renomeia um arquivo.
	RenameFile(oldPath, newPath string) error
	// GetDir retorna o diretório de um caminho.
	GetDir(path string) string
}

// fileService implementa FileService.
type fileService struct{}

// NewFileService cria uma nova instância de FileService.
func NewFileService() FileService {
	return &fileService{}
}

// EnsureDir garante que o diretório existe, criando-o se necessário.
func (s *fileService) EnsureDir(path string) error {
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return fmt.Errorf("erro ao criar diretório %s: %w", path, err)
	}
	return nil
}

// CreateFile cria um novo arquivo.
func (s *fileService) CreateFile(path string) (*os.File, error) {
	file, err := os.Create(path)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar arquivo %s: %w", path, err)
	}
	return file, nil
}

// RemoveFile remove um arquivo.
func (s *fileService) RemoveFile(path string) error {
	return os.Remove(path)
}

// RenameFile renomeia um arquivo.
func (s *fileService) RenameFile(oldPath, newPath string) error {
	if err := os.Rename(oldPath, newPath); err != nil {
		return fmt.Errorf("erro ao renomear arquivo: %w", err)
	}
	return nil
}

// GetDir retorna o diretório de um caminho.
func (s *fileService) GetDir(path string) string {
	return filepath.Dir(path)
}
