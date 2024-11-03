package main

import (
	"blackjack/common"
	"fmt"
	"math/rand"
	"net"
	"net/rpc"
	"time"
)

type Blackjack struct {
	Deck   *common.Deck
	Player *common.Player
	Dealer *common.Player
}

func (b *Blackjack) InitGame(playerName string, reply *common.Message) error {
	b.Deck = common.NewDeck()                                   // 创建新的牌堆
	b.Deck.Shuffle()                                            // 洗牌
	b.Player = &common.Player{Cards: []common.Card{}, Total: 0} // 创建玩家
	b.Dealer = &common.Player{Cards: []common.Card{}, Total: 0} // 创建玩家
	(*reply).Str = "Game initialized\n"
	(*reply).Player = b.Player
	(*reply).Dealer = b.Dealer
	return nil
}

func (b *Blackjack) Continue(player *common.Player, reply *common.Message) error {
	b.Player = &common.Player{Cards: []common.Card{}, Total: 0} // 创建玩家
	b.Dealer = &common.Player{Cards: []common.Card{}, Total: 0} // 创建玩家

	if len(b.Deck.Cards) < 10 {
		(*reply).Str = "Deck is not enough, shuffle the deck."
		b.Deck = common.NewDeck() // 创建新的牌堆
		b.Deck.Shuffle()          // 洗牌
	} else {
		(*reply).Str = "Game Continue\n"
	}
	(*reply).Player = b.Player
	(*reply).Dealer = b.Dealer
	return nil
}

func (b *Blackjack) EndGame(player *common.Player, reply *common.Message) error {
	b.Deck = nil   // 清空牌堆
	b.Player = nil // 清空玩家
	b.Dealer = nil
	(*reply).Str = "Game has been ended and all data has been reset.\n"
	return nil
}

func (b *Blackjack) DealCards(player *common.Player, reply *common.Message) error {
	// 直接对 player 进行修改
	b.Player.Cards = append(b.Player.Cards, b.Deck.Deal(2)...)
	b.Player.Total = calculateTotal(b.Player.Cards)

	b.Dealer.Cards = append(b.Dealer.Cards, b.Deck.Deal(2)...)
	b.Dealer.Total = calculateTotal(b.Dealer.Cards)
	(*reply).Str = fmt.Sprintf("You drew: %v, Total: %d\nDealer drews: %v\n",
		b.Player.Cards, b.Player.Total,
		b.Dealer.Cards[:1])
	(*reply).Player = b.Player
	(*reply).Dealer = b.Dealer
	return nil
}

func (b *Blackjack) Hit(player *common.Player, reply *common.Message) error {
	b.Player.Cards = append(b.Player.Cards, b.Deck.Deal(1)...)
	b.Player.Total = calculateTotal(b.Player.Cards)
	(*reply).Str = fmt.Sprintf("You drew: %v, Total: %d\nDealer drews: %v\n",
		b.Player.Cards, b.Player.Total,
		b.Dealer.Cards[:1])
	(*reply).Player = b.Player
	(*reply).Dealer = b.Dealer
	return nil
}

func (b *Blackjack) Stand(player *common.Player, reply *common.Message) error {
	for b.Dealer.Total < 17 {
		b.Dealer.Cards = append(b.Dealer.Cards, b.Deck.Deal(1)...)
		b.Dealer.Total = calculateTotal(b.Dealer.Cards)
	}
	result := determineWinner(b.Player.Total, b.Dealer.Total)
	(*reply).Str = fmt.Sprintf("You drew: %v, Total: %d\nDealer drews: %v, Total: %d\nResult: %s\n",
		b.Player.Cards, b.Player.Total,
		b.Dealer.Cards, b.Dealer.Total,
		result)
	(*reply).Player = b.Player
	(*reply).Dealer = b.Dealer
	return nil
}

func calculateTotal(cards []common.Card) int {
	total := 0
	aces := 0
	for _, card := range cards {
		if card.Value > 10 {
			total += 10
		} else {
			total += card.Value
		}
		if card.Value == 1 { // Ace
			aces++
		}
	}
	for aces > 0 && total <= 11 {
		total += 10 // Count ace as 11 if it doesn't bust
		aces--
	}
	return total
}

func determineWinner(playerTotal, dealerTotal int) string {
	if playerTotal > 21 {
		return "You bust! Dealer wins!"
	} else if dealerTotal > 21 {
		return "Dealer busts! You win!"
	} else if playerTotal > dealerTotal {
		return "You win!"
	} else if playerTotal < dealerTotal {
		return "Dealer wins!"
	}
	return "It's a tie!"
}

func main() {
	rand.Seed(time.Now().UnixNano())
	blackjack := new(Blackjack)
	rpc.Register(blackjack)
	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer listener.Close()
	fmt.Println("Server is running on port 1234...")
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}
		go rpc.ServeConn(conn)
	}
}
