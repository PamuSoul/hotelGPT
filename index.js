function submitForm(action) {
  let form = document.getElementById("loginregister");
  let formData = new FormData(form);

  let url =
    action === "register"
      ? "http://localhost:8080/api/v1/account/register"
      : "http://localhost:8080/api/v1/account/login";

  fetch(url, {
    method: "POST",
    body: formData,
  })
    .then((response) => response.json())
    .then((data) => {
      if (data.message === "登入成功") {  // 確保後端回應成功
        window.location.href = "/chat.html"; // 跳轉到儀表板頁面
      }
      if (data.message === "創建帳號成功"){
        alert("創建成功")
      }
    const username = data.username; 
      localStorage.setItem('username', username);
    }) 
        
    .catch((error) => {
      console.error("錯誤:", error);
    });
}
