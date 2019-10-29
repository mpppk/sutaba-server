package presenter

import (
	"fmt"
	"strings"

	"golang.org/x/xerrors"

	"github.com/mpppk/messagen/messagen"
	domain "github.com/mpppk/sutaba-server/pkg/domain/service"
)

type MessagenType string

const (
	RootType              = "Root"
	BeginningMessageType  = "BeginningMessage"
	TweetCheckType        = "TweetCheck"
	SutabaDescriptionType = "SutabaDescription"
	GoodEmojiType         = "GoodEmoji"
	BadEmojiType          = "BadEmoji"
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
	TargetNameType        = "TargetName"
	RamenSuffixType       = "RamenSuffix"
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
		ClassType: result.Class,
	}

	if confidence > 0.8 {
		state[ConfidenceType] = ConfidenceHighValue
	} else if confidence > 0.5 {
		state[ConfidenceType] = ConfidenceMediumValue
	} else {
		state[ConfidenceType] = ConfidenceLowValue
	}

	messages, err := generator.Generate(RootType, state, 1)
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
				w(BeginningMessageType) +
					w(TweetCheckType) + w(SutabaDescriptionType) + "スタバ❗️❗️" + w(GoodEmojiType) + "\n" +
					w(LastMessageType),
			},
			Constraints: map[string]string{ClassType: ClassSutabaValue},
		},
		{
			Type: RootType,
			Templates: []string{
				w(BeginningMessageType) +
					w(TweetCheckType) + w(SutabaDescriptionType) + "ラーメン❗️❗️" + w(GoodEmojiType) + "\n" +
					"ズルズルズルズル❗❗️❗️❗" + w(LastMessageType),
			},
			Constraints: map[string]string{
				ClassType:            ClassRamenValue,
				ConfidenceType + "/": ConfidenceHighValue + "|" + ConfidenceMediumValue,
			},
		},
		{
			Type: RootType,
			Templates: []string{
				toTemplate("%sアナタのツイート💕は❌スタバ法❌第%s条🙋"+
					"「%sをスタバなうツイート💕してゎイケナイ❗️」"+
					"に違反しています😡今スグ消しなｻｲ❗️❗️❗️❗️✌️👮🔫",
					BeginningMessageType, RuleNumType, TargetNameType),

				toTemplate(
					"%sそのツイートはツイッター保護法第%s条🌟"+
						"「%sで偽スタバなうツイをしてはいけない❗️😡👊🏻」"+
						"に違反しているゾ😤😤😤💢💢💢！！！！！"+
						"ぃますぐ😩😩😩そのツイートを削除しなさい💢💢💢！！！😇😇😇",
					BeginningMessageType, RuleNumType, TargetNameType),

				toTemplate(
					"ピピーッ❗そのツイートゎ☆Twitterスタバ部のオキテ第%s条・"+
						"「%sのツイートをスタバなうツイート💕しナイ❗❗」"+
						"に違反してるゾ❗ただちに消しナさぃ❗", RuleNumType, TargetNameType),
			},
			Constraints: map[string]string{
				ClassType + "/":      ClassRamenValue + "|" + ClassOtherValue,
				ConfidenceType + "/": ConfidenceHighValue + "|" + ConfidenceMediumValue,
			},
		},
		{
			Type: RootType,
			Templates: []string{
				toTemplate("%sアナタのツイート💕は❌スタバ法❌第%s条🙋"+
					"「ラーメン...?に似た何かしらをスタバなうツイート💕してゎイケナイ❗️」"+
					"に違反しています...多分❗気が向いたら消しなｻｲ️✌️👮🔫",
					BeginningMessageType, RuleNumType),
				toTemplate("%sそのツイートはツイッター保護法第%s条🌟"+
					"「ラーメン...?のようなアレで偽スタバなうツイをしてはいけない❗️😡👊🏻」"+
					"に違反している....ような気がする🤔🤔🤔"+
					"お時間がある時で結構ですのでそのツイートを削除しなさぃ！！！😇😇😇",
					BeginningMessageType, RuleNumType,
				),
				toTemplate("ピピーッ❗そのツイートゎ☆Twitterスタバ部のオキテ第%s条・"+
					"「ラーメン...?いやつけ麺...?なんかその辺のツイートをスタバなうツイート💕しナイ❗❗」"+
					"に違反してる予感がぁりマス❗お手すきの際に消しナさぃ❗", RuleNumType,
				),
			},
			Constraints: map[string]string{
				ClassType:      ClassRamenValue,
				ConfidenceType: ConfidenceLowValue,
			},
		},
		{
			Type: RootType,
			Templates: []string{toTemplate("%sアナタのツイート💕は❌スタバ法❌第%s条🙋\n"+
				"「スタバぢゃないツイートをスタバなうツイート💕してゎイケナイ❗️」\n"+
				"に違反しています😡今スグ消しなｻｲ❗️❗️❗️❗️✌️👮🔫\n", BeginningMessageType, RuleNumType)},
			Constraints: map[string]string{
				ClassType:            ClassOtherValue,
				ConfidenceType + "/": ConfidenceHighValue + "|" + ConfidenceMediumValue},
		},
		{
			Type:      RootType,
			Templates: []string{"...何?何これ?スタバではないと思うけども...いやマジで何?"},
			Constraints: map[string]string{
				ClassType:      ClassOtherValue,
				ConfidenceType: ConfidenceLowValue,
			},
		},
		{
			Type: RootType,
			Templates: []string{toTemplate(
				`{"class": "%s", "confidence": "%s"}`, ClassType, ConfidenceType)},
			Constraints: map[string]string{DebugType + ":1": DebugOnValue},
		},
		{
			Type: BeginningMessageType,
			Templates: []string{
				"ピピーッ❗️🔔⚡️スタバ警察です❗️👊👮❗️\n",
				"ピピーっ！👮👮スタバ警察です🚨🚨🚨🙅🙅🙅🙅\n",
			},
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
			Templates:   []string{"完膚無きまでに", "完全に", "どこからどう見ても", "ここ10年で最高の", "間違いなく"},
			Constraints: map[string]string{ConfidenceType: ConfidenceHighValue},
		},
		{
			Type:        SutabaDescriptionType,
			Templates:   []string{"おそらく", "おおむね", "比較的", "おおよそ"},
			Constraints: map[string]string{ConfidenceType: ConfidenceMediumValue},
		},
		{
			Type:        SutabaDescriptionType,
			Templates:   []string{"たぶん", "どことなく", "もしかすると"},
			Constraints: map[string]string{ConfidenceType: ConfidenceLowValue},
		},
		{
			Type:        TargetNameType,
			Templates:   []string{"ラーメン" + w(RamenSuffixType)},
			Constraints: map[string]string{ClassType: ClassRamenValue},
		},
		{
			Type:        TargetNameType,
			Templates:   []string{"スタバぢゃない画像"},
			Constraints: map[string]string{ClassType: ClassOtherValue},
		},
		{
			Type:        RamenSuffixType,
			Templates:   []string{""},
			Constraints: map[string]string{ClassType + "/": ClassSutabaValue + "|" + ClassOtherValue},
		},
		{
			Type:        RamenSuffixType,
			Templates:   []string{"的なアレ", "みたいな画像", "っぽい画像"},
			Constraints: map[string]string{ClassType: ClassRamenValue},
		},
		{
			Type: LastMessageType,
			Templates: []string{"この調子でグッドなスタバツイートを心がけるようにッ❗️👮‍👮‍",
				"市民の協力に感謝するッッッ👮‍👮‍❗",
				"グッドなスタバに本官もニッコリ☺️☺️☺️️", "やるじゃん😉😉😉",
			},
			Constraints: map[string]string{ClassType: ClassSutabaValue},
		},
		{
			Type: LastMessageType,
			Templates: []string{"アレ法のアレに違反してるゾ❗ただちに消しナさぃ❗",
				"ぃますぐ😩😩😩そのツイートを削除しなさい💢💢💢！！！😇😇😇",
				"今スグ消しなｻｲ❗️❗️❗️❗️✌️👮🔫",
			},
			Constraints: map[string]string{ClassType + "/": ClassRamenValue + "|" + ClassOtherValue},
		},
	}
}
