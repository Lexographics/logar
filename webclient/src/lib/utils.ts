
export function getApiUrl() : string {
  if(process.env.PUBLIC_DEV_API) {
    return process.env.PUBLIC_DEV_API;
  }
  

  const apiUrl = getCookie("api-url");
  if(apiUrl) {
    return apiUrl;
  }

  return "";
}

export function getBasePath() : string {
  if(process.env.PUBLIC_DEV_API) {
    return "";
  }

  const basePath = getCookie("base-path");
  if(basePath) {
    return basePath;
  }
  return "";
}

export function setCookie(name: string, val: string) {
  const date = new Date();
  const value = val;

  date.setTime(date.getTime() + (7 * 86400 * 1000));
  document.cookie = name+"="+value+"; expires="+date.toUTCString()+"; path=/";
}

export function getCookie(name: string) {
  const value = "; " + document.cookie;
  const parts = value.split("; " + name + "=");
  
  if (parts.length == 2) {
      return parts.pop().split(";").shift();
  }
}

export function deleteCookie(name: string) {
  const date = new Date();
  date.setTime(date.getTime() + (-1 * 86400 * 1000));
  document.cookie = name+"=; expires="+date.toUTCString()+"; path=/";
}