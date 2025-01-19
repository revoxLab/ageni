import {range} from 'lodash-es'
import {defineApi, defineMutation, defineQuery} from './define'
import {Response, StudioMessage} from './type'

export * from './common'
export * from './message'
export type * from './type'

export const {
  useRequest: useConversationListQuery,
  request: conversationListQuery,
} = defineApi({
  method: 'post',
  url: '/v1/studio/conversation_list',
  mockParams: {
    bot_id: 1,
    page_size: 10,
    page: 1,
  },
  mockData: {
    code: 0,
    data: [
      {
        id: 1,
        title: 'title',
      },
      {
        id: 2,
        title: 'title',
      },
    ],
    message: 'process success',
  },
})

export const createStudioConversationMutation = defineMutation({
  url: '/v1/studio/create_conversation',
  mockParams: {} as {
    bot_id: number
    title: string
    type?: 'draft'
  },
  mockData: {
    code: 0,
    data: {
      conversation_id: 1,
      bot_id: 1,
      created_at: '2023-04-28T10:00:00Z',
      title: 'xxx',
    },
    message: 'process success',
  },
})

export const useStudioConversationQuery = defineQuery<
  {conversation_id: number; page: number; page_size: number},
  Response<{
    conversation_id: number
    messages?: StudioMessage[]
  }>
>({
  method: 'post',
  url: '/v1/studio/conversation_history',
  mockParams: {
    conversation_id: 1,
    page_size: 10,
    page: 1,
  },
  mockData: {
    code: 0,
    data: {
      conversation_id: 11,
      messages: [],
    },
    message: 'success',
  },
})

export const deleteConversationMutation = defineMutation({
  url: '/v1/studio/conversation_delete',
  mockParams: {
    id: 1,
  },
})

export const useAgentListQuery = defineQuery({
  method: 'post',
  url: '/v1/studio/bot_list',
  mockParams: {} as {
    tab: string
    page: number
    page_size: number
    keywords: string
    pick_type?: number
  },
  mockData: ({page, page_size}) => ({
    code: 0,
    data: {
      bots: range(0, page_size).map((i) => ({
        id: 1,
        status: 1,
        name: 'smart_wallet' + ((page - 1) * page_size + i),
        image:
          'https://s3.ap-northeast-1.amazonaws.com/static.readon.me/topic_img/wallet.svg',
        desc: 'a smart wallet to web3',
        users: 18,
        conversations: 721,
        created_at: {
          seconds: 1724435592,
        },
        creator: {
          name: 'Name',
          head_pic: 'https://static.readon.me/avatar/avatar3.png',
        },
        creator_id: 187,
        welcome_msg:
          'Welcome to your Web3 Agent! 🚀\n\nSeamlessly connect your wallet and experience the power of decentralized transactions at your fingertips. Whether you’re swapping tokens, staking, or sending crypto, our agent is here to make your journey smooth and secure. Dive into the future of finance with ease and confidence—your crypto transactions are just a click away!\n\nYou can start by connecting your wallet or exploring what I can do for you.',
        linked_plugin: range(0, 5).map(() => ({
          name: 'Slack',
          image:
            'https://s3.ap-northeast-1.amazonaws.com/static.readon.me/topic_img/wallet.svg',
        })) as Maybe<{name: string; image: string}[]>,
        tab: 'Social',
      })),
    },
    message: 'success',
  }),
})

export const useAgentDetailQuery = defineQuery({
  method: 'post',
  url: '/v1/studio/bot_detail',
  mockParams: {
    bot_id: 1,
  },
  mockData: {
    code: 0,
    data: {
      id: 1,
      name: 'Smart Wallet',
      status: 1,
      image: 'https://static.readon.me/logo/logo.webp',
      desc: 'a smart wallet to web3',
      users: 14,
      conversations: 567,
      created_at: {
        seconds: 1724435592,
      },
      linked_plugin: [
        {
          name: 'REVOX Basics',
          image:
            'https://s3.ap-northeast-1.amazonaws.com/static.readon.me/topic_img/REVOX.jpeg',
        },
        {
          name: 'Metamask Connector',
          image:
            'https://s3.ap-northeast-1.amazonaws.com/static.readon.me/topic_img/metamask.png',
        },
        {
          name: 'Pancake Connector',
          image:
            'https://s3.ap-northeast-1.amazonaws.com/static.readon.me/topic_img/pancakeswap.jpeg',
        },
      ] as Maybe<{name: string; image: string}[]>,
      work_modes: [
        {
          name: 'Single agent',
        },
      ],
      configuration: [
        {
          name: 'Defi Knowledge Base',
        },
        {
          name: 'Swap Workflow',
        },
      ],
      creator: {
        name: 'revox_official',
        head_pic:
          'https://s3.ap-northeast-1.amazonaws.com/static.readon.me/topic_img/REVOX.jpeg',
      },
      creator_id: 187,
      welcome_msg:
        'Welcome to your Web3 Agent! ?\n\nSeamlessly connect your wallet and experience the power of decentralized transactions at your fingertips. Whether you’re swapping tokens, staking, or sending crypto, our agent is here to make your journey smooth and secure. Dive into the future of finance with ease and confidence—your crypto transactions are just a click away!\n\nYou can start by connecting your wallet or exploring what I can do for you.',
      guide_info: [
        'Connect my wallet with Metamask',
        'Switch to BSC',
        'I want to buy USDT with my BNB',
      ] as Maybe<string[]>,
      tab: 'Social',
    },
    message: 'success',
  },
})

export const {useRequest: usePluginListQuery, request: pluginListQuery} =
  defineApi({
    method: 'post',
    url: '/v1/studio/plugin_list',
    mockParams: {} as Partial<{
      tab: string
      page: number
      page_size: number
      keywords: string
      ids: number[]
    }>,
    mockData: ({page = 1, page_size = 100}) => ({
      code: 0,
      data: {
        list: range(0, page_size).map((i) => ({
          id: 2,
          name: 'Google Search' + ((page - 1) * page_size + i),
          desc: 'a plugin that can search on google.',
          image:
            'https://s3.ap-northeast-1.amazonaws.com/static.readon.me/topic_img/REVOX.jpeg',
          linked_agent: range(0, 5).map(() => ({
            name: 'Slack',
            image:
              'https://s3.ap-northeast-1.amazonaws.com/static.readon.me/topic_img/wallet.svg',
          })) as Maybe<{name: string; image: string}[]>,
          creator: {
            name: 'Name',
            head_pic: 'https://static.readon.me/avatar/avatar3.png',
          },
          creator_id: 187,
          created_at: {
            seconds: 1724435592,
          },
          methods: [method] as (typeof method)[] | undefined,
          depend_plugins: [1, 3],
        })),
        message: 'success',
      },
    }),
  })

export const {data: pluginList} = usePluginListQuery()

const method = {
  id: 16,
  plugin_id: 4,
  name: 'search',
  description: 'Perform a search operation',
  http_sub_path: '/search',
  http_method: 'POST',
  input_schema:
    '{"parameters": [{"name": "query", "type": "string", "visible": true, "required": true, "description": "The search query", "input_method": "body", "default_value": null}, {"name": "max_results", "type": "integer", "visible": true, "required": true, "description": "Maximum number of results to return", "input_method": "body", "default_value": null}, {"enum": ["basic"], "name": "search_depth", "type": "string", "visible": true, "required": true, "description": "Depth of the search", "input_method": "body"}]}',
  input_example:
    '{\n    "query": "Who is Leo Messi?",\n    "search_depth": "basic",\n    "max_results": 5\n}',
  output_schema:
    '[{"name": "query", "type": "string", "description": "The search query submitted by the user"}, {"name": "follow_up_questions", "type": "array", "items": {"type": "string"}, "description": "Potential follow-up questions related to the query", "default_value": null}, {"name": "answer", "type": "string", "description": "A direct answer to the query, if available", "default_value": null}, {"name": "images", "type": "array", "items": {"type": "string"}, "description": "Any images related to the query", "default_value": []}, {"name": "results", "type": "array", "items": {"type": "object", "properties": [{"name": "title", "type": "string", "description": "Title of the search result"}, {"name": "url", "type": "string", "description": "URL of the search result"}, {"name": "content", "type": "string", "description": "Snippet or summary of the search result content"}, {"name": "score", "type": "number", "description": "Relevance score of the search result"}, {"name": "raw_content", "type": "string", "description": "Raw content of the search result, if available", "default_value": null}]}, "description": "Search results related to the query"}, {"name": "response_time", "type": "number", "description": "Time taken to process the query and generate results"}]',
  output_example:
    '{\n    "query": "Who is Leo Messi?",\n    "follow_up_questions": null,\n    "answer": null,\n    "images": [],\n    "results": [\n        {\n            "title": "Lionel Messi | Biography, Barcelona, PSG, Ballon d\'Or, Inter Miami ...",\n            "url": "https://www.britannica.com/biography/Lionel-Messi",\n            "content": "In early 2009 Messi capped off a spectacular 2008–09 season by helping FC Barcelona capture the club’s first “treble” (winning three major European club titles in one season): the team won the La Liga championship, the Copa del Rey (Spain’s major domestic cup), and the Champions League title. Messi’s play continued to rapidly improve over the years, and by 2008 he was one of the most dominant players in the world, finishing second to Manchester United’s Cristiano Ronaldo in the voting for the 2008 Ballon d’Or. At the 2014 World Cup, Messi put on a dazzling display, scoring four goals and almost single-handedly propelling an offense-deficient Argentina team through the group stage and into the knockout rounds, where Argentina then advanced to the World Cup final for the first time in 24 years. After Argentina was defeated in the Copa final—the team’s third consecutive finals loss in a major tournament—Messi said that he was quitting the national team, but his short-lived “retirement” lasted less than two months before he announced his return to the Argentine team. Messi helped Barcelona capture another treble during the 2014–15 season, leading the team with 43 goals scored over the course of the campaign, which resulted in his fifth world player of the year honour.",\n            "score": 0.98564,\n            "raw_content": null\n        },\n        {\n            "title": "Lionel Messi and the unmistakeable sense of an ending",\n            "url": "https://www.nytimes.com/athletic/5637953/2024/07/15/lionel-messi-argentina-ending-injury/",\n            "content": "First, he sank to the ground, grimacing. Play continued for a few seconds and then came the communal gasp. Lionel Messi was down. And Lionel Messi is not a player who goes down for nothing ...",\n            "score": 0.98369,\n            "raw_content": null\n        },\n        {\n            "title": "Lionel Messi: Biography, Soccer Player, Inter Miami CF, Athlete",\n            "url": "https://www.biography.com/athletes/lionel-messi",\n            "content": "The following year, after Messi heavily criticized the referees in the wake of a 2-0 loss to Brazil in the Copa America semifinals, the Argentine captain was slapped with a three-game ban by the South American Football Confederation.\\n So, at the age of 13, when Messi was offered the chance to train at soccer powerhouse FC Barcelona’s youth academy, La Masia, and have his medical bills covered by the team, Messi’s family picked up and moved across the Atlantic to make a new home in Spain. Famous Athletes\\nDennis Rodman\\nBrett Favre\\nTiger Woods\\nJohn McEnroe\\nKurt Warner\\nSandy Koufax\\n10 Things You Might Not Know About Travis Kelce\\nPeyton Manning\\nJames Harden\\nKobe Bryant\\nStephen Curry\\nKyrie Irving\\nA Part of Hearst Digital Media\\n Their marriage, a civil ceremony dubbed by Argentina’s Clarín newspaper as the “wedding of the century,” was held at a luxury hotel in Rosario, with a number of fellow star soccer players and Colombian pop star Shakira on the 260-person guest list.\\n In 2013, the soccer great came back to earth somewhat due to the persistence of hamstring injuries, but he regained his record-breaking form by becoming the all-time leading scorer in La Liga and Champions League play in late 2014.\\n",\n            "score": 0.97953,\n            "raw_content": null\n        },\n        {\n            "title": "Lionel Messi - Wikipedia",\n            "url": "https://en.wikipedia.org/wiki/Lionel_Messi",\n            "content": "He scored twice in the last group match, a 3–2 victory over Nigeria, his second goal coming from a free kick, as they finished first in their group.[423] Messi assisted a late goal in extra time to ensure a 1–0 win against Switzerland in the round of 16, and played in the 1–0 quarter-final win against Belgium as Argentina progressed to the semi-final of the World Cup for the first time since 1990.[424][425] Following a 0–0 draw in extra time, they eliminated the Netherlands 4–2 in a penalty shootout to reach the final, with Messi scoring his team\'s first penalty.[426]\\nBilled as Messi versus Germany, the world\'s best player against the best team, the final was a repeat of the 1990 final featuring Diego Maradona.[427] Within the first half-hour, Messi had started the play that led to a goal, but it was ruled offside. \\"[582] Moreover, several pundits and footballing figures, including Maradona, questioned Messi\'s leadership with Argentina at times, despite his playing ability.[583][584][585] Vickery states the perception of Messi among Argentines changed in 2019, with Messi making a conscious effort to become \\"more one of the group, more Argentine\\", with Vickery adding that following the World Cup victory in 2022 Messi would now be held in the same esteem by his compatriots as Maradona.[581]\\nComparisons with Cristiano Ronaldo\\nAmong his contemporary peers, Messi is most often compared and contrasted with Portuguese forward Cristiano Ronaldo, as part of an ongoing rivalry that has been compared to past sports rivalries like the Muhammad Ali–Joe Frazier rivalry in boxing, the Roger Federer–Rafael Nadal rivalry in tennis, and the Prost–Senna rivalry from Formula One motor racing.[586][587]\\nAlthough Messi has at times denied any rivalry,[588][589] they are widely believed to push one another in their aim to be the best player in the world.[160] Since 2008, Messi has won eight Ballons d\'Or to Ronaldo\'s five,[590] seven FIFA World\'s Best Player awards to Ronaldo\'s five, and six European Golden Shoes to Ronaldo\'s four.[591] Pundits and fans regularly argue the individual merits of both players.[160][592] On 11 July, Messi provided his 20th assist of the league season for Arturo Vidal in a 1–0 away win over Real Valladolid, equalling Xavi\'s record of 20 assists in a single La Liga season from 2008 to 2009;[281][282] with 22 goals, he also became only the second player ever, after Thierry Henry in the 2002–03 FA Premier League season with Arsenal (24 goals and 20 assists), to record at least 20 goals and 20 assists in a single league season in one of Europe\'s top-five leagues.[282][283] Following his brace in a 5–0 away win against Alavés in the final match of the season on 20 May, Messi finished the season as both the top scorer and top assist provider in La Liga, with 25 goals and 21 assists respectively, which saw him win his record seventh Pichichi trophy, overtaking Zarra; however, Barcelona missed out on the league title to Real Madrid.[284] On 7 March, two weeks after scoring four goals in a league fixture against Valencia, he scored five times in a Champions League last 16-round match against Bayer Leverkusen, an unprecedented achievement in the history of the competition.[126][127] In addition to being the joint top assist provider with five assists, this feat made him top scorer with 14 goals, tying José Altafini\'s record from the 1962–63 season, as well as becoming only the second player after Gerd Müller to be top scorer in four campaigns.[128][129] Two weeks later, on 20 March, Messi became the top goalscorer in Barcelona\'s history at 24 years old, overtaking the 57-year record of César Rodríguez\'s 232 goals with a hat-trick against Granada.[130]\\nDespite Messi\'s individual form, Barcelona\'s four-year cycle of success under Guardiola – one of the greatest eras in the club\'s history – drew to an end.[131] He still managed to break two longstanding records in a span of seven days: a hat-trick on 16 March against Osasuna saw him overtake Paulino Alcántara\'s 369 goals to become Barcelona\'s top goalscorer in all competitions including friendlies, while another hat-trick against Real Madrid on 23 March made him the all-time top scorer in El Clásico, ahead of the 18 goals scored by former Real Madrid player Alfredo Di Stéfano.[160][162] Messi finished the campaign with his worst output in five seasons, though he still managed to score 41 goals in all competitions.[161][163] For the first time in five years, Barcelona ended the season without a major trophy; they were defeated in the Copa del Rey final by Real Madrid and lost the league in the last game to Atlético Madrid, causing Messi to be booed by sections of fans at the Camp Nou.[164]",\n            "score": 0.97579,\n            "raw_content": null\n        },\n        {\n            "title": "The life and times of Lionel Messi",\n            "url": "https://www.nytimes.com/athletic/4783674/2023/08/18/lionel-messi-profile-soccer/",\n            "content": "For Messi, it is major trophy number 44.. Despite turning 36 in June, he is as influential as ever. Here is the complete story of Lionel Andres Messi, widely regarded as one of the greatest ...",\n            "score": 0.96961,\n            "raw_content": null\n        }\n    ],\n    "response_time": 0.88\n}',
  status: 1,
}

export const usePluginDetailQuery = defineQuery({
  method: 'post',
  url: '/v1/studio/plugin_detail',
  mockParams: {
    id: 1,
  },
  mockData: {
    code: 0,
    data: {
      plugin: {
        id: 1,
        name: 'Wallet Base',
        desc: 'a base wallet plugins, support basic wallet method.',
        image:
          'https://s3.ap-northeast-1.amazonaws.com/static.readon.me/topic_img/REVOX.jpeg',
        creator: {
          name: 'Name',
          head_pic: 'https://static.readon.me/avatar/avatar3.png',
        },
        created_at: {
          seconds: 1724435592,
        },
        creator_id: 187,
        linked_agent: [
          {
            name: 'smart_wallet',
            image: 'https://static.readon.me/logo/logo.webp',
          },
        ] as Maybe<{name: string; image: string}[]>,
        methods: [method] as (typeof method)[] | undefined,
        depend_plugins: [1, 3],
      },
    },
    message: 'success',
  },
})

export const usePluginTabsQuery = defineQuery({
  url: '/v1/studio/plugin_tabs',
  mockData: {
    code: 0,
    data: [
      'Most used',
      'Efficiency Tools',
      'Web3',
      'Entertainment',
      'Game',
      'Lifestyle',
      'Education',
    ],
    message: 'success',
  },
})

export const useAgentTabsQuery = defineQuery({
  url: '/v1/studio/bot_tabs',
  mockData: {
    code: 0,
    data: [
      'Most used',
      'News',
      'Web3',
      'Web Search',
      'Social',
      'Education',
      'Finance',
    ],
    message: 'success',
  },
})

export const useUserAgentList = defineQuery({
  url: '/v1/studio/user_bot_list',
  method: 'post',
  mockParams: {
    page: 1,
    page_size: 10,
  },
  mockData: {
    code: 0,
    data: {
      bots: [
        {
          id: 1,
          name: 'smart_wallet',
          image:
            'https://s3.ap-northeast-1.amazonaws.com/static.readon.me/topic_img/wallet.svg',
          desc: 'a smart wallet to web3',
          status: 1,
          created_at: {
            seconds: 1724435592,
          },
          linked_plugin: [
            {
              id: 1,
              name: 'Wallet Base',
              desc: 'a base wallet plugins, support basic wallet method.',
              image:
                'https://s3.ap-northeast-1.amazonaws.com/static.readon.me/topic_img/REVOX.jpeg',
            },
          ],
        },
        {
          id: 2,
          name: 'smart_wallet2',
          image:
            'https://s3.ap-northeast-1.amazonaws.com/static.readon.me/topic_img/wallet.svg',
          desc: 'a smart wallet to web3',
          status: 1,
          created_at: {
            seconds: 1724435592,
          },
          linked_plugin: [
            {
              id: 1,
              name: 'Wallet Base',
              desc: 'a base wallet plugins, support basic wallet method.',
              image:
                'https://s3.ap-northeast-1.amazonaws.com/static.readon.me/topic_img/REVOX.jpeg',
            },
          ],
        },
      ],
    },
    message: 'success',
  },
})

export const useAgentDraftQuery = defineQuery({
  method: 'post',
  url: '/v1/studio/get_bot_draft',
  mockParams: {
    bot_id: 1,
  },
  mockData: {
    code: 0,
    data: {
      draft: {
        id: 3,
        bot_id: 12,
        creator_id: 1524,
        prompt: 'draft string',
        plugins: [1, 2, 3, 4],
        welcome_msg: 'draft welcone string',
        guide_info: ['draft 1 string', 'draft 2 string'],
        debug_conversation_id: 1206,
        created_at: {
          seconds: 1726819524,
        },
        updated_at: {
          seconds: 1726822571,
        },
        model_settings: {
          model: 'gpt-4o',
          temperature: 1,
          top_p: 1,
          rounds: 10,
          max_length: 4095,
        },
      },
    },
    message: 'success',
  },
})

export const createAgentMutation = defineMutation({
  url: '/v1/studio/create_bot',
  mockParams: {
    name: 'string',
  },
  mockData: {
    code: 0,
    message: 'success',
    data: {
      id: 42,
    },
  },
})

export const updateAgentMutation = defineMutation({
  url: '/v1/studio/update_bot_info',
  mockParams: {
    bot_id: 1,
    name: 'string',
    description: 'string',
    type: 'string',
  },
  mockData: {
    code: 0,
    message: 'success',
    data: {
      id: 42,
    },
  },
})

export const publishAgentMutation = defineMutation({
  url: '/v1/studio/publish_bot',
  mockParams: {
    bot_id: 1,
    prompt: 'string',
    plugins: [1, 2, 3],
    welcome_msg: 'string',
    guide_info: ['string', 'string'],
    model_settings: {
      model: 'gpt-4o',
      temperature: 1,
      top_p: 1,
      rounds: 10,
      max_length: 4095,
    },
  },
  mockData: {
    code: 0,
    data: {
      id: 1,
    },
    message: 'success',
  },
})

export const saveAgentDraftMutation = defineMutation({
  url: '/v1/studio/save_bot_draft',
  mockParams: {
    bot_id: 1,
    prompt: 'string',
    plugins: [1, 2, 3],
    welcome_msg: 'string',
    guide_info: ['string', 'string'],
    model_settings: {
      model: 'gpt-4o',
      temperature: 1,
      top_p: 1,
      rounds: 10,
      max_length: 4095,
    },
  },
  mockData: {
    code: 0,
    data: {
      draft_id: 1,
    },
    message: 'success',
  },
})
