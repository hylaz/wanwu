import requests
import json
import time
from datetime import datetime
from settings import DOC_STATUS_INIT_URL


if __name__ == '__main__':
    res = requests.get(DOC_STATUS_INIT_URL)
    print(f"======={datetime.now().strftime('%Y-%m-%d %H:%M:%S')} DOC_STATUS_INIT_URL 请求结果:=======")
    print(res.status_code)
    result_data = json.loads(res.text)
    print(result_data)
