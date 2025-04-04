export const timeToISO = (time) => {
  return time.toISOString().split(".")[0];
}

export const strToISO = (str) => {
  const time = new Date(str)
  return time.toISOString().split(".")[0];
}

export const formatTime = (rawTime) => {
  const date = new Date(rawTime);
  const yyyy = date.getFullYear();
  const MM = String(date.getMonth() + 1).padStart(2, "0");
  const dd = String(date.getDate()).padStart(2, "0");
  const hh = String(date.getHours()).padStart(2, "0");
  const mm = String(date.getMinutes()).padStart(2, "0");
  return `${yyyy}-${MM}-${dd} ${hh}:${mm}`;
};
