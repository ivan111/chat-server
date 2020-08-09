package main

import (
	_ "github.com/lib/pq"
	"math/rand"
	"strings"
	"time"
)

const sqlGetRandNames = `
SELECT name
FROM names
ORDER BY random()
LIMIT $1
`

func normilizeMessage(msg string) string {
	return strings.Trim(msg, " 　!！、。.?？ー～")
}

type responseFunc func(*Message, *room)

var responseFuncs = map[string]responseFunc{
	"みくじ":        omikuji,
	"おみくじ":       omikuji,
	"おは":         goodMorning,
	"おはよ":        goodMorning,
	"おっは":        goodMorning,
	"オッハ":        goodMorning,
	"おは4":        goodMorning,
	"おは４":        goodMorning,
	"おはよん":       goodMorning,
	"おはよう":       goodMorning,
	"おはようございます":  goodMorning,
	"こんにちわ":      hello,
	"こんにちは":      hello,
	"今日わ":        hello,
	"今日は":        hello,
	"こんちくわ":      hello,
	"こにゃにゃちわ":    hello,
	"こん":         hello,
	"こんばんわ":      goodEvening,
	"こんばんは":      goodEvening,
	"今晩わ":        goodEvening,
	"今晩は":        goodEvening,
	"わんばんこ":      goodEvening,
	"初めまして":      niceToMeetYou,
	"はじめまして":     niceToMeetYou,
	"誰かいますか":     isAnyoneThere,
	"誰かいます":      isAnyoneThere,
	"誰かいる":       isAnyoneThere,
	"誰かおる":       isAnyoneThere,
	"いますか":       isAnyoneThere,
	"いますか誰か":     isAnyoneThere,
	"ヤッホ":        parrotReturn,
	"やっほ":        parrotReturn,
	"おめでとう":      parrotReturn,
	"おめでとうございます": parrotReturn,
	"誕生日おめでとう":   parrotReturn,
	"お誕生日おめでとう":  parrotReturn,
	"誕生日おめでとうございます":  parrotReturn,
	"お誕生日おめでとうございます": parrotReturn,
}

var goodMorningMessages = []string{
	"おはようございます＞ALL",
	"おはようございます＞みなさん1",
	"おはようございます＞$1",
	"おはようございます >$1",
	"おはようございます >$1",
	"おはようございます >$1さん",
	"$1さん、おはようございます",
	"おはようございます",
	"おはようございます",
	"おはようございます",
	"おはようございます",
	"おはようございます",
	"おはようございます",
	"おはようございます",
	"おはようございます。$1さん",
	"おはようございます。 >$1",
	"$1様、おはようございます。",
	"おはようございます。",
	"おはよう!$1さん",
	"おはよう$1さん。元気してる？",
	"おはよう＞$1さん",
	"$1さん、おはよう",
	"おはよう!$1さん",
	"おはよう＞$1",
	"$1さん、おはよう",
	"おはよう!$1さん",
	"おはよう＞$1",
	"$1さん、おはよう",
	"$1ちゃん おはよう",
	"$1君 おはよう",
	"おはよう",
	"おはよう",
	"おはよう",
	"おはよう",
	"おはよう。",
	"おはよう。",
	"おはよう。",
	"おはよう。",
	"おはようさん$1さん",
	"おはよ",
	"('-'*)ｵﾊﾖ♪",
	"(^o^)ﾉ ＜ おはよー",
	"ぉっはょ～!!（*＾－＾）ノ",
	"おっは～",
}

func goodMorning(msg *Message, myRoom *room) {
	sayResponse(msg.Name, myRoom, goodMorningMessages)
}

var helloMessages = []string{
	"こんにちは＞ALL",
	"こんにちわ＞みなさん",
	"こんにちは。＞$1",
	"こんにちは！ >$1",
	"こんにちは >$1",
	"こんにちは >$1",
	"$1さん、こんにちは",
	"こんにちは＞$1",
	"こんにちわ＞$1",
	"こんにちは。＞$1",
	"こんにちは！ >$1さん",
	"こんにちは >$1",
	"こんにちは >$1",
	"$1さん、こんにちは",
	"こんにちは",
	"こんにちは",
	"こんにちは",
	"こんにちは",
	"こんにちは",
	"こんにちは",
	"こんにちは",
	"こんにちは",
	"こんにちは",
	"こんにちは",
	"Hello!",
	"今日は",
	"$1様、こんにちは。",
	"こんにちは$1さん。元気してる？",
	"こんちは＞$1",
	"$1さん、こにゃにゃちわ",
	"こんにちは!$1さん",
	"こんにちは＞$1さん",
	"$1さん、こんにちは",
	"こんにちは!$1さん",
	"こんにちは＞$1さん",
	"$1さん、こんにちは",
	"$1ちゃん こんにちは",
	"||ゝω･)ﾉ こんちゎぁぁ",
	"こんにちは_φ(･ω･｀)",
	"(★´･З･)ﾉ　こんにちは",
}

func hello(msg *Message, myRoom *room) {
	sayResponse(msg.Name, myRoom, helloMessages)
}

var goodEveningMessages = []string{
	"こんばんは＞ALL",
	"こんばんわ＞みなさん",
	"こんばんは。＞$1",
	"こんばんは！ >$1",
	"こんばんは >$1",
	"こんばんは >$1",
	"$1さん、こんばんは",
	"こんばんは＞$1",
	"こんばんわ＞$1",
	"こんばんは。＞$1",
	"こんばんは！ >$1さん",
	"こんばんは >$1",
	"こんばんは >$1",
	"$1さん、こんばんは",
	"こんばんは",
	"こんばんは",
	"こんばんは",
	"こんばんは",
	"こんばんは",
	"こんばんは",
	"こんばんは",
	"こんばんは",
	"こんばんは",
	"こんばんは",
	"Hello!",
	"今晩は",
	"$1様、こんばんは。",
	"こんばんは$1さん。元気してる？",
	"わんばんこ＞$1",
	"こんばんは!$1さん",
	"こんばんは＞$1さん",
	"$1さん、こんばんは",
	"こんばんは!$1さん",
	"こんばんは＞$1さん",
	"$1さん、こんばんは",
	"$1ちゃん こんばんは",
	"こんばんは_φ(･ω･｀)",
	"(★´･З･)ﾉ　こんばんは",
}

func goodEvening(msg *Message, myRoom *room) {
	sayResponse(msg.Name, myRoom, goodEveningMessages)
}

var niceToMeetYouMessages = []string{
	"初めまして",
	"初めまして$1さん",
	"初めまして>$1",
	"$1さん、初めまして",
	"初めまして",
	"初めまして$1さん",
	"初めまして＞$1",
	"初めまして>$1",
	"初めまして >$1",
	"初めまして >$1さん",
	"$1さん、初めまして",
	"はじめまして",
	"はじめまして$1さん",
	"はじめまして>$1",
	"$1さん、はじめまして",
	"はじめまして",
	"はじめまして$1さん",
	"はじめまして$1さん。よろしくお願いいたします。",
	"はじめまして＞$1",
	"はじめまして>$1",
	"はじめまして >$1",
	"はじめまして。よろしくね >$1",
	"はじめまして >$1さん",
	"$1さん、はじめまして",
	"はじめまして",
	"はじめまして$1さん",
	"はじめまして>$1",
	"$1さん、はじめまして",
	"はじめまして",
	"はじめまして$1さん",
	"はじめまして＞$1",
	"はじめまして>$1",
	"はじめまして >$1",
	"はじめまして >$1さん",
	"$1さん、はじめまして",
	"はじめまして(*´-ω-))ﾍﾟｺﾘ",
	"(*￣ω￣)ノはじめまして",
	"はじめまして(ｏ'д')从('д'ｏ)よろしくネ",
}

func niceToMeetYou(msg *Message, myRoom *room) {
	sayResponse(msg.Name, myRoom, niceToMeetYouMessages)
}

var isAnyoneThereMessages = []string{
	"$1さん、いますよ",
	"いますよ。$1さん",
	"いますよ >$1",
	"いますよ>$1",
	"いますよ＞$1",
	"いますよ",
	"いますよ",
	"いますよ",
	"$1さん、いますよ",
	"いますよ。$1さん",
	"いますよ >$1",
	"いますよ>$1",
	"いますよ＞$1",
	"いますよ",
	"いますよ",
	"いますよ",
	"$1さん、いますよ",
	"いますよ。$1さん",
	"いますよ >$1",
	"いますよ>$1",
	"いますよ＞$1",
	"いますよ",
	"いますよ",
	"いますよ",
	"どうしました？",
	"何ですか？",
	"はいはーい",
	"いるよ",
}

func isAnyoneThere(msg *Message, myRoom *room) {
	sayResponse(msg.Name, myRoom, isAnyoneThereMessages)
}

func parrotReturn(msg *Message, myRoom *room) {
	sayResponse(msg.Name, myRoom, []string{msg.Message})
}

func omikuji(msg *Message, myRoom *room) {
	i := rand.Intn(100)
	var res string

	switch {
	case i == 0:
		res = "超大吉"
	case i == 1:
		res = "スーパーハイパーメガマックス大吉"
	case i == 2:
		res = "超大凶"
	case i == 3:
		res = "スーパーハイパーメガマックス大凶"
	case i < 25:
		res = "大吉"
	case i < 35:
		res = "中吉"
	case i < 48:
		res = "小吉"
	case i < 72:
		res = "吉"
	case i < 91:
		res = "末吉"
	default:
		res = "凶"
	}

	sendMessage(myRoom, "おみくじ", res, "#000")
}

func sayResponse(fromName string, myRoom *room, messagesArr []string) {
	n := rand.Intn(20) + 5

	rows, err := db.Query(sqlGetRandNames, n)
	if err != nil {
		return
	}

	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return
		}

		ci := rand.Intn(len(colors))
		color := colors[ci]

		ms := (time.Second * 3) + (time.Millisecond * time.Duration(rand.Intn(20000)))

		mi := rand.Intn(len(messagesArr))
		msg := messagesArr[mi]
		msg = strings.Replace(msg, "$1", fromName, 1)

		go func() {
			time.Sleep(ms)
			sendMessage(myRoom, name, msg, color)
		}()
	}
}
