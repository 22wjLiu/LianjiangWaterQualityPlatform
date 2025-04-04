import request from "@/util/request";
export const submit = (data) => {
  return request.post("/regist", data, { 
    needToken: false 
  });
};

export const getVerify = (id) => {
  return request.get(`/verify/${id}`, { 
    needToken: false 
  });
};
