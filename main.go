package main

import(
    "flag"
)

/*------------------------*
 * COMMAND LINE INTERFACE *
 *------------------------*/
func main() {
    // set up flags
    var listenPort string
    flag.StringVar(&listenPort, "l", "", "")
    flag.StringVar(&listenPort, "listen", "", "")

    var seedData string
    flag.StringVar(&seedData, "s", "", "")
    flag.StringVar(&seedData, "seed", "", "")

    var helpFlag bool
    flag.BoolVar(&helpFlag, "h", false, "")
    flag.BoolVar(&helpFlag, "help", false, "")

    var publicFlag bool
    flag.BoolVar(&publicFlag, "p", false, "")
    flag.BoolVar(&publicFlag, "public", false, "")

    flag.Parse()

    listenPort   = ":" + listenPort

    if helpFlag {
        showGlobalHelp()
        return
    }

    myNode := newNode()
    myNode.run(listenPort, seedData, publicFlag)
}
