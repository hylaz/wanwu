import service from "@/utils/request"
import {USER_API} from "@/utils/requestConstants"

export const setPlatformInfo = (type, data, config) => {
  return service({
    url: `${USER_API}/custom/${type}`,
    method: "post",
    data,
    config
  });
}