package main

import (
	"blackjack/common"
	"fmt"
	"net/rpc"
)

func main() {
	client, err := rpc.Dial("tcp", "localhost:1234")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer client.Close()

	var player common.Player
	var reply common.Message
INIT:
	err = client.Call("Blackjack.InitGame", "player", &reply)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	//reply.PrintMessage()
	fmt.Println(reply.Str)
CONTINUE:
	err = client.Call("Blackjack.DealCards", &player, &reply)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	//reply.PrintMessage()
	fmt.Println(reply.Str)

	// Game loop
	for {
		fmt.Print("Do you want to hit or stand? (h/s): ")
		var choice string
		fmt.Scanln(&choice)

		if choice == "h" {
			err = client.Call("Blackjack.Hit", &player, &reply)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			//reply.PrintMessage()
			fmt.Println(reply.Str)
			if reply.Player.Total > 21 {
				err = client.Call("Blackjack.Stand", &player, &reply)
				if err != nil {
					fmt.Println("Error:", err)
					return
				}
				//reply.PrintMessage()
				fmt.Println("Game Ended:\n", reply.Str)
				break
			}
		} else {
			err = client.Call("Blackjack.Stand", &player, &reply)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			//reply.PrintMessage()
			fmt.Println(reply.Str)
			break
		}
	}
	fmt.Print("Do you want to continue, start a new game or stop? (c/n/s): ")
	var choice string
	fmt.Scanln(&choice)

	if choice == "c" {
		err = client.Call("Blackjack.Continue", &player, &reply)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		goto CONTINUE
	} else if choice == "n" {
		err = client.Call("Blackjack.EndGame", &player, &reply)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		goto INIT
	} else if choice == "s" {
		err = client.Call("Blackjack.EndGame", &player, &reply)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
	}
}
