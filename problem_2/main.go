/*
Write a program that demonstrates quorum election. The program should have a
specified number of members in the quorum and start an interactive mode for the
quorum election game.
a. Game steps:
i. Start the quorum with N members.
ii. Elect one of the members as the quorum leader.
iii. Each member sends heartbeat signals to each other to ensure they are
alive.
iv. Identify a member that has failed to respond to the heartbeat by voting.
1. Remove the failed member from the quorum.
2. If the failed member was the leader, go back to step ii.
b. Each member should have an ID starting from 0, 1, 2, and so on.
c. The command "kill 0" should make member 0 unresponsive to others.
d. There are multiple quorum mechanisms available, and you can design a
better one according to your requirements. (Hint: Consensus Algorithm, or
Centralization)

launch binary with specified number of member
./main 3
> Starting quorum with 3 members
> Member 0: Hi
> Member 1: Hi
> Member 2: Hi
> Member 0: I want to be leader
> Member 2: Accept member 0 to be leader
> Member 1: I want to be leader
> Member 1: Accept member 0 to be leader
> Member 0 voted to be leader: (2 > 3/2)
> kill 1
> Member 0: failed heartbeat with Member 1
> Member 2: failed heartbeat with Member 1
> Member 1: kick out of quorum: (2 > current/2)
> kill 2
> Member 0: failed heartbeat with Member 1
> Member 0: no response from other users(timeout)
> Member 2: kick out of quorum: leader decision
> Quorum failed: (1 > total/2)
*/
// Raft or Paxos

// https://thesecretlivesofdata.com/raft/

package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type State string
type Member struct {
	id            int
	state         State
	hasVoted      bool
	socket        chan string
	ack           chan string
	alive         bool
	lastHeartbeat time.Time
	sendHeartbeat chan string
}

var electionMutex sync.Mutex
var leaderExists bool

const (
	Follower  State = "Follower"
	Candidate State = "Candidate"
	Leader    State = "Leader"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Please enter a number")
		return
	}

	size, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Println("The parameter is not a valid number:", err)
		return
	}
	fmt.Println("Starting quorum with", size, "members")

	wg := sync.WaitGroup{}
	leaderExists = false
	members := make([]*Member, size)

	for i := 0; i < size; i++ {
		wg.Add(1)
		members[i] = &Member{
			id:            i,
			state:         Follower,
			socket:        make(chan string),
			ack:           make(chan string),
			sendHeartbeat: make(chan string),
			alive:         true,
		}
		go members[i].citizen(&wg, members)
	}
	wg.Wait()
	killMember(members)
}

func (m *Member) citizen(wg *sync.WaitGroup, allMembers []*Member) {
	defer wg.Done()
	m.state = Follower
	m.hasVoted = false
	// Let the voting record return to zero

	if m.state == Follower && m.alive == true {
		fmt.Println("Member", m.id, ": Hi")
		time.Sleep(100 * time.Millisecond)
		if leaderExists == false {
			m.follower(allMembers) // Enter the follower state and prepare to start the election
		}
	}
}

func (m *Member) follower(allMembers []*Member) {
	votedChan := make(chan bool)
	go func(m *Member, votedChan chan bool) {
		for i := 0; i < 1; i++ {
			msg := <-m.socket
			var candidateId int
			fmt.Sscanf(msg, "Please vote for me, I am %d", &candidateId)
			if !m.hasVoted && m.alive {
				m.ack <- fmt.Sprintf("Accept member %d to be leader", candidateId)
				fmt.Printf("Member %d: Accept member %d to be leader\n", m.id, candidateId)
				m.hasVoted = true
				votedChan <- true
			}
		}
	}(m, votedChan)

	// Make sure that once you have voted, you cannot become a Candidate to start an election again
	select {
	case <-votedChan:
		return
	case <-time.After(100 * time.Millisecond):
		if !m.hasVoted {
			m.state = Candidate
			fmt.Println("Member", m.id, ": I want to be leader")
			electionMutex.Lock()
			defer electionMutex.Unlock()
			m.leaderElection(allMembers)
			if leaderExists == true && m.state != Leader {
				m.state = Follower
			}
		}
	}
	if m.state == Follower && leaderExists == true {
		go m.receiveHeartbeat(allMembers) // Leaders do heartbeat checks
	}
}

func (m *Member) leaderElection(allMembers []*Member) {
	if leaderExists == false {
		getVotes := 1
		for i := 0; i < len(allMembers); i++ {
			if m.id != i {
				msg := fmt.Sprintf("Please vote for me, I am %d", m.id)
				allMembers[i].socket <- msg
			}
		} // Tell others that I want to be a leader
		for i := 0; i < len(allMembers); i++ {
			if m.id != i {
				select {
				case msg := <-allMembers[i].ack:
					if msg == fmt.Sprintf("Accept member %d to be leader", m.id) {
						getVotes++
					}
				case <-time.After(2 * time.Second):
					_ = fmt.Sprintf("Member %d did not vote within the time limit, skip", i)
				}
			}
		}

		numSurvivor := 0
		for _, survivor := range allMembers {
			if survivor.alive == true {
				numSurvivor++
			}
		} // Count the survivors to confirm the legal number of voters
		if getVotes >= (numSurvivor / 2) {
			fmt.Println("Member", m.id, "voted to be leader: ", getVotes, ">=", (numSurvivor / 2))
			leaderExists = true
			m.state = Leader
			if m.state == Leader {
				go m.heartbeatCheck(allMembers)
			}
		} else {
			fmt.Println("Member", m.id, ": not enough votes to become leader,only have:", getVotes)
		}
	}
}

func (m *Member) heartbeatCheck(allMembers []*Member) {
	for {
		for _, member := range allMembers {
			if member.id != m.id && member.alive {
				member.sendHeartbeat <- fmt.Sprintf("heartbeat from leader:Member%d", m.id)
			}
		}
		time.Sleep(1 * time.Second)
	} // Send a heartbeat check to all alive members
}

func (m *Member) receiveHeartbeat(allMembers []*Member) {
	for {
		select {
		case msg := <-m.sendHeartbeat:
			for _, member := range allMembers {
				if msg == fmt.Sprintf("heartbeat from leader:Member%d", member.id) {
					m.lastHeartbeat = time.Now()
				}
			} // Receive heartbeat checks from the leader
		case <-time.After(2 * time.Second):
			if time.Since(m.lastHeartbeat) > 3*time.Second {
				if leaderExists == false {
					fmt.Println("The leader no longer exists and needs to be re-elected")
					break
				}
			} // Timeout without receiving a heartbeat check from the leader
		}
	}
}

func killMember(allMembers []*Member) {
	wg1 := sync.WaitGroup{}
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("-> ")
		text, _ := reader.ReadString('\n')

		text = strings.TrimSpace(text)

		if text == "hi" {
			fmt.Println("hello")
		} else if text == "exit" {
			break
		} else {
			parts := strings.Fields(text)
			idStr := parts[1]
			id, err := strconv.Atoi(idStr)
			if err != nil {
				fmt.Println("ID is not a valid number")
			} else {
				fmt.Println("The member you want to kill is:", id)
				for _, member := range allMembers {
					if member.id == id {
						if member.state == Leader { // Confirm whether the member being killed is the leader
							member.alive = false
							leaderExists = false
							for _, reElectionMembers := range allMembers {
								wg1.Add(1)
								go reElectionMembers.citizen(&wg1, allMembers) // If it is the leader be killed, re-elect
							}
						} else {
							member.alive = false
						}
					}
				}
			}
		}
	}
}
