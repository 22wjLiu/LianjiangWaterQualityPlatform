import request from "@/util/request";

export const download = (pramas) => {
  return request.get("/download" + pramas, {
    needToken: true,
    responseType: "arraybuffer",
  });
};

export const upload = (system, data) => {
  return request.post(`/upload/${system}`, data, {
    needToken: true,
  });
};

export const getFileList = (pramas) => {
  return request.get("/files" + pramas, {
    needToken: true,
  });
};

export const deleteFiles = (ids) => {
  return request.delete("/files", {
    data: { ids },
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

export const updateFileName = (id, body) => {
  return request.put(`/fileName/${id}`, body, {
    needToken: true,
  });
};