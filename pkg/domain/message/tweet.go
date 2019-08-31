package message

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/mpppk/sutaba-server/pkg/infra/classifier"

	"golang.org/x/xerrors"
)

func PredToText(predict *classifier.ImagePredictResponse) (string, error) {
	conf, err := strconv.ParseFloat(predict.Confidence, 32)
	if err != nil {
		return "", xerrors.Errorf("failed to parse confidence(%s) to float: %w", predict.Confidence)
	}
	confidence := float32(conf)
	predStr := ""
	switch predict.Pred {
	case "sutaba":
		if confidence > 0.8 {
			predStr = GetRandomSutabaHighConfidenceText()
		} else if confidence > 0.5 {
			predStr = GetRandomSutabaMiddleConfidenceText()
		} else {
			predStr = GetRandomSutabaLowConfidenceText()
		}
	case "ramen":
		if confidence > 0.8 {
			predStr = GetRandomRamenHighConfidenceText(confidence)
		} else if confidence > 0.5 {
			predStr = GetRandomRamenMiddleConfidenceText(confidence)
		} else {
			predStr = GetRandomRamenLowConfidenceText(confidence)
		}
	case "other":
		if confidence > 0.8 {
			predStr = GetRandomOtherHighConfidenceText(confidence)
		} else if confidence > 0.5 {
			predStr = GetRandomOtherMiddleConfidenceText(confidence)
		} else {
			predStr = GetRandomOtherLowConfidenceText()
		}
	}

	return predStr, nil
}

func pickRandomStr(texts []string) string {
	rand.Seed(time.Now().Unix())
	return texts[rand.Intn(len(texts))]
}

func GetRandomSutabaHighConfidenceText() string {
	list := []string{"ムムッこのツイート🤔🤔🤔", "オオッこのツイート🤔🤔🤔", "ヤヤッこのツイート🤔🤔🤔",
		"アァーーーッ️️❗❗何だこれはーーッ❗️❗"}
	list1 := []string{"完膚無きまでに", "完全に", "どこからどう見ても", "ここ10年で最高の", "間違いなく"}
	list2 := []string{"この調子でグッドなスタバツイートを心がけるようにッ❗️👮‍♂👮‍♂️️", "市民の協力に感謝するッッッ👮‍👮‍❗❗"}
	return fmt.Sprintf("ピピーッ❗️🔔⚡️スタバ警察です❗️👊👮❗️ %s...%sスタバ❗️❗️😂😂😂\n%s",
		pickRandomStr(list),
		pickRandomStr(list1),
		pickRandomStr(list2),
	)
}

func GetRandomSutabaMiddleConfidenceText() string {
	list := []string{"ムムッこのツイート🤔🤔🤔", "オオッこのツイート🤔🤔🤔", "ヤヤッこのツイート🤔🤔🤔",
		"アァーーーッ️️❗❗何だこれはーーッ❗️❗"}
	list1 := []string{"おおむね", "比較的", "おおよそ"}
	list2 := []string{"😉😉😉", "😙😙😙", "☺️☺️☺️", "🤗🤗🤗"}
	list3 := []string{"グッドなスタバに本官もニッコリ☺️☺️☺️️", "やるじゃん😉😉😉"}
	return fmt.Sprintf("ピピーッ❗️🔔⚡️スタバ警察です❗️👊👮❗️ %s...%sスタバ%s\n%s",
		pickRandomStr(list),
		pickRandomStr(list1),
		pickRandomStr(list2),
		pickRandomStr(list3),
	)
}

func GetRandomSutabaLowConfidenceText() string {
	list := []string{"ムムッこのツイート🤔🤔🤔", "オオッこのツイート🤔🤔🤔", "ヤヤッこのツイート🤔🤔🤔",
		"アァーーーッ️️❗❗何だこれはーーッ❗️❗"}
	list1 := []string{"たぶん", "どことなく", "もしかすると"}
	list2 := []string{"😐😐😐", "😑😑😑"}
	list3 := []string{"ほのかなスタバみを感じる...!", "自信ないけど..おじさん最近目が遠くてな~~~👮‍♀️💦"}
	return fmt.Sprintf("ピピーッ❗️🔔⚡️スタバ警察です❗️👊👮❗️ %s...%s...スタバ......???%s\n%s",
		pickRandomStr(list),
		pickRandomStr(list1),
		pickRandomStr(list2),
		pickRandomStr(list3),
	)
}

func GetRandomRamenHighConfidenceText(confidence float32) string {
	list0 := []string{"ムムッこのツイート🤔🤔🤔", "オオッこのツイート🤔🤔🤔", "ヤヤッこのツイート🤔🤔🤔",
		"アァーーーッ️️❗❗何だこれはーーッ❗️❗"}
	list1 := []string{"完膚無きまでに", "完全に", "どこからどう見ても", "間違いなく"}
	list2 := []string{
		"アレ法のアレに違反してるゾ❗ただちに消しナさぃ❗",
		"ぃますぐ😩😩😩そのツイートを削除しなさい💢💢💢！！！😇😇😇",
		"今スグ消しなｻｲ❗️❗️❗️❗️✌️👮🔫"}
	elm := fmt.Sprintf("%s...%sラーメン❗️❗️😂😂😂\nズルズルズルズル❗❗️❗️❗%s",
		pickRandomStr(list0),
		pickRandomStr(list1),
		pickRandomStr(list2),
	)
	list := []string{
		"ピピーッ❗️🔔⚡️スタバ警察です❗️👊👮❗️" +
			"アナタのツイート💕は❌スタバ法❌第" +
			fmt.Sprintf("%4.f", confidence*1000) + "条🙋" +
			"「ラーメンをスタバなうツイート💕してゎイケナイ❗️」" +
			"に違反しています😡今スグ消しなｻｲ❗️❗️❗️❗️✌️👮🔫",
		"ピピーっ！👮👮スタバ警察です🚨🚨🚨🙅🙅🙅🙅" +
			"そのツイートはツイッター保護法第" +
			fmt.Sprintf("%4.f", confidence*1000) + "条🌟" +
			"「ラーメンで偽スタバなうツイをしてはいけない❗️😡👊🏻」" +
			"に違反しているゾ😤😤😤💢💢💢！！！！！" +
			"ぃますぐ😩😩😩そのツイートを削除しなさい💢💢💢！！！😇😇😇",
		"ピピーッ❗そのツイートゎ☆" +
			"Twitterスタバ部のオキテ第" +
			fmt.Sprintf("%4.f", confidence*1000) + "条・" +
			"「ラーメンのツイートをスタバなうツイート💕しナイ❗❗」" +
			"に違反してるゾ❗ただちに消しナさぃ❗",
		elm, elm, elm,
	}
	return pickRandomStr(list)
}

func GetRandomRamenMiddleConfidenceText(confidence float32) string {
	list0 := []string{"ムムッこのツイート🤔🤔🤔", "オオッこのツイート🤔🤔🤔", "ヤヤッこのツイート🤔🤔🤔",
		"アァーーーッ️️❗❗何だこれはーーッ❗️❗"}
	list1 := []string{"おおむね", "比較的", "おおよそ"}
	list2 := []string{
		"アレ法のアレに違反してるゾ❗ただちに消しナさぃ❗",
		"ぃますぐ😩😩😩そのツイートを削除しなさい💢💢💢！！！😇😇😇",
		"今スグ消しなｻｲ❗️❗️❗️❗️✌️👮🔫"}
	elm := fmt.Sprintf("%s...%sラーメン❗️❗️😂😂😂\nズルズルズルズル❗❗️❗️❗%s",
		pickRandomStr(list0),
		pickRandomStr(list1),
		pickRandomStr(list2),
	)
	list := []string{
		"ピピーッ❗️🔔⚡️スタバ警察です❗️👊👮❗️" +
			"アナタのツイート💕は❌スタバ法❌第" +
			fmt.Sprintf("%4.f", confidence*1000) + "条🙋" +
			"「ラーメン的なアレをスタバなうツイート💕してゎイケナイ❗️」" +
			"に違反しています😡今スグ消しなｻｲ❗️❗️❗️❗️✌️👮🔫",
		"ピピーっ！👮👮スタバ警察です🚨🚨🚨🙅🙅🙅🙅" +
			"そのツイートはツイッター保護法第" +
			fmt.Sprintf("%4.f", confidence*1000) + "条🌟" +
			"「ラーメンみたいな画像で偽スタバなうツイをしてはいけない❗️😡👊🏻」" +
			"に違反しているゾ😤😤😤💢💢💢！！！！！" +
			"ぃますぐ😩😩😩そのツイートを削除しなさい💢💢💢！！！😇😇😇",
		"ピピーッ❗そのツイートゎ☆" +
			"Twitterスタバ部のオキテ第" +
			fmt.Sprintf("%4.f", confidence*1000) + "条・" +
			"「ラーメンっぽい画像をスタバなうツイート💕しナイ❗❗」" +
			"に違反してるゾ❗ただちに消しナさぃ❗",
		elm, elm, elm,
	}
	return pickRandomStr(list)
}

func GetRandomRamenLowConfidenceText(confidence float32) string {
	list0 := []string{"ムムッこのツイート🤔🤔🤔", "オオッこのツイート🤔🤔🤔", "ヤヤッこのツイート🤔🤔🤔",
		"アァーーーッ️️❗❗何だこれはーーッ❗️❗"}
	list1 := []string{"たぶん", "どことなく", "もしかすると"}
	list2 := []string{"😐😐😐", "😑😑😑"}
	list3 := []string{
		"そのうち消しナさぃ❗",
		"近いうちにそのツイートを削除しなさい😇😇😇",
		"今スグとは言わないので消しなｻｲ❗️❗️❗️❗️✌️👮🔫"}
	elm := fmt.Sprintf("%s...%s...ラーメン......???%s%s",
		pickRandomStr(list0),
		pickRandomStr(list1),
		pickRandomStr(list2),
		pickRandomStr(list3),
	)
	list := []string{
		"ピピーッ❗️🔔⚡️スタバ警察です❗️👊👮❗️" +
			"アナタのツイート💕は❌スタバ法❌第" +
			fmt.Sprintf("%4.f", confidence*1000) + "条🙋" +
			"「ラーメン...?に似た何かしらをスタバなうツイート💕してゎイケナイ❗️」" +
			"に違反しています...多分❗気が向いたら消しなｻｲ️✌️👮🔫",
		"ピピーっ！👮👮スタバ警察です🚨🚨🚨🙅🙅🙅🙅" +
			"そのツイートはツイッター保護法第" +
			fmt.Sprintf("%4.f", confidence*1000) + "条🌟" +
			"「ラーメン...?のようなアレで偽スタバなうツイをしてはいけない❗️😡👊🏻」" +
			"に違反している....ような気がする🤔🤔🤔" +
			"お時間がある時で結構ですのでそのツイートを削除しなさい！！！😇😇😇",
		"ピピーッ❗そのツイートゎ☆" +
			"Twitterスタバ部のオキテ第" +
			fmt.Sprintf("%4.f", confidence*1000) + "条・" +
			"「ラーメン...?いやつけ麺...?なんかその辺の...のツイートをスタバなうツイート💕しナイ❗❗」" +
			"に違反してる..という気がしないでもないゾ❗お手すきの際に消しナさぃ❗",
		elm, elm, elm,
	}
	return pickRandomStr(list)
}

func GetRandomOtherHighConfidenceText(confidence float32) string {
	rand.Seed(time.Now().Unix())
	if rand.Intn(20) < 2 {
		return "ピピーッ❗️🔔⚡️スタバ警察です❗️👊👮❗" +
			"アナタのツイート💕は❌スタバ法❌第..." +
			"ピピーッ❗️❗🔔⚡️スタバ警察です❗️👊👮❗アナ..." +
			"ピピーッ❗️❗ピピピピーッッ❗❗❗❗🔔⚡️スタ..." +
			"ピピピッピピピピーッ❗️❗❗❗ピピピピーッッッッ❗❗❗❗❗️"
	}
	list := []string{
		"ピピーッ❗️🔔⚡️スタバ警察です❗️👊👮❗️" +
			"アナタのツイート💕は❌スタバ法❌第" +
			fmt.Sprintf("%4.f", confidence*1000) + "条🙋" +
			"「スタバぢゃないツイートをスタバなうツイート💕してゎイケナイ❗️」" +
			"に違反しています😡今スグ消しなｻｲ❗️❗️❗️❗️✌️👮🔫",
		"ピピーっ！👮👮スタバ警察です🚨🚨🚨🙅🙅🙅🙅" +
			"そのツイートはツイッター保護法第" +
			fmt.Sprintf("%4.f", confidence*1000) + "条🌟" +
			"「偽スタバなうツイをしてはいけない❗️😡👊🏻」" +
			"に違反しているゾ😤😤😤💢💢💢！！！！！" +
			"ぃますぐ😩😩😩そのツイートを削除しなさい💢💢💢！！！😇😇😇",
		"ピピーッ❗そのツイートゎ☆" +
			"Twitterスタバ部のオキテ第" +
			fmt.Sprintf("%4.f", confidence*1000) + "条・" +
			"「スタバぢゃないツイートをスタバなうツイート💕しナイ❗❗」" +
			"に違反してるゾ❗ただちに消しナさぃ❗",
	}
	return pickRandomStr(list)
}

func GetRandomOtherMiddleConfidenceText(confidence float32) string {
	list := []string{"おおむね", "比較的", "おおよそ", "だいぶ", "かなり"}
	templates := []string{
		"ピピーッ❗️🔔⚡️スタバ警察です❗️👊👮❗️" +
			"アナタのツイート💕は❌スタバ法❌第" +
			fmt.Sprintf("%4.f", confidence*1000) + "条🙋" +
			"「%sスタバぢゃないツイートをスタバなうツイート💕してゎイケナイ❗️」" +
			"に違反しています😡今スグ消しなｻｲ❗️❗️❗️❗️✌️👮🔫",
		"ピピーっ！👮👮スタバ警察です🚨🚨🚨🙅🙅🙅🙅" +
			"そのツイートはツイッター保護法第" +
			fmt.Sprintf("%4.f", confidence*1000) + "条🌟" +
			"「%s偽スタバなうツイをしてはいけない❗️😡👊🏻」" +
			"に違反しているゾ😤😤😤💢💢💢！！！！！" +
			"ぃますぐ😩😩😩そのツイートを削除しなさい💢💢💢！！！😇😇😇",
		"ピピーッ❗そのツイートゎ☆" +
			"Twitterスタバ部のオキテ第" +
			fmt.Sprintf("%4.f", confidence*1000) + "条・" +
			"「%sスタバぢゃないツイートをスタバなうツイート💕しナイ❗❗」" +
			"に違反してるゾ❗ただちに消しナさぃ❗",
	}
	return fmt.Sprintf(pickRandomStr(templates), pickRandomStr(list))
}

func GetRandomOtherLowConfidenceText() string {
	list := []string{
		"...何?何これ?スタバではないと思うけども...いやマジで何?",
		"...何?何これ?スタバではないと思うけども...いやマジで何?",
		"...何?何これ?スタバではないと思うけども...いやマジで何?",
		"はい",
	}
	return pickRandomStr(list)
}
