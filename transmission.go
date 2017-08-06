package main 

/*
A Transmission is simply a message that will be sent throughout the network.  
It includes the actual message as well as the addresses of nodes who have
already received the Transmission
*/

type Transmission struct {
    Block Block
    BeenSent bool
    Sender string
}

func (t *Transmission) updateBeenSent() {
    t.BeenSent = true
}

func (t *Transmission) updateSender(address string){
    t.Sender = address
}