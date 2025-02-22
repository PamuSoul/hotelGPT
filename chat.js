document.addEventListener("DOMContentLoaded", function () {
  const username = localStorage.getItem("username");

  if (username) {
    fetch(`http://localhost:8080/api/v1/account/history?username=${username}`, {
      method: "GET",
      headers: { "Content-Type": "application/json" },
    })
      .then((response) => response.json())
      .then((data) => {
        const chatContainer = document.querySelector(".chat-container");

        // 遍歷伺服器回傳的歷史紀錄，並將每條訊息加到 .chat-container 中
        data.history.forEach((message) => {
        
          const userMessageDiv = document.createElement("div");
          userMessageDiv.classList.add("usermessage");
          userMessageDiv.innerHTML = `<span class="sender">${message.username}:</span><span>${message.message}</span>`;
          chatContainer.appendChild(userMessageDiv);

          const gptMessageDiv = document.createElement("div");
          gptMessageDiv.classList.add("gptmessage");
          gptMessageDiv.innerHTML = `<span class="sender">GPT:</span><span>${message.gptmessage}</span>`;
          chatContainer.appendChild(gptMessageDiv);
        });

        chatContainer.scrollTop = chatContainer.scrollHeight;
      })
      .catch((error) => console.error("Error fetching history:", error));
  }
});

document
  .getElementById("chatFrom")
  .addEventListener("submit", function (event) {
    event.preventDefault();

    const message = document.getElementById("messageinput").value;
    const username = localStorage.getItem("username");
    fetch("http://localhost:8080/api/v1/account/chat", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ username, message }),
    })
      .then((response) => response.json())
      .then((data) => {
        const chatContainer = document.querySelector(".chat-container");

        const userMessageDiv = document.createElement("div");
        userMessageDiv.classList.add("usermessage");
        userMessageDiv.innerHTML = `<span class="sender">${data.username}:</span><span>${data.message}</span>`;
        chatContainer.appendChild(userMessageDiv);

        const gptMessageDiv = document.createElement("div");
        gptMessageDiv.classList.add("gptmessage");
        gptMessageDiv.innerHTML = `<span class="sender">GPT:</span><span>${data.gptmessage}</span>`;
        chatContainer.appendChild(gptMessageDiv);

        chatContainer.scrollTop = chatContainer.scrollHeight;
        // 清空輸入框
        document.getElementById("messageinput").value = "";
      })
      .catch((error) => console.error("Error:", error));
  });
