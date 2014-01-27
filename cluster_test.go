package cluster

import (
	"testing"
	"fmt"
	"sync"
	"time"
	"strconv"
)

func TestCluster(t *testing.T) {
        wg := new(sync.WaitGroup)
        fmt.Println("ok");
        //go startServer(5,wg)
        for i:=1;i<4;i++{
        wg.Add(1)
        wg.Add(1)
        //wg.Add(1)
        //wg.Add(1)
        go startServer(i,wg)
        //go startServer(2,wg)
}
        wg.Wait()
}

func startServer(id int,wg *sync.WaitGroup) {

        server := New(id, "./config.json")
        //server.Outbox() <- &Envelope{Pid:BROADCAST, Msg: "hello there"}
        ch := make(chan []byte) //stdin channel variable
        go generateMsg(ch, server,wg)
        for {
                select {
                case envelope := <-server.Inbox():
                        fmt.Printf("Received msg at %d from %d: '%s'\n",server.Pid(), envelope.Sen_Pid, envelope.Msg)

                case msg:=<-ch:
                        //fmt.Printf("Received msg from %d: '%s'\n", envelope.Pid, envelope.Msg)
                        //temp:=2
                        //fmt.Println(Envelope{Pid:temp, Msg: string(msg)});
                        server.Outbox() <- &Envelope{Pid: BROADCAST, Msg: string(msg)}
                        //server.Outbox() <- &Envelope{Pid:2, Msg: string(msg)}
                case <-time.After(10 * time.Second):
                        wg.Done()
                        println("Waited and waited. Ab thak gaya\n")

              }
      }
}


//reference http://stackoverflow.com/questions/9680812/how-can-i-make-net-read-wait-for-input-in-golang

// Start a goroutine to read from our net connection
func generateMsg(ch chan []byte, ser ServerType,wg *sync.WaitGroup) {

        for i:=0;i<1001;i++{
                time.After(100)
                msg:= strconv.Itoa(ser.MyPid)+"."+strconv.Itoa(i)
                ch <- []byte(msg)
        }
        wg.Done()
}

