package main

import (
	"community/pkg/aiapi/deepseek"
	"github.com/zeromicro/go-zero/core/logx"
)

const TestMessage = `
	The user will provide some exam text. Please parse the "question" and "answer" and output them in JSON format.
	你需要考虑多方面，比如天气，路程，精力，路途风景，城市文化等

	EXAMPLE INPUT:
	“我计划下个月去云南大理玩 4 天，想住靠洱海的民宿，预算大概 3000 元 / 人。第一天想先去大理古城逛吃，下午去崇圣寺三塔；第二天希望环洱海玩，最好能租电动车，晚上想在双廊看日落；第三天去喜洲古镇，想体验做扎染，顺便买些手工乳扇；第四天上午去蝴蝶泉，下午坐高铁返程。对了，听说大理紫外线强，还有要注意高原反应吗？”

	EXAMPLE JSON OUTPUT:
	{
  "travelAnalysisResult": 
{
    "basicItineraryInfo": {
      "destination": "云南大理",
      "itineraryDays": "4天",
      "travelTime": "下个月（未指定具体日期）",
      "budgetPerPerson": "3000元",
      "coreDemands": "环洱海游玩、体验白族文化、入住洱海周边民宿",
      "isAiSupplemented": false
    },
    "dailyDetailedItinerary": [
      {
        "itineraryDate": "行程第1天",
        "itineraryTheme": "古城与文化地标",
        "activities": [
          {
            "projectName": "大理古城",
            "estimatedDuration": "3小时（上午）",
            "coreExperience": "漫步古城街巷，品尝白族特色小吃（如烤乳扇、饵块）",
            "isUserSpecified": true,
            "isAiSupplemented": false
          },
          {
            "projectName": "崇圣寺三塔",
            "estimatedDuration": "2小时（下午）",
            "coreExperience": "参观唐代佛塔，拍摄三塔倒影",
            "isUserSpecified": true,
            "supplementaryTips": "建议通过官方公众号提前预约门票，避开12:00-14:00人流高峰",
            "isAiSupplemented": true
          }
        ],
        "diningRecommendations": {
          "lunch": "推荐：大理古城内的“石井私房菜”",
          "dinner": "推荐：古城南门附近的白族石板烧",
          "isAiSupplemented": true
        }
      },
      {
        "itineraryDate": "行程第2天",
        "itineraryTheme": "洱海深度游",
        "activities": [
          {
            "projectName": "洱海骑行（电动车）",
            "estimatedDuration": "5小时（全天）",
            "coreExperience": "从古城出发，途经才村码头、磻溪村S弯、喜洲古镇，最终抵达双廊",
            "isUserSpecified": true,
            "supplementaryTips": "电动车租金约50-80元/天，需提前检查电量，携带防晒用品",
            "isAiSupplemented": true
          },
          {
            "projectName": "双廊古镇观日落",
            "estimatedDuration": "1小时（傍晚18:00-19:00）",
            "coreExperience": "在双廊玉几岛或太阳宫附近观赏洱海日落",
            "isUserSpecified": true,
            "isAiSupplemented": false
          }
        ],
        "diningRecommendations": {
          "lunch": "推荐：喜洲古镇内的“喜林苑”（招牌菜为黄焖鸡）",
          "dinner": "推荐：双廊的“云途海景餐厅”",
          "isAiSupplemented": true
        }
      },
      {
        "itineraryDate": "行程第3天",
        "itineraryTheme": "古镇与非遗体验",
        "activities": [
          {
            "projectName": "喜洲古镇",
            "estimatedDuration": "2小时（上午）",
            "coreExperience": "参观白族古建筑，漫步四方街",
            "isUserSpecified": true,
            "isAiSupplemented": false
          },
          {
            "projectName": "白族扎染体验",
            "estimatedDuration": "1.5小时（下午）",
            "coreExperience": "跟随当地手艺人学习扎染工艺，制作小方巾或T恤",
            "isUserSpecified": true,
            "supplementaryTips": "推荐“蓝续扎染坊”，体验费用约80-120元/人，作品可带走",
            "isAiSupplemented": true
          },
          {
            "projectName": "购买手工乳扇",
            "estimatedDuration": "30分钟",
            "coreExperience": "在喜洲古镇主街购买新鲜手工乳扇，可搭配玫瑰酱食用",
            "isUserSpecified": true,
            "isAiSupplemented": false
          }
        ],
        "diningRecommendations": {
          "lunch": "推荐：喜洲古镇内的“转角楼饭店”（尝试白族特色“土八碗”）",
          "isAiSupplemented": true
        }
      },
      {
        "itineraryDate": "行程第4天",
        "itineraryTheme": "自然景观与返程",
        "activities": [
          {
            "projectName": "蝴蝶泉公园",
            "estimatedDuration": "1.5小时（上午）",
            "coreExperience": "参观蝴蝶泉，观赏泉边古树；春季可看到蝴蝶群（其他季节蝴蝶较少）",
            "isUserSpecified": true,
            "supplementaryTips": "门票约40元/人，建议上午前往，避开午后高温",
            "isAiSupplemented": true
          },
          {
            "projectName": "高铁返程",
            "estimatedDuration": "未指定（根据目的地不同有所差异）",
            "coreExperience": "从大理站乘坐高铁离开",
            "isUserSpecified": true,
            "isAiSupplemented": false
          }
        ]
      }
    ],
    "transportationPlan": {
      "outboundTransport": "未指定（仅提及返程为高铁）",
      "localTransport": [
        {
          "transportMode": "电动车租赁",
          "applicableScenario": "洱海游玩（第2天）",
          "costReference": "50-80元/天",
          "isUserSpecified": true,
          "isAiSupplemented": false
        },
        {
          "transportMode": "出租车/网约车",
          "applicableScenario": "景点间短途往返（如古城到崇圣寺三塔）",
          "costReference": "15-30元/次",
          "isUserSpecified": false,
          "isAiSupplemented": true
        }
      ],
      "returnTransport": {
        "transportMode": "高铁",
        "departureLocation": "大理站",
        "destination": "未指定",
        "isUserSpecified": true,
        "isAiSupplemented": false
      }
    },
    "accommodationSuggestions": {
      "accommodationType": "洱海周边民宿",
      "locationPreference": "洱海附近（未指定具体区域；AI补充：推荐双廊或才村码头周边，观景视角佳）",
      "budgetRange": "未指定单晚价格（根据4天行程及人均3000元总预算，推算单晚约750元/人）",
      "recommendedOptions": [
        "双廊“六阅·无所”民宿：一线海景房，含早餐，靠近双廊古镇",
        "才村“大理苍海觅踪客栈”：性价比高，步行5分钟可达洱海边码头"
      ],
      "isUserSpecified": true,
      "isAiSupplemented": true
    },
    "travelNotes": [
      {
        "noteType": "气候防护",
        "details": "大理紫外线较强，需携带高倍数防晒霜（SPF50+）、遮阳帽及太阳镜",
        "isUserSpecified": true,
        "isAiSupplemented": false
      },
      {
        "noteType": "高原反应",
        "details": "大理海拔约2000米，多数人无明显高原反应；若出现头痛、乏力等症状，建议减少剧烈活动、多喝温水，严重时及时就医",
        "isUserSpecified": true,
        "isAiSupplemented": true
      },
      {
        "noteType": "文化礼仪",
        "details": "进入白族民居需脱鞋，避免随意触碰祭祀物品；拍摄当地居民前建议先征得同意",
        "isUserSpecified": false,
        "isAiSupplemented": true
      },
      {
        "noteType": "物品准备",
        "details": "洱海周边风力较大，建议携带薄外套；扎染体验可能弄脏衣物，可准备旧衣服或围裙",
        "isUserSpecified": false,
        "isAiSupplemented": true
      }
    ],
    "featuredExperienceRecommendations": [
      {
        "experienceName": "洱海游船",
        "recommendationReason": "从下关码头乘坐游船前往南诏风情岛，欣赏洱海全景；春季可观赏苍山雪景",
        "suitableDay": "行程第2天（洱海游玩当天）",
        "isAiSupplemented": true
      },
      {
        "experienceName": "白族三道茶体验",
        "recommendationReason": "在大理古城或喜洲古镇品尝“一苦二甜三回味”的三道茶，深入了解白族文化",
        "suitableDay": "行程第1天（大理古城游玩期间）",
        "isAiSupplemented": true
      },
      {
        "experienceName": "苍山徒步",
        "recommendationReason": "若时间充裕，可增加苍山洗马潭大索道徒步项目，观赏高山杜鹃（4-5月最佳）",
        "suitableDay": "可调整至行程第3天下午（需缩短喜洲古镇游玩时间）",
        "isAiSupplemented": true
      }
    ]
  }
}
`

func main() {
	deepSeekClient := deepseek.NewClient("sk-4be9e334908f4c2087cf68df508883f3") // 测试key
	resp, err := deepSeekClient.SetModel(deepseek.ModelDeepSeekReasoner).SetMaxTokens(4096).SetTemperature(0.3).SetTopP(0.5).
		AddSystemMessage(TestMessage).AddUserMessage("我想从北京去西安玩").SetResponseFormat(&deepseek.ResponseFormat{Type: deepseek.ResponseFormatJSON}).Send()
	if err != nil {
		logx.Errorf(err.Error())
		return
	}
	logx.Infof(resp.Choices[0].Message.Content)
}
