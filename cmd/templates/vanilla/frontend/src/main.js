
const runtime = require('@wailsapp/runtime');

// We need to wait for runtime.Init to complete before 
// running our JS
runtime.Init(() => {

    // Ensure the default app div is 100% wide/high
    var app = document.getElementById("app");
    app.style.width = "100%";
    app.style.height = "100%";

    // Inject html
    app.innerHTML = `
    <div class="logo"></div>
    <div class="container">
        <button id="button">Click Me!</button>
        <div id="result"/>
    </div>
    `;

    // Connect button to Go method
    document.getElementById("button").onclick = () => {
       window.backend.basic().then((result) => {
           document.getElementById("result").innerText = result;
       })
    } 

});