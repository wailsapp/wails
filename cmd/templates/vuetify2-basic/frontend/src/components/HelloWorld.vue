<template>
  <v-container fluid class="px-0">
    <v-layout>
      <v-flex xs12 sm6 offset-sm3>
        <v-card raised="raised" class="pa-4 ma-4">
          <v-layout justify-center align-center class="pa-4 ma-4">
            <v-img :src="require('../assets/images/logo.png')"></v-img>
          </v-layout>
          <v-card-actions>
            <v-layout justify-center align-center class="px-0">
              <v-btn color="blue" @click="getMessage">Press Me</v-btn>
            </v-layout>
          </v-card-actions>
        </v-card>
      </v-flex>
    </v-layout>
    <div class="text-xs-center">
      <v-dialog v-model="dialog" width="500">
        <v-card>
          <v-card-title class="headline" primary-title>Message from Go</v-card-title>
          <v-card-text>{{message}}</v-card-text>
          <v-divider></v-divider>
          <v-card-actions>
            <v-spacer></v-spacer>
            <v-btn color="primary" text @click="dialog = false">Awesome</v-btn>
          </v-card-actions>
        </v-card>
      </v-dialog>
    </div>
  </v-container>
</template>

<script>
  export default {
    data () {
      return {
        message: " ",
        raised: true,
        dialog: false
      }
    },
    methods: {
      getMessage: function () {
        var self = this
        window.backend.basic().then(result => {
          self.message = result
          self.dialog = true
        })
      }
    }
  }
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
  h1 {
    margin-top: 2em;
    position: relative;
    min-height: 5rem;
    width: 100%;
  }

  a:hover {
    font-size: 1.7em;
    border-color: blue;
    background-color: blue;
    color: white;
    border: 3px solid white;
    border-radius: 10px;
    padding: 9px;
    cursor: pointer;
    transition: 500ms;
  }

  a {
    font-size: 1.7em;
    border-color: white;
    background-color: #121212;
    color: white;
    border: 3px solid white;
    border-radius: 10px;
    padding: 9px;
    cursor: pointer;
  }
</style>
