import request from "@/util/request";

export const getDataTableInfos = () => {
  return request.get("/dataTableInfos", {
    needToken: true,
  });
}

export const getDataBackStage = (start, end, pramas) => {
  let path = "/dataBackStage";

  path += start ? `/${start}` : "/null";
  path += end ? `/${end}` : "/null";
  path += pramas;

  return request.get(path, {
    needToken: true,
  });
};

export const deleteDataBackStage = (body) => {
  return request.delete("/dataBackStage", {
    data: body,
    needToken: true,
  });
};

export const updateDataBackstage = (pramas, body) => {
  return request.put(`/dataBackStage${pramas}`, body, {
    needToken: true,
  });
};