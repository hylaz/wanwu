import concurrent.futures
from pathlib import Path
from typing import List, NamedTuple, Optional, Union, cast

class FileEncodeType(NamedTuple):
    encoding: Optional[str]
    confidence: float
    language: Optional[str]


def detect_encodings(
    file_path: Union[str, Path], timeout: int = 5
) -> List[FileEncodeType]:

    import chardet

    file_path = str(file_path)

    def read_and_detect(file_path: str) -> List[dict]:
        with open(file_path, "rb") as f:
            rawdata = f.read()
        return cast(List[dict], chardet.detect_all(rawdata))

    with concurrent.futures.ThreadPoolExecutor() as executor:
        future = executor.submit(read_and_detect, file_path)
        try:
            encodings = future.result(timeout=timeout)
        except concurrent.futures.TimeoutError:
            raise TimeoutError(
                f"Timeout error {file_path}"
            )

    if all(encoding["encoding"] is None for encoding in encodings):
        raise RuntimeError(f"can not detect encoding for {file_path}")
    return [FileEncodeType(**enc) for enc in encodings if enc["encoding"] is not None]