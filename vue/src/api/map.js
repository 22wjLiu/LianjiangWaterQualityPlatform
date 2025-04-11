import request from "@/util/request";

export const getMapTables = () => {
  return request.get("/mapTables", {
    needToken: true,
  });
};

export const getMapInfos = (id, pramas) => {
  let path = "/mapInfos";

  path += id ? `/${id}` : "/null";
  path += pramas;
  
  return request.get(path, {
    needToken: true,
  });
};

export const getMapVersions = (start, end, pramas) => {
  let path = "/mapVersions";

  path += start ? `/${start}` : "/null";
  path += end ? `/${end}` : "/null";
  path += pramas;

  return request.get(path, {
    needToken: true,
  });
};

export const createMapVersion = (pramas, body) => {
  return request.post(`/mapVersion${pramas}`, body, {
    needToken: true,
  });
};

export const createMap = (id, body) => {
  let path = "/createMap";

  path += id ? `/${id}` : "/null";

  return request.post(path, body, {
    needToken: true,
  });
};

export const deleteMap = (id, ids) => {
  let path = "/deleteMap";

  path += id ? `/${id}` : "/null";

  return request.delete(path, {
    data: { ids },
    needToken: true,
  });
};

export const updateMap = (id, curMapId, body) => {
  let path = "/updateMap";

  path += id ? `/${id}` : "/null";
  path += curMapId ? `/${curMapId}` : "/null";

  return request.put(path, body, {
    needToken: true,
  });
};

export const deleteMapVersion = (ids) => {
  return request.delete("/mapVersion", {
    data: { ids },
    needToken: true,
  });
};

export const changeMapVersion = (pramas) => {
  return request.put(`/changeMapVersion${pramas}`, {
    needToken: true,
  });
};
