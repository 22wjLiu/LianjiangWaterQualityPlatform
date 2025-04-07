import request from "@/util/request";
import { timeToISO } from "@/util/timeFormater";

export const getUsers = (start, end, pramas) => {
  let path = "/users";

  path += start ? `/${start}` : "/null";
  path += end ? `/${end}` : "/null";
  path += pramas;

  return request.get(path, {
    needToken: true,
  });
};

export const updateUser = (id, body) => {
  return request.put(`/user/${id}`, body, {
    needToken: true,
  });
};

export const deleteUsers = (ids) => {
  return request.delete("/users", {
    data: { ids },
    needToken: true,
  });
};
