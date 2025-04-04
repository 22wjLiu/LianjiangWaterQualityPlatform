import request from "@/util/request";
export const forecast = (params) => {
  return request.get(`/forecast?${params}`, {
    needToken: true
  });
};
