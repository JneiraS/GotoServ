package utils

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	cts "github.com/JneiraS/GotoServ/internal/constants"
)

// UpdateOrAddCSVRecord met à jour une ligne existante (par agent) ou ajoute une nouvelle ligne dans le CSV.
func UpdateOrAddCSVRecord(AssignmentsCSV, agent, scope, keywords string) error {
	file, err := os.OpenFile(AssignmentsCSV, os.O_RDWR, 0644)
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

	// Upsert par (scope, keywords), pas par agent
	found := false
	targetScope := strings.TrimSpace(scope)
	targetKeywords := strings.TrimSpace(keywords)

	for i, row := range records {
		if i == 0 {
			continue // header
		}
		if len(row) < 3 {
			continue
		}

		rowScope := strings.TrimSpace(row[1])
		rowKeywords := strings.TrimSpace(row[2])

		if strings.EqualFold(rowScope, targetScope) && rowKeywords == targetKeywords {
			records[i][0] = agent    // agent
			records[i][1] = scope    // scope
			records[i][2] = keywords // keywords
			found = true
			break
		}
	}

	if !found {
		records = append(records, []string{agent, scope, keywords})
	}

	if err := file.Truncate(0); err != nil {
		return err
	}
	if _, err := file.Seek(0, 0); err != nil {
		return err
	}

	writer := csv.NewWriter(file)
	writer.Comma = ';'
	if err := writer.WriteAll(records); err != nil {
		return err
	}
	writer.Flush()
	return writer.Error()
}

// UpdateKeywordsForAgent met à jour uniquement les keywords d'une ligne identifiée par agent.
// Retourne false si l'agent n'existe pas dans le CSV.
func UpdateKeywordsForAgent(assignmentsCSV, agent, keywords string) (bool, error) {
	file, err := os.OpenFile(assignmentsCSV, os.O_RDWR, 0644)
	if err != nil {
		return false, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ';'
	records, err := reader.ReadAll()
	if err != nil {
		return false, err
	}

	found := false
	target := strings.TrimSpace(agent)
	for i, row := range records {
		if i == 0 {
			continue // header
		}
		if len(row) < 3 {
			continue
		}
		if strings.EqualFold(strings.TrimSpace(row[0]), target) {
			records[i][2] = keywords
			found = true
			break
		}
	}

	if !found {
		return false, nil
	}

	if err := file.Truncate(0); err != nil {
		return false, err
	}
	if _, err := file.Seek(0, 0); err != nil {
		return false, err
	}

	writer := csv.NewWriter(file)
	writer.Comma = ';'
	if err := writer.WriteAll(records); err != nil {
		return false, err
	}
	writer.Flush()
	return true, writer.Error()
}

// ConvertCSVToJSONGeneric convertit n'importe quel CSV en JSON (tableau d'objets).
// Les clés JSON sont les noms de colonnes de l'en-tête CSV.
func ConvertCSVToJSONGeneric(AssignmentsCSV, AssignmentsJSON string, delimiter rune) error {
	rows, err := ReadCSVAsMaps(AssignmentsCSV, delimiter)
	if err != nil {
		return err
	}

	out, err := json.MarshalIndent(rows, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal json: %w", err)
	}

	if err := os.WriteFile(AssignmentsJSON, out, 0o644); err != nil {
		return fmt.Errorf("write json: %w", err)
	}

	return nil
}

// ReadCSVAsMaps lit un CSV et retourne []map[colonne]valeur.
func ReadCSVAsMaps(AssignmentsCSV string, delimiter rune) ([]map[string]string, error) {
	if delimiter == 0 {
		delimiter = ';'
	}

	f, err := os.Open(AssignmentsCSV)
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
// It reads from AssignmentsCSV and writes the output to AssignmentsJSON.
// If the conversion fails, the program exits with a fatal error.
func CreatJsonFromCsv() {
	if err := ConvertCSVToJSONGeneric(cts.AssignmentsCSV, cts.AssignmentsJSON, ';'); err != nil {
		log.Fatal(err)
	}
}
