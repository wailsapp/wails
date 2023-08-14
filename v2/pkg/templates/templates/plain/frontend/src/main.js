// Get input + focus
let nameElement = document.getElementById("name");
nameElement.focus();
import './main.css';

// Setup the greet function
window.greet = function () {
    // Get name
    let name = nameElement.value;

    // Check if the input is empty
    if (name === "") return;

    // Call App.Greet(name)
    try {
        window.go.main.App.Greet(name)
            .then((result) => {
                // Update result with data back from App.Greet()
                document.getElementById("result").innerText = result;
            })
            .catch((err) => {
                console.error(err);
            });
    } catch (err) {
        console.error(err);
    }
};

nameElement.onkeydown = function (e) {
    if (e.keyCode == 13) {
        window.greet();
    }
};
