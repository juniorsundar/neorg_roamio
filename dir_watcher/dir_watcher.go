package dir_watcher

import(
    "os"
    "path/filepath"

	"github.com/fsnotify/fsnotify"
    "github.com/juniorsundar/neorg_roamio/logger"
)

// Initialise watcher for subdirectories recursively.
//
// @param watcher: Main `watcher` from main loop.
//
// @param dir: Root directory as source of recursive search.
func AddWatchDirRecursively(watcher *fsnotify.Watcher, dir string) error {
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

// Routine function to observe the assigned directory for changes
//
// @param watcher: Directory watcher in play.
//
// @param verbose: Verbosity flag.
func DirWatcher(watcher *fsnotify.Watcher, verbose bool) {
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

			if verbose {
				logger.LogInfo.Println("Event:", event.Name)
			}

			if event.Op&fsnotify.Create == fsnotify.Create {
				if err == nil && file.IsDir() {
					err = AddWatchDirRecursively(watcher, event.Name)
					if err != nil {
						logger.LogErr.Fatalln(err)
					}
				}
			}

			if event.Op&fsnotify.Remove == fsnotify.Remove {
				if os.IsNotExist(err) {
					watcher.Remove(event.Name)
					logger.LogWarn.Printf("Removed %s from watcher.", event.Name)
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

