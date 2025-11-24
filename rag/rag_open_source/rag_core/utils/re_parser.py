import re


def extract_pattern(text, pattern):
    result = []
    matched_strs = re.findall(pattern, text)
    for match in matched_strs:
        matched_strs.append(match.strip())

    return result

if __name__ == "__main__":
    text = "测试使用发布日期：2025年08月08日， 测试2025-02-02"
    date_patterns = [
        # 匹配 "发布日期：2025年7月8日" 或 "发布日期: 2025年07月08日"
        r'发布日期[：:]\s*\d{4}\s*年\s*\d{1,2}\s*月\s*\d{1,2}\s*日',

        # 匹配 "发布日期：2025-07-08" 或 "发布日期: 2025-7-8"
        r'发布日期[：:]\s*\d{4}-\d{1,2}-\d{1,2}',

        # 匹配 "发布日期：2025/07/08" 或 "发布日期: 2025/7/8"
        r'发布日期[：:]\s*\d{4}/\d{1,2}/\d{1,2}',

        # 匹配 "发布日期：2025.07.08"
        r'发布日期[：:]\s*\d{4}\.\d{1,2}\.\d{1,2}',

        # 匹配独立的 "2025年7月8日" 或 "2025年07月08日"
        r'\d{4}\s*年\s*\d{1,2}\s*月\s*\d{1,2}\s*日',

        # 匹配独立的 "2025-07-08" 或 "2025-7-8"
        r'\b\d{4}-\d{1,2}-\d{1,2}\b',

        # 匹配独立的 "2025/07/08" 或 "2025/7/8"
        r'\b\d{4}/\d{1,2}/\d{1,2}\b',

        # 匹配独立的 "2025.07.08"
        r'\b\d{4}\.\d{1,2}\.\d{1,2}\b'
    ]

    # 按顺序遍历所有模式
    for pattern in date_patterns:
        matches = extract_pattern(text, pattern)
        for match in matches:
            print(match.strip())

