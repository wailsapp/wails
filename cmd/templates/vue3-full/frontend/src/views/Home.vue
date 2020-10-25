<template>
  <div class="home">
    <img @click="getMessage" alt="Vue logo" src="../assets/appicon.png" :style="{ height: '400px' }"/>
    <HelloWorld :msg="message" />
  </div>
</template>

<script lang="ts">
import { ref, defineComponent } from "vue";
import HelloWorld from "@/components/HelloWorld.vue"; // @ is an alias to /src

interface Backend {
  basic(): Promise<string>;
}

declare global {
  interface Window {
    backend: Backend;
  }
}

export default defineComponent({
  name: "Home",
  components: {
    HelloWorld,
  },
  setup() {
    
    const message = ref("Click the Icon");

    const getMessage = () => {
      window.backend.basic().then(result => {
        message.value = result;
      });
    }

    return { message: message, getMessage: getMessage };
  },
});
</script>
