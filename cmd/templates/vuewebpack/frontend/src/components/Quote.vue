<template>
  <div class="container">
    <blockquote v-if="quote != null" :cite="quote.person">{{ quote.text }}</blockquote>
    <p></p>
    <button @click="getNewQuote()">Get new Quote</button>
  </div>
</template>

<script>
import "../assets/css/quote.css";
import { eventBus } from "../main";

export default {
  data() {
    return {
      quote: null
    };
  },
  methods: {
    getNewQuote: function() {
      var self = this;
      backend.QuotesCollection.GetQuote().then(result => {
        self.quote = result;
      });
    }
  },
  created() {
    eventBus.$on("ready", this.getNewQuote);
  }
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
.hello {
  margin-top: 2em;
  position: relative;
  width: 100%;
}
h3 {
  margin: 40px 0 0;
}
ul {
  list-style-type: none;
  padding: 0;
}
li {
  display: inline-block;
  margin: 0 10px;
}
a {
  color: gold;
}
</style>
