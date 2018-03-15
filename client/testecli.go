package main

import (
	"encoding/gob"
	"fmt"
	"math/rand"
	"net"

	"github.com/Seriyin/lab0-go/bank"
)

func main() {
	n := rand.Intn(12) + 4
	ch := make([]chan int64, n, n)
	for i := range ch {
		ch[i] = make(chan int64, 1)
	}
	ip := net.IPv4(127, 0, 0, 1)
	tcp := new(net.TCPAddr)
	tcp.IP = ip
	tcp.Port = 22556
	master, err := net.DialTCP("tcp", nil, tcp)
	if err != nil {
		panic(err)
	}
	defer master.Close()
	dec := gob.NewDecoder(master)
	enc := gob.NewEncoder(master)
	for i := 0; i < n; i++ {
		f, err := net.DialTCP("tcp", nil, tcp)
		if err != nil {
			// handle error
			panic(err)
		}
		go spamOps(f, ch[i])
	}
	r := int64(0)
	for _, c := range ch {
		for i := range c {
			r += i
		}
	}
	rep := new(bank.Reply)
	enc.Encode(bank.Message{Op: 0, Mov: 0})
	dec.Decode(rep)
	fmt.Printf("Got %d, Expected %d\n", r, rep.Balance)
}

func spamOps(conn net.Conn, ch chan int64) {
	defer conn.Close()
	dec := gob.NewDecoder(conn)
	enc := gob.NewEncoder(conn)
	r := new(bank.Reply)
	sum := int64(0)
	for i := rand.Intn(30000) + 50000; i > 0; i-- {
		enc.Encode(bank.Message{Op: 1, Mov: rand.Int63n(400) - 200})
		dec.Decode(r)
		if r.Res {
			sum += r.Balance
		} else {
			fmt.Printf("Rejected %d\n", r.Balance)
		}
	}
	ch <- sum
	close(ch)
}
