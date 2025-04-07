import request from "@/util/request";

export const getMapData = () => {
  return request.get("/map", {
    needToken: true,
  });
};

export const getMapTables = () => {
  return request.get("/mapTables", {
    needToken: true,
  });
};

export const getCurMaps = (pramas) => {
  return request.get(`/curMaps${pramas}`, {
    needToken: true,
  });
};
