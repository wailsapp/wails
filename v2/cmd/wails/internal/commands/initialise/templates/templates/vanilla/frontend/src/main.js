// Get input + focus
let nameElement = document.getElementById("name");
nameElement.focus();

// Setup the greet function
window.greet = function () {

  // Get name
  let name = nameElement.value;

  // Call App.Greet(name)
  window.go.main.App.Greet(name).then((result) => {
    // Update result with data back from App.Greet()
    document.getElementById("result").innerText = result;
  });
};

nameElement.onkeydown = function (e) {
  console.log(e)
  if (e.keyCode == 13) {
    window.greet()
  }
}