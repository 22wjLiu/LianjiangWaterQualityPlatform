import request from "@/util/request";

export const getLineData = (name, system, params) => {
  return request.get(`/data/${name}/${system}?${params}`, {
    needToken: true,
  });
};

export const getPieData = (name, system, field) => {
  return request.get(`/data/rowall/${name}/${system}/${field}`, {
    needToken: true,
  });
};

export const getNameList = (id) => {
  return request.get(`/map/${id}`, {
    needToken: true,
  });
};

export const getStationName = () => {
  return request.get("/stationName", {
    needToken: true,
  });
};

export const getTimeRange = (name, system) => {
  return request.get(`/timeRange/${name}/${system}`, {
    needToken: true,
  });
};

export const getActiveMapInfosByStationName = (mapType, stationName) => {
  return request.get(`/mapInfos/${mapType}/${stationName}`, {
    needToken: true,
  });
};
