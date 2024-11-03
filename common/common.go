package common

import (
	"fmt"
	"time"

	"golang.org/x/exp/rand"
)

type Card struct {
	Suit  string
	Value int
}

// 定义牌堆的结构体
type Deck struct {
	Cards []Card
}

type Player struct {
	Cards []Card
	Total int
}

type Message struct {
	Str    string
	Player *Player
	Dealer *Player
}

func NewDeck() *Deck {
	suits := []string{"♠️", "♥️", "♣️", "♦️"}
	values := []int{2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 1}
	var cards []Card

	for _, suit := range suits {
		for _, value := range values {
			cards = append(cards, Card{Suit: suit, Value: value})
		}
	}

	return &Deck{Cards: cards}
}

func (d *Deck) Shuffle() {
	rand.Seed(uint64(time.Now().UnixNano())) // 以当前时间为随机种子
	for i := range d.Cards {
		j := rand.Intn(len(d.Cards))                    // 随机选择一个索引
		d.Cards[i], d.Cards[j] = d.Cards[j], d.Cards[i] // 交换位置
	}
}

// 发牌
func (d *Deck) Deal(numCards int) []Card {
	if numCards > len(d.Cards) {
		numCards = len(d.Cards)
	}

	hand := d.Cards[:numCards]   // 获取要发的牌
	d.Cards = d.Cards[numCards:] // 从牌堆中移除已发的牌
	//fmt.Printf("%v\n", hand)
	//fmt.Printf("%v\n", d.Cards)
	return hand
}

func (m *Message) PrintMessage() {
	fmt.Printf("Player:%v\nDealer:%v\n%s\n",
		m.Player.Cards, m.Dealer.Cards, m.Str)
}
