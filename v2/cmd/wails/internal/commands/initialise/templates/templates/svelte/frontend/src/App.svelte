<script>
  let name = "";

  let greeting = "Please enter your name below 👇";

  function greet() {
    // Check if the input is empty
    if (name === "") return;

    // Call App.Greet(name)
    try {
      window.go.main.App.Greet(name)
        .then((result) => {
          // Update result with data back from App.Greet()
          greeting = result;
        })
        .catch((err) => {
          console.error(err);
        });
    } catch (err) {
      console.error(err);
    }
  }
</script>

<div class="logo" />
<div class="result" id="result">{greeting}</div>
<div class="input-box" id="input" data-wails-no-drag>
  <input
    class="input"
    id="name"
    type="text"
    autocomplete="off"
    bind:value={name}
    on:keydown={(e) => e.key === "Enter" && greet()}
  />
  <button class="btn" on:click={greet}>Greet</button>
</div>

<style>
  .logo {
    display: block;
    width: 50%;
    height: 50%;
    margin: auto;
    padding: 10% 0 0;
    background-position: center;
    background-repeat: no-repeat;
    background-image: url("./assets/images/logo-universal.png");
    background-size: 100% 100%;
    background-origin: content-box;
  }
  .result {
    height: 20px;
    line-height: 20px;
    margin: 1.5rem auto;
  }
  .input-box .btn {
    width: 60px;
    height: 30px;
    line-height: 30px;
    border-radius: 3px;
    border: none;
    margin: 0 0 0 20px;
    padding: 0 8px;
    cursor: pointer;
  }
  .input-box .btn:hover {
    background-image: linear-gradient(to top, #cfd9df 0%, #e2ebf0 100%);
    color: #333333;
  }

  .input-box .input {
    border: none;
    border-radius: 3px;
    outline: none;
    height: 30px;
    line-height: 30px;
    padding: 0 10px;
    background-color: rgba(240, 240, 240, 1);
    -webkit-font-smoothing: antialiased;
  }

  .input-box .input:hover {
    border: none;
    background-color: rgba(255, 255, 255, 1);
  }

  .input-box .input:focus {
    border: none;
    background-color: rgba(255, 255, 255, 1);
  }
</style>
