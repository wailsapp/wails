<script setup>
import { ref, onMounted } from 'vue'
import {GreetService} from "../../bindings/changeme";
import {Events} from "@wailsio/runtime";

const name = ref('')
const result = ref('Please enter your name below 👇')
const time = ref('Listening for Time event...')

const doGreet = () => {
  let localName = name.value;
  if (!localName) {
    localName = 'anonymous';
  }
  GreetService.Greet(localName).then((resultValue) => {
    result.value = resultValue;
  }).catch((err) => {
    console.log(err);
  });
}

onMounted(() => {
  Events.On('time', (timeValue) => {
    time.value = timeValue.data;
  });
})

defineProps({
  msg: String,
})

</script>

<template>
  <h1>{{ msg }}</h1>

  <div class="result">{{ result }}</div>
  <div class="card">
    <div class="input-box">
      <input class="input" v-model="name" type="text" autocomplete="off"/>
      <button class="btn" @click="doGreet">Greet</button>
    </div>
  </div>

  <div class="footer">
    <div><p>Click on the Wails logo to learn more</p></div>
    <div><p>{{ time }}</p></div>
  </div>
</template>
