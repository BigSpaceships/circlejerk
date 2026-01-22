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

async function main() {
  console.log(await getQueue())
}

main()
