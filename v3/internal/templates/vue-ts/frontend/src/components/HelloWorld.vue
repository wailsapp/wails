<script setup lang="ts">
import { ref, onMounted } from 'vue'
import {GreetService} from "../../bindings/changeme";
import {Events} from "@wailsio/runtime";

defineProps<{ msg: string }>()

const name = ref('')
const result = ref('Please enter your name below 👇')
const time = ref('Listening for Time event...')

const doGreet = () => {
  let localName = name.value;
  if (!localName) {
    localName = 'anonymous';
  }
  GreetService.Greet(localName).then((resultValue: string) => {
    result.value = resultValue;
  }).catch((err: Error) => {
    console.log(err);
  });
}

onMounted(() => {
  Events.On('time', (timeValue: { data: string }) => {
    time.value = timeValue.data;
  });
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
