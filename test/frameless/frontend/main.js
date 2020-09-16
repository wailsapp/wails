// Get input + focus
var nameElement = document.getElementById("name");
nameElement.focus();

// Stup the greet function
window.greet = function () {

  // Get name
  var name = nameElement.value;

  // Call Basic.Greet(name)
  window.backend.main.Basic.Greet(name).then((result) => {
    // Update result with data back from Basic.Greet()
    document.getElementById("result").innerText = result;
  });
}

window.closeme = function () {
  console.log('here');
  window.backend.main.Basic.Close();
}