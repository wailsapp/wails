import {ready} from '@wails/runtime';

ready(() => {
    // Get input + focus
    let nameElement = document.getElementById("name");
    nameElement.focus();

    // Setup the greet function
    window.greet = function () {

        // Get name
        let name = nameElement.value;

        // Call App.Greet(name)
        window.backend.main.App.Greet(name).then((result) => {
            // Update result with data back from App.Greet()
            document.getElementById("result").innerText = result;
        });
    };
});