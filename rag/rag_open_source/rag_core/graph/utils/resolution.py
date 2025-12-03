import networkx as nx
import editdistance
import re
from typing import List, Dict, Tuple, Set
from graph.config import get_config
from graph.utils import call_llm_api


class LLMEntityResolver:
    def __init__(self, config=None):
        if config is None:
            config = get_config()
        """Init method definition."""
        self._llm_client = call_llm_api.LLMCompletionCall(config.construction.LLM_MODEL,
                                                          config.construction.LLM_BASE_URL,
                                                          config.construction.LLM_API_KEY)

    def resolve_by_name_and_llm(
            self,
            old_graph: nx.Graph,
            new_graph: nx.Graph
    ) -> Dict[str, str]:
        """
        结合名称相似性和LLM判断进行实体解析
        返回映射关系: {new_entity_name: old_entity_name}
        """
        # 1. 基于名称找到候选对
        candidate_pairs = self._find_similar_entities(old_graph, new_graph)

        # 2. 使用LLM辅助决策
        confirmed_mappings = self._llm_decision(candidate_pairs, old_graph, new_graph)

        return confirmed_mappings

    def _find_similar_entities(
            self,
            old_graph: nx.Graph,
            new_graph: nx.Graph
    ) -> List[Tuple[str, str]]:
        """基于名称相似性找到候选实体对"""
        temp_nodes = []
        candidates = []
        for new_node in new_graph.nodes():
            new_attrs = new_graph.nodes[new_node]
            new_type = new_attrs.get('label')
            new_schema_type = new_attrs['properties'].get('schema_type', '')

            for old_node in old_graph.nodes():
                old_attrs = old_graph.nodes[old_node]
                old_type = new_attrs.get('label')
                old_schema_type = new_attrs['properties'].get('schema_type', '')

                # 类型相同且名称相似
                if new_type == old_type and new_schema_type == old_schema_type and self._is_name_similar(new_node, old_node):
                    candidates.append((new_node, old_node))

        return candidates

    def _is_name_similar(self, name1: str, name2: str) -> bool:
        """判断两个名称是否相似"""
        # 英文使用编辑距离
        if self._is_english(name1) and self._is_english(name2):
            distance = editdistance.eval(name1.lower(), name2.lower())
            return distance <= min(len(name1), len(name2)) // 2

        # 中文使用字符重叠
        chars1, chars2 = set(name1), set(name2)
        overlap = len(chars1 & chars2)
        total = len(chars1 | chars2)
        return overlap / total > 0.6 if total > 0 else False

    def _is_english(self, text: str) -> bool:
        """判断是否为英文文本"""
        if not text:
            return False
        pattern = re.compile(r"[`a-zA-Z0-9\s.,':;/\"?<>!\(\)\-]")
        eng_count = sum(1 for char in text if pattern.fullmatch(char))
        return eng_count / len(text) > 0.8 if text else False

    def _llm_decision(
            self,
            candidate_pairs: List[Tuple[str, str]],
            old_graph: nx.Graph,
            new_graph: nx.Graph
    ) -> Dict[str, str]:
        """使用LLM辅助决策确认实体对是否相同"""
        if not candidate_pairs:
            return {}

        # 构造提示词
        prompt = self._build_llm_prompt(candidate_pairs, old_graph, new_graph)

        # 调用LLM
        response = self._llm_client.call_api(prompt)

        # 解析结果
        return self._parse_llm_response(response, candidate_pairs)

    def _build_llm_prompt(
            self,
            candidate_pairs: List[Tuple[str, str]],
            old_graph: nx.Graph,
            new_graph: nx.Graph
    ) -> str:
        """构建LLM提示词"""
        prompt_lines = [
            "请判断以下实体对或者属性对是否含义相同，",
            "示例输入：",
            "属性对 1: A属性:文物年代： 战国 \n B属性文物年代：战国时期",
            "属性对 2: A属性:出土时间： 1965 - 1966 年 \n B属性出土时间：1968年",
            "属性对 3: A属性:文物年份：唐代 \n B属性文物年份：明代",
            "属性对 4: A属性:出土时间：1959年 \n B属性出土时间：1965年",
            "实体对 5: A实体:甘肃省博物馆 类型：省博物馆 \n B实体湖北省博物馆 类型：省博物馆",
            "示例输出：",
            "pair 1: 是",
            "pair 2: 否",
            "pair 3: 否",
            "pair 4: 否",
            "pair 5: 否",
            "",
            "真实输入数据如下：",
            ""
        ]
        for idx, (new_name, old_name) in enumerate(candidate_pairs, 1):
            new_attrs = new_graph.nodes[new_name]['properties']
            old_attrs = old_graph.nodes[old_name]['properties']

            if new_graph.nodes[new_name].get('label') == "entity":
                prompt_lines.append(f"实体对 {idx}:")
                prompt_lines.append(f"  A实体: {new_name}")
                prompt_lines.append(f"    类型: {new_attrs.get('schema_type', '')}")
                prompt_lines.append(f"  B实体: {old_name}")
                prompt_lines.append(f"    类型: {old_attrs.get('schema_type', '')}")
                prompt_lines.append("")
            if new_graph.nodes[new_name].get('label') == "attribute":
                prompt_lines.append(f"属性对 {idx}:")
                prompt_lines.append(f"  A属性: {new_name}")
                prompt_lines.append(f"  B属性: {old_name}")
                prompt_lines.append("")
        prompt_lines.append("请回答每个实体/属性对是否相同，格式为:")
        prompt_lines.append("pair 1: 是/否")
        prompt_lines.append("pair 2: 是/否")
        prompt_lines.append("...")

        return "\n".join(prompt_lines)

    def _parse_llm_response(
            self,
            response: str,
            candidate_pairs: List[Tuple[str, str]]
    ) -> Dict[str, str]:
        """解析LLM响应结果"""
        mappings = {}
        lines = response.strip().split('\n')

        for line in lines:
            if ':' in line:
                parts = line.split(':')
                if len(parts) >= 2:
                    pair_info = parts[0].strip()
                    decision = parts[1].strip()

                    # 提取实体对编号
                    import re
                    match = re.search(r'pair (\d+)', pair_info)
                    if match and decision.startswith('是'):
                        idx = int(match.group(1)) - 1
                        if 0 <= idx < len(candidate_pairs):
                            new_name, old_name = candidate_pairs[idx]
                            if new_name != old_name:
                                mappings[new_name] = old_name

        return mappings

