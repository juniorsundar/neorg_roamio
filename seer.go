package main

import (
    "os"
    "path/filepath"
    "strings"
    "sync"

    "github.com/fsnotify/fsnotify"
    "github.com/juniorsundar/neorg_roamio/logger"
    "github.com/juniorsundar/neorg_roamio/local"
)

var (
    mu            sync.Mutex
    fileList      []string
)

// Initialise watcher for subdirectories recursively.
//
// Parameters:
//   - watcher: Directory watcher in play.
//   - dir: Root directory as source of recursive search.
func addWatchDirRecursively(watcher *fsnotify.Watcher, dir string) error {
    return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }

        if info.IsDir() {
            if info.Name()[0] == '.' {
                logger.LogWarn.Println("Skipping hidden directory:", path)
                return filepath.SkipDir
            }

            err = watcher.Add(path)
            if err != nil {
                return err
            }
            logger.LogWarn.Printf("Watching: %s", path)
        }
        return nil
    })
}

// Routine function to observe the assigned directory for changes.
//
// Parameters:
//   - watcher: Directory watcher in play.
func dirWatcher(watcher *fsnotify.Watcher) {
    for {
        // Select is like Switch in that it observes what comes out of the
        // channels. In this case, the two channels of interest are
        // watcher.Events/Write. If any one of them gets a message then their
        // respective case is triggered.
        select {
        case event, ok := <-watcher.Events:
            if !ok {
                return
            }
            file, err := os.Stat(event.Name)

            if event.Op&fsnotify.Create == fsnotify.Create {
                logger.LogInfo.Printf("Created %s", event.Name)
                if err == nil && file.IsDir() {
                    err = addWatchDirRecursively(watcher, event.Name)
                    if err != nil {
                        logger.LogErr.Fatalln(err)
                    }
                }
            }

            if event.Op&fsnotify.Remove == fsnotify.Remove {
                if os.IsNotExist(err) {
                    watcher.Remove(event.Name)
                    logger.LogWarn.Printf("Removed %s", event.Name)
                }
            }

            if event.Op&fsnotify.Rename == fsnotify.Rename {
                if os.IsNotExist(err) {
                    watcher.Remove(event.Name)
                    logger.LogWarn.Printf("Renamed %s", event.Name)
                }
            }
        case err, ok := <-watcher.Errors:
            if !ok {
                return
            }
            logger.LogErr.Println(err)
        }
    }
}

// Depth traverse watched directories for files
//
// Parameters:
//   - watcher: Directory watcher in play
//   - extension: Filetype we are interested in listing
func listFilesRecursively(watcher *fsnotify.Watcher, extension string) []string {
    watchedFolders := watcher.WatchList()

    var wg sync.WaitGroup
    wg.Add(len(watchedFolders))

    for _, folder := range watchedFolders {
        go func(folder string, extension string, wg *sync.WaitGroup) {
            defer wg.Done()
            dirPath := strings.Split(folder, "/")
            wsRelativeFolder := strings.Join(dirPath[len(strings.Split(local.ConfigData.Workspace.Root, "/")):], "/")

            dirList, err := os.ReadDir(folder)
            if err != nil {
                logger.LogErr.Fatalf("Directory doesn't exist.")
            }

            var fileStringList []string
            for _, file := range dirList {
                stringSplit := strings.Split(file.Name(), ".")
                if len(stringSplit) > 1 && stringSplit[1] == "norg" {
                    fullPath := filepath.Join(wsRelativeFolder, file.Name())
                    fileStringList = append(fileStringList, fullPath)
                }
            }

            mu.Lock()
            fileList = append(fileList, fileStringList...)
            mu.Unlock()

        }(folder, extension, &wg)
    }

    wg.Wait()
    return fileList
}
