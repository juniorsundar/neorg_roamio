package main

import (
	"bufio"
	"encoding/json"
	"os"
	// "path/filepath"
	"strings"

	"github.com/juniorsundar/neorg_roamio/logger"
	"github.com/juniorsundar/neorg_roamio/local"
	_ "github.com/mattn/go-sqlite3"
)

type Node struct {
	Address string `json:"address"`
	Line    int    `json:"line"`
	Title   string `json:"title"`
	Type    string `json:"type"`
}

type Cache struct {
	Nodes map[string]Node `json:"nodes"`
}

type NorgMeta struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Authors     string   `json:"authors"`
	Categories  []string `json:"categories"`
	Created     string   `json:"created"`
	Updated     string   `json:"updated"`
	Version     string   `json:"version"`
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
	_, err := os.Stat(local.ConfigData.Workspace.Root + "/.roamioCache.json")
	if os.IsNotExist(err) {
		logger.LogWarn.Printf("%s file doesn't exist.", local.ConfigData.Workspace.Root+"/.roamioCache.json")
		err := os.WriteFile(local.ConfigData.Workspace.Root+"/.roamioCache.json", []byte(""), 0666)
		if err != nil {
			logger.LogErr.Println(err)
		}
		return false
	} else {
		logger.LogInfo.Println("Found file " + local.ConfigData.Workspace.Root + "/.roamioCache.json")
		return true
	}
}

func buildCache() error {
	return nil
}

func invalidateCache(relativeFileList []string) error {
	cacheFile, err := os.Open(local.ConfigData.Workspace.Root + "/.roamioCache.json")
	if err != nil {
		return err
	}
	defer cacheFile.Close()

	// To hold the full JSON structure
	var result Cache

	decoder := json.NewDecoder(cacheFile)
	err = decoder.Decode(&result)
	if err != nil {
		return err
	}

	// invalidate file nodes only for the moment
	// TODO also look at non-files like blocks
	fileNodeNames := make(map[string]bool)
	for _, fileNode := range relativeFileList {
		fileNodeNames[fileNode] = true
	}

	// Iterate over the nodes in the cache and check against the fileNodeNames map
	for address, node := range result.Nodes {
		if node.Type == "file" {
			if _, exists := fileNodeNames[address]; exists {
				fileNodeNames[address] = false // Mark as found
			}
		}
	}

	for fileNode, notFound := range fileNodeNames {
		if notFound {
			logger.LogErr.Printf("Missing %s", fileNode)
            norgMeta, err := extractNorgMetadata(local.ConfigData.Workspace.Root+"/"+fileNode)
            if err != nil {
                return err
            }

            newNode := Node{
                Address: fileNode,
                Line: 1,
                Title: norgMeta.Title,
                Type: "file",
            }
            result.Nodes[fileNode] = newNode

		} else {
			logger.LogWarn.Printf("Found %s", fileNode)
		}
	}

    // Marshal the cache struct into JSON
    jsonData, err := json.MarshalIndent(result, "", "    ") // Pretty print with indentation
    if err != nil {
        return err
    }

    // Write the JSON data to a file
    err = os.WriteFile(local.ConfigData.Workspace.Root + "/.roamioCache.json", jsonData, 0644)
    if err != nil {
        return err
    }

	return nil
}
