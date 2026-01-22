async function getQueue() {
  try {
    const response = await fetch(window.location.origin + "/api/queue")

    if (!response.ok) {
      throw new Error(response.status);
    }

    const json = await response.json()

    return json
  } catch (e) {
    throw new Error(e)
  }
}

function enterQueue(type) {
  fetch(window.location.origin + "/api/enter", {
    method: "POST",
    body: JSON.stringify({
      type: type
    })
  })
}

function parseJwt(token) {
  if (!token) {
    return;
  }
  const base64Url = token.split('.')[1];
  const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/');

  const jsonPayload = decodeURIComponent(window.atob(base64).split('').map(function(c) {
    return '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2);
  }).join(''));

  return JSON.parse(jsonPayload);
}

async function getUserInfo() {
  const cookie = await cookieStore.get("Auth");

  return parseJwt(cookie.value).UserInfo;
}

function setUserInfo(userInfo) {
  const username = userInfo.preferred_username;
  const name = userInfo.name;
  document.getElementById("profile-pic").src = `https://profiles.csh.rit.edu/image/${username}`;
  document.getElementById("profile-name").innerText = name;
}

async function main() {
  console.log(await getQueue())
  userInfo = await getUserInfo()
  setUserInfo(userInfo)
}

main()
