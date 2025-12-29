package view

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/alan-gomes1/jira-reporter/internal/model"
)

// docxGenerator implementa ReportGenerator para formato DOCX.
type docxGenerator struct {
	fileService fileRenamer
}

// fileRenamer interface mínima para operações de arquivo.
type fileRenamer interface {
	RenameFile(oldPath, newPath string) error
}

// NewDOCXGenerator cria um novo gerador DOCX.
func NewDOCXGenerator(fileService fileRenamer) ReportGenerator {
	return &docxGenerator{
		fileService: fileService,
	}
}

// Generate converte HTML para DOCX usando LibreOffice.
// args[0]: caminho do arquivo HTML de entrada
// args[1]: caminho do arquivo DOCX de saída
func (g *docxGenerator) Generate(
	writer io.Writer, data *model.ReportData, args ...string,
) error {
	if len(args) < 2 {
		return fmt.Errorf(
			"argumentos insuficientes: necessário htmlPath e docxPath",
		)
	}

	htmlPath := args[0]
	docxPath := args[1]

	loPath, err := g.findLibreOffice()
	if err != nil {
		return err
	}

	return g.convert(loPath, htmlPath, docxPath)
}

// Format retorna o formato suportado.
func (g *docxGenerator) Format() model.ReportFormat {
	return model.FormatDOCX
}

// findLibreOffice localiza o executável do LibreOffice.
func (g *docxGenerator) findLibreOffice() (string, error) {
	loPath, err := exec.LookPath("libreoffice")
	if err == nil {
		return loPath, nil
	}

	// Tenta encontrar o soffice (nome alternativo)
	loPath, err = exec.LookPath("soffice")
	if err == nil {
		return loPath, nil
	}

	return "", fmt.Errorf(
		"LibreOffice não está instalado. Exemplo de instalação: " +
			"sudo apt-get install libreoffice-writer",
	)
}

// convert executa a conversão de HTML para DOCX.
func (g *docxGenerator) convert(loPath, htmlPath, docxPath string) error {
	outputDir := filepath.Dir(docxPath)

	cmd := exec.Command(
		loPath,
		"--headless",
		"--convert-to",
		"docx:MS Word 2007 XML",
		"--outdir",
		outputDir,
		htmlPath,
	)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf(
			"erro ao executar LibreOffice: %w - %s", err, stderr.String(),
		)
	}

	// LibreOffice gera o arquivo com o mesmo nome do HTML mas com extensão .docx
	generatedDocx := strings.TrimSuffix(htmlPath, ".html") + ".docx"
	if generatedDocx != docxPath {
		if err := g.fileService.RenameFile(generatedDocx, docxPath); err != nil {
			// Fallback para renomeação direta
			if err := os.Rename(generatedDocx, docxPath); err != nil {
				return fmt.Errorf("erro ao renomear arquivo DOCX: %w", err)
			}
		}
	}

	return nil
}
