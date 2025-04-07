import request from "@/util/request";

export const getDataLog = (start, end, pramas) => {
  let path = "/history/data";

  path += start ? `/${start}` : "/null";
  path += end ? `/${end}` : "/null";
  path += pramas;

  return request.get(path, {
    needToken: true,
  });
};

export const deleteDataLog = (ids) => {
  return request.delete("/history/data", {
    data: { ids },
    needToken: true,
  });
};

export const getFileLog = (start, end, pramas) => {
  let path = "/history/file";

  path += start ? `/${start}` : "/null";
  path += end ? `/${end}` : "/null";
  path += pramas;

  return request.get(path, {
    needToken: true,
  });
};

export const deleteFileLog = (ids) => {
  return request.delete("/history/file", {
    data: { ids },
    needToken: true,
  });
};

export const getMapLog = (start, end, pramas) => {
  let path = "/history/map";

  path += start ? `/${start}` : "/null";
  path += end ? `/${end}` : "/null";
  path += pramas;

  return request.get(path, {
    needToken: true,
  });
};

export const deleteMapLog = (ids) => {
  return request.delete("/history/map", {
    data: { ids },
    needToken: true,
  });
};
