<template>
  <div class="home">
    <img
      @click="getMessage"
      alt="Vue logo"
      src="../assets/appicon.png"
      :style="{ height: '400px' }"
    />
    <HelloWorld :msg="message" />
  </div>
</template>

<script>
import { ref } from "vue";
import HelloWorld from "@/components/HelloWorld.vue";

export default {
  name: "Home",
  components: {
    HelloWorld,
  },
  setup() {
    const message = ref("Click the Icon");

    const getMessage = () => {
      window.backend
        .basic()
        .then((res) => {
          message.value = res;
        })
        .catch((err) => {
          console.log(err);
        });
    };

    return {
      message: message,
      getMessage: getMessage,
    };
  },
};
</script>
