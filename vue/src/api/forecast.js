import request from "@/util/request";

export const forecast = (start, end, pramas) => {
  let path = "/forecast";

  path += start ? `/${start}` : "/null";
  path += end ? `/${end}` : "/null";
  path += pramas;

  return request.get(path, {
    needToken: true,
  });
};
