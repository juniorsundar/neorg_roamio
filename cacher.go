package main

import (
	"bufio"
	"database/sql"
	"os"
	"path"
	"strings"

	"github.com/juniorsundar/neorg_roamio/local"
	"github.com/juniorsundar/neorg_roamio/logger"
	_ "github.com/mattn/go-sqlite3"
)

type NorgMeta struct {
	Title       string   `json:"title"`
    Description string   `json:"description"`
    Authors     string   `json:"authors"`
	Categories  []string `json:"categories"`
	Created     string   `json:"created"`
	Updated     string   `json:"updated"`
	Version     string   `json:"version"`
}

type Unit struct {
	Id          string   `json:"id"`
	Title       string   `json:"title"`
	Address     string   `json:"address"`
	Line        int      `json:"line"`
	Type        string   `json:"type"`
	Categories  []string `json:"categories"`
	Created     string   `json:"created"`
	Description string   `json:"description"`
}

func extractNorgMetadata(norgFile string) (NorgMeta, error) {
	var norgMeta NorgMeta
	file, err := os.Open(norgFile)
	if err != nil {
		return norgMeta, err
	}
	defer file.Close()

	// Create a scanner to read the file line by line
	scanner := bufio.NewScanner(file)

	lineNumber := 0
	var titleLine string

	// Iterate through the file lines
	for scanner.Scan() {
		lineNumber++
		line := scanner.Text()

		// Check if we are on the second line
		if lineNumber == 2 {
			titleLine = line
			if titleLine == "" {
				return norgMeta, err
			}
			break
		}
	}

	if strings.HasPrefix(titleLine, "title: ") {
		norgMeta.Title = strings.TrimPrefix(titleLine, "title: ")
	}

	return norgMeta, nil
}

func cacheExists() bool {
	local.GetDatabase()

	// Check if the database file exists.
	_, err := os.Stat(local.DatabasePath)
	if os.IsNotExist(err) {
		logger.LogWarn.Printf("'%s' database doesn't exist. Creating it.", local.DatabasePath)

		// Create the database file.
		file, err := os.Create(local.DatabasePath)
		if err != nil {
			logger.LogErr.Printf("Error creating database file: %v", err)
			return false
		}
		file.Close()
		// Open the database to ensure it's a valid SQLite database
		db, err := sql.Open("sqlite3", local.DatabasePath)
		if err != nil {
			logger.LogErr.Printf("Error opening database: %v", err)
			return false
		}
		defer db.Close()

		return false
	} else if err != nil {
		// Handle potential errors from os.Stat
		logger.LogErr.Printf("Error checking database file: %v", err)
		return false
	} else {
		logger.LogInfo.Println("Found database " + local.DatabasePath)
	}

	// Open the database to ensure it's a valid SQLite database
	db, err := sql.Open("sqlite3", local.DatabasePath)
	if err != nil {
		logger.LogErr.Printf("Error opening database: %v", err)
		return false
	}
	defer db.Close()

	return true
}

func buildCache(relativeFileList []string) error {
	db, err := sql.Open("sqlite3", local.DatabasePath)
	if err != nil {
		logger.LogErr.Println("Error opening database!")
		return err
	}
	defer db.Close()

	createTableSQL := `
    CREATE TABLE IF NOT EXISTS nodes (
    id TEXT PRIMARY KEY,
    title TEXT,
    address TEXT,
    line INTEGER,
    type TEXT,
    categories TEXT, 
    created TEXT,
    description TEXT
    );
    `
	_, err = db.Exec(createTableSQL)
	if err != nil {
		logger.LogErr.Println("Error creating table!")
		return err
	}

	insertSQL := `
    INSERT INTO nodes (id, title, address, line, type, categories, created, description)
    VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
    `

	stmt, err := db.Prepare(insertSQL)
	if err != nil {
		logger.LogErr.Println("Error preparing SQL statement!")
		return err
	}
	defer stmt.Close()

	for _, relativeFile := range relativeFileList {
		fullPath := path.Join(local.ConfigData.Workspace.Root, relativeFile)

		norgMeta, err := extractNorgMetadata(fullPath)
		if err != nil {
			logger.LogWarn.Printf("Error extracting metadata from %s: %v", fullPath, err)
			continue
		}

		entry := Unit{
			Id:          relativeFile,
			Title:       norgMeta.Title,
			Address:     fullPath,
			Line:        1,
			Type:        "node",
			Categories:  norgMeta.Categories,
			Created:     norgMeta.Created,
			Description: norgMeta.Description,
		}

		_, err = stmt.Exec(
			entry.Id,
			entry.Title,
			entry.Address,
			entry.Line,
			entry.Type,
			strings.Join(entry.Categories, ","),
			entry.Created,
			entry.Description,
		)

		if err != nil {
			logger.LogWarn.Printf("Error inserting data into databas: %v", err)
			continue
		}
	}
	return nil
}

func invalidateCache(relativeFileList []string) error {
	return nil
}
