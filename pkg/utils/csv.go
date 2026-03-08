package utils

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

const (
	csvPath  = "assignement_fcb.csv"
	jsonPath = "assignement_fcb.json"
)

// UpdateOrAddCSVRecord met à jour une ligne existante (par agent) ou ajoute une nouvelle ligne dans le CSV.
func UpdateOrAddCSVRecord(csvPath, agent, scope, keywords string) error {
	// Lire tout le CSV
	file, err := os.OpenFile(csvPath, os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ';'
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	// Chercher si l'agent existe déjà
	found := false
	for i, row := range records {
		if i == 0 {
			continue // skip header
		}
		if len(row) > 0 && row[0] == agent {
			if len(row) > 1 {
				records[i][1] = scope
			}
			if len(row) > 2 {
				records[i][2] = keywords
			}
			found = true
			break
		}
	}

	// Ajouter une nouvelle ligne si non trouvé
	if !found {
		records = append(records, []string{agent, scope, keywords})
	}

	// Réécrire tout le CSV
	if err := file.Truncate(0); err != nil {
		return err
	}
	if _, err := file.Seek(0, 0); err != nil {
		return err
	}
	writer := csv.NewWriter(file)
	writer.Comma = ';'
	err = writer.WriteAll(records)
	if err != nil {
		return err
	}
	writer.Flush()
	return writer.Error()
}

// ConvertCSVToJSONGeneric convertit n'importe quel CSV en JSON (tableau d'objets).
// Les clés JSON sont les noms de colonnes de l'en-tête CSV.
func ConvertCSVToJSONGeneric(csvPath, jsonPath string, delimiter rune) error {
	rows, err := ReadCSVAsMaps(csvPath, delimiter)
	if err != nil {
		return err
	}

	out, err := json.MarshalIndent(rows, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal json: %w", err)
	}

	if err := os.WriteFile(jsonPath, out, 0o644); err != nil {
		return fmt.Errorf("write json: %w", err)
	}

	return nil
}

// ReadCSVAsMaps lit un CSV et retourne []map[colonne]valeur.
func ReadCSVAsMaps(csvPath string, delimiter rune) ([]map[string]string, error) {
	if delimiter == 0 {
		delimiter = ';'
	}

	f, err := os.Open(csvPath)
	if err != nil {
		return nil, fmt.Errorf("open csv: %w", err)
	}
	defer f.Close()

	r := csv.NewReader(f)
	r.Comma = delimiter
	r.FieldsPerRecord = -1

	header, err := r.Read()
	if err != nil {
		return nil, fmt.Errorf("read header: %w", err)
	}
	for i := range header {
		header[i] = strings.TrimSpace(strings.TrimPrefix(header[i], "\uFEFF"))
	}

	var rows []map[string]string

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("read record: %w", err)
		}

		row := make(map[string]string, len(header))
		for i, col := range header {
			if i < len(record) {
				row[col] = strings.TrimSpace(record[i])
			} else {
				row[col] = ""
			}
		}

		// Colonnes en trop: extra_1, extra_2, ...
		for i := len(header); i < len(record); i++ {
			key := fmt.Sprintf("extra_%d", i-len(header)+1)
			row[key] = strings.TrimSpace(record[i])
		}

		rows = append(rows, row)
	}

	return rows, nil
}

// CreatJsonFromCsv converts a CSV file to a JSON file using a semicolon as the delimiter.
// It reads from csvPath and writes the output to jsonPath.
// If the conversion fails, the program exits with a fatal error.
func CreatJsonFromCsv() {
	if err := ConvertCSVToJSONGeneric(csvPath, jsonPath, ';'); err != nil {
		log.Fatal(err)
	}
}
