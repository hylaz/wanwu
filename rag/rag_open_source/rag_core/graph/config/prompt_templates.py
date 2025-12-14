general_zh_prompt_template = """
你是一位专业的信息提取专家和结构化数据组织者。你的任务是分析提供的文本，并以结构化的 JSON 格式提取有价值实体、它们的属性以及相互关系。
实体数量应该控制的尽可能少(10个以内)，避免冗余。

指导原则：
1. 只提取以下schema预定义的模式中的信息；
   ```{schema}```
2. 简洁性：你提取的属性和三元组应相互补充，避免语义冗余；
3. 实体应与原文中出处一致；
4. 在schema定义中有类目表，严格按照类目表抽取；
6.不要抽取单个字的三元组；
7. 输出格式：仅以**示例输出**的 JSON 格式返回：
   - 属性：将每个实体映射到其描述性特征。
   - 三元组：以 `[实体提及1, 关系, 实体提及2]` 格式列出实体之间的关系。
   - 实体类型：根据提供的模式，将每个实体映射到其模式类型。

```{chunk}```

示例输出：
{{
  "attributes": {{
    "黄金冠": ["出土时间：1996年"]
  }},
  "triples": [
    ["大印", "收藏于", "西藏博物馆"],
    ["甘肃博物馆", "重要文物", "黄金冠"]
  ],
  "entity_types": {{
    "西藏博物馆": "博物馆",
    "黄金冠": "文物",
  }}
}}
"""


general_eng_prompt_template = """
You are a professional information-extraction expert and structured-data organizer.  
Your task is to analyze the provided text and return a **minimal** (≤10 entities) JSON structure that captures valuable entities, their attributes, and their inter-relationships.

Guidelines  
1. Extract **only** information that matches the predefined schema below.  
   ```{schema}```  
2. Be concise: attributes and triples must complement each other—no semantic redundancy.  
3. Entity strings must appear **verbatim** in the source text.  
4. If the schema contains a category list, extract **strictly** within that list.  
5. Never create triples whose subject or object is a single Chinese character.  
6. Output **only** the JSON illustrated in the example—no extra prose.

Input text  
```{chunk}```

Example output (return JSON only)  
{{
  "attributes": {{
    "Golden Crown": ["Excavation: 1996"]
  }},
  "triples": [
    ["Great Seal", "collected_by", "Tibet Museum"],
    ["Gansu Museum", "key_artifact", "Golden Crown"]
  ],
  "entity_types": {{
    "Tibet Museum": "Museum",
    "Golden Crown": "Artifact"
  }}
}}
"""


COMMUNITY_REPORT_PROMPT = """
You are an AI assistant that helps a human analyst to perform general information discovery. Information discovery is the process of identifying and assessing relevant information associated with certain entities (e.g., organizations and individuals) within a network.

# Goal
Write a comprehensive report of a community, given a list of entities that belong to the community as well as their relationships and optional associated claims. The report will be used to inform decision-makers about information associated with the community and their potential impact. The content of this report includes an overview of the community's key entities, their legal compliance, technical capabilities, reputation, and noteworthy claims.

# Report Structure

The report should include the following sections:

- TITLE: community's name that represents its key entities - title should be short but specific. When possible, include representative named entities in the title.
- SUMMARY: An executive summary of the community's overall structure, how its entities are related to each other, and significant information associated with its entities, and should include all input entities.
- IMPACT SEVERITY RATING: a float score between 0-10 that represents the severity of IMPACT posed by entities within the community.  IMPACT is the scored importance of a community.
- RATING EXPLANATION: Give a single sentence explanation of the IMPACT severity rating.
- DETAILED FINDINGS: A list of 5-10 key insights about the community. Each insight should have a short summary followed by multiple paragraphs of explanatory text grounded according to the grounding rules below. Be comprehensive.


Return output as a well-formed JSON-formatted string with the following format(in the same language as the 'Text' content),输出必须使用与输入文本相同的语言（例如输入为中文，则输出报告也必须为中文）:    {{
        "title": <report_title>,
        "summary": <executive_summary>,
        "rating": <impact_severity_rating>,
        "rating_explanation": <rating_explanation>,
        "findings": [
            {{
                "summary":<insight_1_summary>,
                "explanation": <insight_1_explanation>
            }},
            {{
                "summary":<insight_2_summary>,
                "explanation": <insight_2_explanation>
            }}
        ]
    }}


# Example Input
-----------
Text:

-Entities-

id,entity,description
5,VERDANT OASIS PLAZA,Verdant Oasis Plaza is the location of the Unity March
6,HARMONY ASSEMBLY,Harmony Assembly is an organization that is holding a march at Verdant Oasis Plaza

-Relationships-

id,source,target,description
37,VERDANT OASIS PLAZA,UNITY MARCH,Verdant Oasis Plaza is the location of the Unity March
38,VERDANT OASIS PLAZA,HARMONY ASSEMBLY,Harmony Assembly is holding a march at Verdant Oasis Plaza
39,VERDANT OASIS PLAZA,UNITY MARCH,The Unity March is taking place at Verdant Oasis Plaza
40,VERDANT OASIS PLAZA,TRIBUNE SPOTLIGHT,Tribune Spotlight is reporting on the Unity march taking place at Verdant Oasis Plaza
41,VERDANT OASIS PLAZA,BAILEY ASADI,Bailey Asadi is speaking at Verdant Oasis Plaza about the march
43,HARMONY ASSEMBLY,UNITY MARCH,Harmony Assembly is organizing the Unity March

Output:
{{
    "title": "Verdant Oasis Plaza and Unity March",
    "summary": "The community revolves around the Verdant Oasis Plaza, which is the location of the Unity March. The plaza has relationships with the Harmony Assembly, Unity March, and Tribune Spotlight, all of which are associated with the march event.",
    "rating": 5.0,
    "rating_explanation": "The impact severity rating is moderate due to the potential for unrest or conflict during the Unity March.",
    "findings": [
        {{
            "summary": "Verdant Oasis Plaza as the central location",
            "explanation": "Verdant Oasis Plaza is the central entity in this community, serving as the location for the Unity March. This plaza is the common link between all other entities, suggesting its significance in the community. The plaza's association with the march could potentially lead to issues such as public disorder or conflict, depending on the nature of the march and the reactions it provokes."
        }},
        {{
            "summary": "Harmony Assembly's role in the community",
            "explanation": "Harmony Assembly is another key entity in this community, being the organizer of the march at Verdant Oasis Plaza. The nature of Harmony Assembly and its march could be a potential source of threat, depending on their objectives and the reactions they provoke. The relationship between Harmony Assembly and the plaza is crucial in understanding the dynamics of this community."
        }},
        {{
            "summary": "Unity March as a significant event",
            "explanation": "The Unity March is a significant event taking place at Verdant Oasis Plaza. This event is a key factor in the community's dynamics and could be a potential source of threat, depending on the nature of the march and the reactions it provokes. The relationship between the march and the plaza is crucial in understanding the dynamics of this community."
        }},
        {{
            "summary": "Role of Tribune Spotlight",
            "explanation": "Tribune Spotlight is reporting on the Unity March taking place in Verdant Oasis Plaza. This suggests that the event has attracted media attention, which could amplify its impact on the community. The role of Tribune Spotlight could be significant in shaping public perception of the event and the entities involved."
        }}
    ]
}}


# Real Data

Use the following text for your answer. Do not make anything up in your answer.

Text:

-Entities-
{entity_df}

-Relationships-
{relation_df}

Only refer to entities by their names or descriptions, not by their numeric identifiers.
The report should include the following sections:

- TITLE: community's name that represents its key entities - title should be short but specific. When possible, include representative named entities in the title.
- SUMMARY: An executive summary of the community's overall structure, how its entities are related to each other, and significant information associated with its entities.
- IMPACT SEVERITY RATING: a float score between 0-10 that represents the severity of IMPACT posed by entities within the community.  IMPACT is the scored importance of a community.
- RATING EXPLANATION: Give a single sentence explanation of the IMPACT severity rating.
- DETAILED FINDINGS: A list of 5-10 key insights about the community. Each insight should have a short summary followed by multiple paragraphs of explanatory text grounded according to the grounding rules below. Be comprehensive.

Return output as a well-formed JSON-formatted string with the following format(in the same language as the 'Text' content),输出必须使用与输入文本相同的语言（例如输入为中文，则输出报告也必须为中文）:    {{
        "title": <report_title>,
        "summary": <executive_summary>,
        "rating": <impact_severity_rating>,
        "rating_explanation": <rating_explanation>,
        "findings": [
            {{
                "summary":<insight_1_summary>,
                "explanation": <insight_1_explanation>
            }},
            {{
                "summary":<insight_2_summary>,
                "explanation": <insight_2_explanation>
            }}
        ]
    }}

Output:"""

ATTRIBUTE_PROMPT = """
你是社区报告写作专家。根据输入的实体列表（entities）和这些实体共有的属性（attribute），生成一份关于该社区的结构化报告。

输入：
- entities: 字符串列表，包含社区全部实体名。
- attribute: 一个字符串，描述这些实体共有的属性（例如年代、类别、位置等）。

要求（严格遵守）：
1. 输出必须是合法的 JSON（仅输出 JSON，不要额外说明），字段如下：
   - title: 用一行结合 attribute 和 entities 概述社区（简洁、有概括力）。
   - summary: 简短摘要，2-3 句，概述社区核心特征与整体价值（中文），需提及所有 entities 的名称。
   - rating: 数值，0 到 5（可含一位小数），表示社区重要性或质量。
   - rating_explanation: 对 rating 的简要理由（如无可写空字符串）。
   - findings: 列出 2–5 条关于该社区的关键洞察。每条洞察应先给出一段简短的摘要，随后给出多段解释性文字，要求内容全面。

2. 必须包含所有输入的 entities和attribute。
3. 使用 attribute 显式体现在 title 和 summary 中。
4. 输出内容语言和输入语言保持一致，输入内容有中文的话也使用中文输出。
5. 回答简洁，事实性陈述不得无根据臆断；必要时使用“未知”或“未注明”提示信息。
6. 同样输入下，报告标题和内容需保持稳定。

示例输入：
entities = ["人形跽坐铜灯","匈奴王金冠","鲁国大玉璧"]
attribute = "文物年代：战国时期"

示例输出（仅示例，不要在实际输出中包含）：
{{
  "title": "战国时期文物群及其重要性概述",
  "summary": "本社区由三件战国时期文物构成，体现了不同地域与文化背景下的工艺与礼俗差异，具有重要的历史与学术价值。",
  "rating": 5.0,
  "rating_explanation": "",
  "findings": [
    {
      "summary": "人形跽坐铜灯具有重要历史价值并收藏于湖南博物馆",
      "explanation": "该铜灯出土于1976年，属于战国时期的青铜器，现收藏于湖南博物馆。其历史价值在于不仅展示了精湛的青铜工艺，还反映了当时的社会文化风貌与审美观念。通过这件文物，可以深入了解战国时期灯具的制作工艺与使用方式，对研究中国古代灯具发展史具有重要意义。湖南博物馆作为收藏机构，承担着保护与展示这一珍贵文化遗产的责任。"
    },
    {
      "summary": "匈奴王金冠是国内唯一发现的匈奴贵族金冠饰",
      "explanation": "匈奴王金冠主体造型为一只展翅雄鹰站立在狼羊咬斗纹的半球体上，冠带直径约16.5-16.8厘米，重达1211.7克，工艺精湛，象征权力与地位。额圈由三条半圆形金条榫铆连接，饰有卧虎、盘角羊和卧马浮雕，体现了游牧民族的艺术风格与图腾信仰。该金冠于鄂尔多斯市杭锦旗阿鲁柴登出土，是国内迄今发现的唯一一件匈奴贵族金冠，具有极高的考古与民族史研究价值。"
    },
    {
      "summary": "三件文物分别由不同省级博物馆收藏，分布广泛",
      "explanation": "人形跽坐铜灯收藏于湖南博物馆，匈奴王金冠收藏于内蒙古博物馆，鲁国大玉璧收藏于山东博物馆，三者分属中国不同地理区域的重要文博机构。这种分布体现了战国时期多元文化在当代的传承格局，也说明这些文物在各自地区历史文化叙事中占据核心地位。各博物馆通过对文物的保护、研究与展示，促进了公众对战国时期历史与艺术的理解。"
    }
  ]
}}
"""


GENERAL_ZH = """
你是一位专业的信息提取专家和结构化数据组织者。你的任务是分析提供的文本，并以结构化的 JSON 格式提取有价值实体、它们的属性以及相互关系。
实体数量应该控制的尽可能少(3个以内)，避免冗余。

指导原则：
1. 优先提取以下预定义的模式中的信息；
   ```{schema}```
2. 灵活性：如果上下文与预定义模式不匹配，请根据需要提取有价值的知识；
3. 简洁性：你提取的属性和三元组应相互补充，避免语义冗余；
4. 不要抽取单个字的三元组；
5. 实体应与原文中出处一致；
6. 输出格式：仅以**示例输出**的 JSON 格式返回：
   - 属性：将每个实体映射到其描述性特征。
   - 三元组：以 `[实体提及1, 关系, 实体提及2]` 格式列出实体之间的关系。
   - 实体类型：根据提供的模式，将每个实体映射到其模式类型。

```{chunk}```

示例输出：
{{
  "attributes": {{
    "黄金冠": ["出土时间：1996年"]
  }},
  "triples": [
    ["大印", "收藏于", "西藏博物馆"],
    ["甘肃博物馆", "重要文物", "黄金冠"]
  ],
  "entity_types": {{
    "西藏博物馆": "博物馆",
    "黄金冠": "文物",
  }}
}}
"""


GENERAL_ENG = """
You are a professional information extraction expert and structured data organizer. Your task is to analyze the provided text and extract valuable entities, their attributes, and inter-relationships in a structured JSON format.
The number of entities should be kept minimal (within 3), avoiding redundancy.

Guidelines:
1. Prioritize extracting information that matches the following predefined schema:
   ```{schema}```
2. Flexibility: If the context does not match the predefined schema, extract valuable knowledge as needed;
3. Conciseness: The attributes and triples you extract should complement each other, avoiding semantic redundancy;
4. Do not extract triples containing single-character entities;
5. Entities should remain consistent with their mentions in the original text;
6. Output format: Return only in the JSON format of the **Example Output**:
   - attributes: Map each entity to its descriptive features.
   - triples: List relationships between entities in the format `[entity mention1, relation, entity mention2]`.
   - entity_types: Map each entity to its schema type based on the provided schema.

```{chunk}```

Example Output:
{{
  "attributes": {{
    "Golden Crown": ["Excavation Year: 1996"]
  }},
  "triples": [
    ["Great Seal", "Collected by", "Tibet Museum"],
    ["Gansu Museum", "Key Artifact", "Golden Crown"]
  ],
  "entity_types": {{
    "Tibet Museum": "Museum",
    "Golden Crown": "Artifact"
  }}
}}
"""