import request from "@/util/request";
import { timeToISO } from "@/util/timeFormater";

export const getDataLog = (start, end, pramas) => {
  let path = "/history/data";

  path += start ? `/${timeToISO(start)}` : "/null";
  path += end ? `/${timeToISO(end)}` : "/null";
  path += pramas;

  return request.get(path, {
      needToken: true
    }
  );
};

export const deleteDataLog = (pramas) => {
  return request.delete(
    `/history/data${pramas}`, {
      needToken: true
    }
  );
};

export const getFileLog = (start, end, pramas) => {
  let path = "/history/file";

  path += start ? `/${timeToISO(start)}` : "/null";
  path += end ? `/${timeToISO(end)}` : "/null";
  path += pramas;
  
  return request.get(path, {
      needToken: true
    }
  );
};

export const deleteFileLog = (pramas) => {
  return request.delete(
    `/history/file${pramas}`, {
      needToken: true
    }
  );
};

export const getMapLog = (start, end, pramas) => {
  let path = "/history/map";

  path += start ? `/${timeToISO(start)}` : "/null";
  path += end ? `/${timeToISO(end)}` : "/null";
  path += pramas;
  
  return request.get(path, {
      needToken: true
    }
  );
};

export const deleteMapLog = (pramas) => {
  return request.delete(
    `/history/map${pramas}`, {
      needToken: true
    }
  );
};