package main

// import "fmt"
import "flag"
import "os"
import "log"

func main() {
	dirPtr := flag.String("dir", ".", "Roam directory address.")
	portPtr := flag.String("port", "8080", "Roam server port.")
	verbosePtr := flag.Bool("verbose", false, "Enable verbose logging.")

    red := "\033[31m"
    // yellow := "\033[33m"
    reset := "\033[0m"

	flag.Parse()

	if *verbosePtr {
		log.Println("Directory:", *dirPtr)
		log.Println("Port:", *portPtr)
		log.Println("Verbose:", *verbosePtr)
	}

    info, err := os.Stat(*dirPtr)
    if os.IsNotExist(err) {
        log.Fatalf(red + "Directory %s does not exist." + reset, *dirPtr)
    } else if (os.IsPermission(err)) {
        log.Fatalf(red + "Directory %s does not have edit permission." + reset, *dirPtr)
    }
    if !info.IsDir() {
        log.Fatalf(red + "%s is not a directory." + reset, *dirPtr)
    }
}
