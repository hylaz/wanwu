import service from "@/utils/request"
import {USER_API} from "@/utils/requestConstants"

export const getHistoryList = () => {
  return service({
    url: `${USER_API}/exploration/app/history`,
    method: 'get',
  })
}
export const setFavorite = (data) => {
  return service({
    url: `${USER_API}/exploration/app/favorite`,
    method: 'post',
    data
  })
}
export const getExplorList = (params) => {
  return service({
    url: `${USER_API}/exploration/app/list`,
    method: 'get',
    params
  })
}