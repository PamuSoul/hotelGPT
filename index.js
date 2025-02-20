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
      document.getElementById("responseMessage").innerText = data.message;
    })
    .catch((error) => {
      console.error("錯誤:", error);
    });
}
