import request from "@/util/request";

export const login = (data) => {
  return request.post("/login", data, {
    needToken: false,
  });
};
export const verifyEmail = (id) => {
  return request.get(`/verify/${id}`, {
    needToken: false,
  });
};
export const security = (data) => {
  return request.put("/security", data, {
    needToken: false,
  });
};
export const updatepass = (data) => {
  return request.put("/updatepass", data, {
    needToken: true,
  });
};
