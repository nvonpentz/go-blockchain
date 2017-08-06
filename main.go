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
    flag.StringVar(&listenPort, "l", "1999", "")
    flag.StringVar(&listenPort, "listen", "1999", "")

    var seedPort string
    flag.StringVar(&seedPort, "s", "", "")
    flag.StringVar(&seedPort, "seed", "", "")

    var helpFlag bool
    flag.BoolVar(&helpFlag, "h", false, "")
    flag.BoolVar(&helpFlag, "help", false, "")

    var joinFlag bool
    flag.BoolVar(&joinFlag, "j", false, "")
    flag.BoolVar(&joinFlag, "join", false, "")

    flag.Parse()

    listenPort = ":" + listenPort 
    seedPort = ":" + seedPort

    if helpFlag {
        showGlobalHelp()
        return
    }

    myNode := newNode()
    myNode.run(listenPort, seedPort,joinFlag)
}
