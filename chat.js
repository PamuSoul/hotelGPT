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
