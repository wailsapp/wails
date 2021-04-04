// Get input + focus
let nameElement = document.getElementById("name");
nameElement.focus();

// Setup the greet function
window.greet = function () {

  // Get name
  let name = nameElement.value;

  // Call Basic.Greet(name)
  window.backend.main.Basic.Greet(name).then((result) => {
    // Update result with data back from Basic.Greet()
    document.getElementById("result").innerText = result;
  });
}