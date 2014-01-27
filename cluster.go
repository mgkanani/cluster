package Cluster

import (
	"encoding/json"
	"fmt"
	zmq "github.com/pebbe/zmq4"
	"io/ioutil"
	"os"
	"strconv"
)

const (
	BROADCAST = -1
)

type Envelope struct {
	// On the sender side, Pid identifies the receiving peer. If instead, Pid is
	// set to cluster.BROADCAST, the message is sent to all peers. On the receiver side, the
	// Id is always set to the original sender. If the Id is not found, the message is silently dropped
	Pid int

	Sen_Pid int

	// An id that globally and uniquely identifies the message, meant for duplicate detection at
	// higher levels. It is opaque to this package.
	MsgId int64

	// the actual message.
	Msg interface{}
}

type Server interface {

	// Id of this server
	Pid() int

	// array of other servers' ids in the same cluster
	Peers() []int

	// the channel to use to send messages to other peers
	// Note that there are no guarantees of message delivery, and messages
	// are silently dropped
	Outbox() chan *Envelope

	// the channel to receive messages from other peers.
	Inbox() chan *Envelope
}

type jsonobject struct {
	Object ServersType
}

type ServersType struct {
	Servers []ServerType
}
type ServerType struct {
	Socket    string
	MyPid     int
	PeerIds   []int
	Map       map[string]string
	SocketMap map[string]*zmq.Socket
	in        chan *Envelope
	out       chan *Envelope
}

func New(pid int, filenm string) ServerType {
	file, e := ioutil.ReadFile(filenm)
	if e != nil {
		fmt.Printf("File error: %v\n", e)
		os.Exit(1)
	}
	var servers jsonobject
	err := json.Unmarshal(file, &servers)
	if err != nil {
		fmt.Println("error:", err)
	}

	ret_ser := ServerType{Socket: "", MyPid: pid}
	peers := make([]int, len(servers.Object.Servers)-1)
	Mp := make(map[string]string)
	SocketMp := make(map[string]*zmq.Socket)
	for i, j := 0, 0; i < len(servers.Object.Servers); i++ {
		serv := servers.Object.Servers[i]
		Mp[strconv.Itoa(serv.MyPid)] = serv.Socket
		if serv.MyPid == pid {
			//fmt.Printf("%v", serv, Mp)
			ret_ser.Socket = serv.Socket
			ret_ser.MyPid = serv.MyPid
		} else {
			peers[j] = serv.MyPid
			SocketMp[strconv.Itoa(serv.MyPid)], _ = zmq.NewSocket(zmq.DEALER)
			SocketMp[strconv.Itoa(serv.MyPid)].Connect("tcp://" + Mp[strconv.Itoa(serv.MyPid)])
			j++
		}
	}
	ret_ser.PeerIds = peers
	ret_ser.Map = Mp
	ret_ser.SocketMap = SocketMp
	//fmt.Println(peers)

	ret_ser.in, ret_ser.out = make(chan *Envelope), make(chan *Envelope)
	go ret_ser.SendMsg(ret_ser.out)
	go ret_ser.GetMsg(ret_ser.in)

	return ret_ser
}

func (ser ServerType) SendMsg(enve chan *Envelope) {
	for {
		x := <-enve
		x.Sen_Pid = ser.MyPid
		//fmt.Println(ser.Map, ser.PeerIds)
		//fmt.Println("In OutBox(SendMsg) Message:-", *x)
		if x.Pid == -1 {
			//BroadCast
			for i := 0; i < len(ser.PeerIds); i++ {
				str := strconv.Itoa(ser.PeerIds[i])
				//fmt.Println("Map[", str, "]=", ser.Map[strconv.Itoa(ser.PeerIds[i])])
				/*				dealer, _ := zmq.NewSocket(zmq.DEALER)
								dealer.Connect("tcp://" + ser.Map[strconv.Itoa(ser.PeerIds[i])])
				*/
				data, _ := json.Marshal(*x)
				//				dealer.Send(string(data), 0)
				ser.SocketMap[str].Send(string(data), 0)
			}
		} else {
			//Unicast
			_, ok := ser.Map[strconv.Itoa(x.Pid)]
			if ok && x.Pid!=ser.MyPid{
				//fmt.Println("Send Unicast\n")
				str_pid := strconv.Itoa(x.Pid)
				//dealer, _ := zmq.NewSocket(zmq.DEALER)
				//dealer.Connect("tcp://" + ser.Map[strconv.Itoa(x.Pid)])
				data, _ := json.Marshal(*x)
				//dealer.Send(string(data), 0)
				ser.SocketMap[str_pid].Send(string(data), 0)
			}
		}
	}
}

func (ser ServerType) GetMsg(enve chan *Envelope) {
	str := "tcp://" + ser.Socket
	//fmt.Println(*flagid,"\nServerdata:-",str,"\nPid=",server.MyPid);
	//fmt.Println("\nServerdata:-",str,"\nPid=",server.MyPid);
	//context, err :=zmq.NewContext();

	socket, err := zmq.NewSocket(zmq.DEALER)
	if err != nil {
		fmt.Println("error:", err)
	}
	socket.Bind(str)

	for {
		//x:= <-enve
		var x Envelope
		//fmt.Println("In InBox(GetMsg)\n")

		data, err := socket.Recv(0)
		err = json.Unmarshal([]byte(data), &x)
		if err != nil {
			fmt.Println("error:", err)
		}

		if x.Pid == -1 {
			//BroadCast
			enve <- &x
		} else if x.Pid == ser.MyPid{
			//fmt.Println("Receive Unicast\n")
			//Unicast
			enve <- &x
		}
		//fmt.Println("In InBox(GetMsg)\n", x)
	}
}

func (ser ServerType) Pid() int {
	return ser.MyPid
}
func (ser ServerType) Peers() []int {
	return ser.PeerIds
}
func (ser ServerType) Outbox() chan *Envelope {
	return ser.out
}

func (ser ServerType) Inbox() chan *Envelope {
	return ser.in
}