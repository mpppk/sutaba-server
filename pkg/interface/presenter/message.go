package presenter

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"golang.org/x/xerrors"

	"github.com/mpppk/messagen/messagen"
	domain "github.com/mpppk/sutaba-server/pkg/domain/service"
)

type MessagenType string

const (
	RootType              = "Root"
	TweetCheckType        = "TweetCheck"
	SutabaDescriptionType = "SutabaDescription"
	GoodEmojiType         = "GoodEmoji"
	LastMessageType       = "LastMessage"
	ExclamationType       = "Exclamation"
	ThinkingEmojiType     = "ThinkingEmoji"
	ConfidenceType        = "Confidence"
	ConfidenceHighValue   = "High"
	ConfidenceMediumValue = "Medium"
	ConfidenceLowValue    = "Low"
	ClassType             = "Class"
	ClassSutabaValue      = "sutaba"
	ClassRamenValue       = "ramen"
	ClassOtherValue       = "other"
	DebugType             = "Debug"
	DebugOnValue          = "on"
	DebugOffValue         = "off"
	RuleNumType           = "RuleNum"
)

func generateResultMessage(result *domain.ClassifyResult) (string, error) {
	confidence := float32(result.Confidence)
	generator, err := messagen.New(nil)
	if err != nil {
		return "", xerrors.Errorf("failed to generate messagen instance: %w", err)
	}

	if err := generator.AddDefinition(getMessagenDefinitions()...); err != nil {
		return "", xerrors.Errorf("failed to add definitions to messagen: %w", err)
	}

	state := map[string]string{
		"Class": result.Class,
	}

	if confidence > 0.8 {
		state["Confidence"] = "High"
	} else if confidence > 0.5 {
		state["Confidence"] = "Medium"
	} else {
		state["Confidence"] = "Low"
	}

	fmt.Println("state", state)

	messages, err := generator.Generate("Root", state, 1)
	if err != nil {
		return "", err
	}
	return messages[0], nil
}

func toTemplateVariable(v string) string {
	return "{{." + v + "}}"
}

func toTemplate(tmpl string, variables ...string) string {
	var newVariables []interface{}
	for _, v := range variables {
		newVariables = append(newVariables, toTemplateVariable(v))
	}
	return fmt.Sprintf(tmpl, newVariables...)
}

func w(v string) string {
	return toTemplateVariable(v)
}

func getMessagenDefinitions() []*messagen.Definition {
	return []*messagen.Definition{
		{
			Type: RootType,
			Templates: []string{
				toTemplate("ピピーッ❗️🔔⚡️スタバ警察です❗️👊👮❗️\n%s%sスタバ❗️❗️%s\n%s",
					TweetCheckType, SutabaDescriptionType, GoodEmojiType, LastMessageType),
			},
			Constraints: map[string]string{ClassType: ClassSutabaValue},
		},
		{
			Type: RootType,
			Templates: []string{toTemplate("ピピーッ❗️🔔⚡️スタバ警察です❗️👊👮❗️\n"+
				"アナタのツイート💕は❌スタバ法❌第%s条🙋\n"+
				"「スタバぢゃないツイートをスタバなうツイート💕してゎイケナイ❗️」\n"+
				"に違反しています😡今スグ消しなｻｲ❗️❗️❗️❗️✌️👮🔫\n", RuleNumType)},
			Constraints: map[string]string{
				ClassType:            ClassSutabaValue,
				ConfidenceType + "/": ConfidenceHighValue + "|" + ConfidenceMediumValue},
		},
		{
			Type: RootType,
			Templates: []string{toTemplate(
				`{"class": "%s", "confidence": "%s"}`, ClassType, ConfidenceType)},
			Constraints: map[string]string{DebugType + ":1": DebugOnValue},
		},
		{
			Type: TweetCheckType,
			Templates: []string{
				w(ExclamationType) + "このツイート" + strings.Repeat(w(ThinkingEmojiType), 3) + "...",
			},
		},
		{
			Type:      TweetCheckType,
			Templates: []string{"アアーーー❗️なんだこれはーーー❗️❗️"},
			Weight:    0.5,
		},
		{Type: ExclamationType, Templates: []string{"ムムッ", "ヤヤッ", "オオッ"}},
		{Type: ThinkingEmojiType, Templates: []string{"🤔", "🤨"}},
		{Type: GoodEmojiType, Templates: []string{"😆", "😂"}},
		{
			Type:        SutabaDescriptionType,
			Templates:   []string{"完全に", "間違いなく"},
			Constraints: map[string]string{ConfidenceType: ConfidenceHighValue},
		},
		{
			Type:        SutabaDescriptionType,
			Templates:   []string{"おそらく", "多分"},
			Constraints: map[string]string{ConfidenceType: ConfidenceMediumValue},
		},
		{
			Type: LastMessageType,
			Templates: []string{"この調子でグッドなスタバツイートを心がけるようにッ❗️👮‍👮‍",
				"市民の協力に感謝するッッッ👮‍👮‍❗"},
		},
	}
}

func classAndConfidenceToText(className string, confidence float32) string {
	predStr := ""
	switch className {
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

	return predStr
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
