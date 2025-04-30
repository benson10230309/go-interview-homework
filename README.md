# go-interview-homework
Golang interview homework: Raft election, RAID simulator &amp; Q&amp;A game

# Go Project - Backend System Engineer (Golang)

This repository includes three Go projects:

1. **Math Game with Students and Teacher**
   Simulates a math quiz competition using goroutines and synchronization.

2. **Quorum Leader Election**
   Demonstrates quorum-based leader election and re-election when a leader fails.

3. **RAID Simulator (RAID0, RAID1, RAID10, RAID5, RAID6)**
   Simulates writing and rebuilding data across multiple RAID levels.

## How to Run

```bash
For question 1:

cd problem_1
go run main.go

For question 2:

cd problem_2
go run main.go 5
# Type "kill leader_number" in console to test re-election

For question 3:

cd problem_3
go run main.go

Make sure inputSample.txt exists in the RAID folder.

For question 1bonus1:

cd problem_1_bonus1
go run main.go

For question 1bonus2:

cd problem_1_bonus2
go run main.go
