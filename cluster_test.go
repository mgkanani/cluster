package cluster

import (
	"fmt"
	"strconv"
	"sync"
	"testing"
	"time"
)

var counter map[int]int //counts total number of messages received for each PIDs.
var msgs map[int]map[string]string // stores msg received for each PIDs, len(msgs[i]) will give total unique msgs received at i.

var total_servers int
var total_msgs int

func TestCluster(t *testing.T) {
	total_servers=4;
	total_msgs=10001;

	counter = make(map[int]int)
	msgs = make(map[int]map[string]string)
	wg := new(sync.WaitGroup)
	//fmt.Println(len(t_str),"ok")
	//go startServer(5,wg)
	for i := 1; i < total_servers+1; i++ {
		msgs[i] = make(map[string]string)
		counter[i] = 0
		wg.Add(1)//increment wg's counter by 1
		wg.Add(1)//increment wg's counter by 1
		go startServer(i, wg)//server will be started in new thread.
	}

	wg.Wait() //This will wait till wg's count becomes zero.
	//this will prints the summary of message received for each PIDs.
	for i := 1; i < total_servers+1; i++ {
		fmt.Println("total msgs received at PID -", i, ":-", counter[i],"\t total unique msgs:-", len(msgs[i]))
	}
}

func startServer(id int, wg *sync.WaitGroup) {
	//Initialises server'data.
	server := New(id, "./config.json")

	ch := make(chan []byte) //channel variable used for retrieving msgs generated by child thread.

	//to test for data transfer of >60000 bytes uncomment below line and comment next-line of below line.
	//go generateLongMsg(ch, server, wg)
	go generateMsg(ch, server, wg)
	for {
		select {
		case envelope := <-server.Inbox():
			//case <-server.Inbox():
			str := envelope.Msg.(string)

			//to print length of msgs uncomment below line.
			//fmt.Printf("Received msg at PID -%d from %d: '%d'\n",server.Pid(), envelope.Pid, len(envelope.Msg.(string)))

			// to print msgs recieved uncomment below line.
			//fmt.Printf("Received msg at PID -%d from %d: '%s'\n",server.Pid(), envelope.Pid, envelope.Msg)

			counter[id]++
			msgs[id][str] = str

		case msg := <-ch:
			//fmt.Println(Envelope{Pid:temp, Msg: string(msg)});
			server.Outbox() <- &Envelope{Pid: BROADCAST, Msg: string(msg)}

		case <-time.After(2 * time.Second):
			//println("Waited and waited. Ab thak gaya\n")
			wg.Done()//waits for 2 seconds,if there is no any events then decrements wg's counter by 1 and exit from this procedure.
			return
		}
	}
}

//reference http://stackoverflow.com/questions/9680812/how-can-i-make-net-read-wait-for-input-in-golang

func generateMsg(ch chan []byte, ser ServerType, wg *sync.WaitGroup) {
	//this will generates total 10001 msgs.
	for i := 0; i < total_msgs; i++ {
		time.After(1)
		msg := strconv.Itoa(ser.MyPid) + "." + strconv.Itoa(i)
		ch <- []byte(msg)
	}
	wg.Done() //decrements wg's counter by 1
}


func generateLongMsg(ch chan []byte, ser ServerType, wg *sync.WaitGroup) {
        //code to make string of length 65536. Uncomment to test for it.
        t_str:="h"
        for i:=0;i<16;i++{
        t_str=t_str+t_str;
        }
        ch <- []byte(t_str)
        wg.Done()
}

