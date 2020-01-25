<template>
  <v-container fluid  class="px-0">
    <v-row>
       <v-col cols="12" style="text-align:center;">

       </v-col>
    </v-row>
  </v-container>
</template>

<script>
export default {
  name:'CheckIntegriy',
  data() {
    return {
      message: " ",
      raised: true,
      dialog: false
    };
  },

  created(){
    this.$root.$on('wailsinit', this.init)
  },
  methods: {
    init(){

    },
    getMessage() {
      // var self = this;

    }
  }
};
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
