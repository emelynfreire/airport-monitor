package main

import "sync"

func main() {
    var wg sync.WaitGroup
    wg.Add(2)

    go func() {
        consume("arrival-flights", "FlightTotemAll")
        wg.Done()
    }()
    go func() {
        consume("departure-flights", "FlightTotemAll")
        wg.Done()
    }()

    wg.Wait()
}

