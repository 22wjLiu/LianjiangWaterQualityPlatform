import request from "@/util/request";

export const download = (path, file) => {
  return request.get(`/download?${path}${file}`, {
    needToken: true,
    responseType: "arraybuffer",
  });
};

export const upload = (system, data) => {
  return request.post(`/upload/${system}`, data, {
    needToken: true,
  });
};

export const getFileList = (path) => {
  return request.get(`/files?${path}`, {
    needToken: true,
  });
};


export const getFileInfos = (start, end, pramas) => {
  let path = "/fileInfos";

  path += start ? `/${start}` : "/null";
  path += end ? `/${end}` : "/null";
  path += pramas;

  return request.get(path, {
    needToken: true,
  });
};
