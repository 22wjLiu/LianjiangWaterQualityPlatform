import request from "@/util/request";

export const getMapTables = () => {
  return request.get("/mapTables", {
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

export const createMapVersion = (pramas, data) => {
  return request.post(`/mapVersion${pramas}`, data, {
    needToken: true,
  });
};

export const deleteMapVersion = (ids) => {
  return request.delete("/mapVersion", {
    data: { ids },
    needToken: true,
  });
};
