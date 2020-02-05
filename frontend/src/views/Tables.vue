<template>
  <v-container fluid class="Tables">
    <v-row style="padding:20px;">
      <template v-if="!isLoadingTables">
        <v-col v-for="(table,index) in tables" :key="index" class="Table" cols="12" lg="6">
          <Table :key="table.keyVue" :table="table" />
        </v-col>
      </template>
      <v-progress-circular
        v-show="isLoadingTables"
        indeterminate
        color="primary"
        style="margin:auto;"
      />
    </v-row>
  </v-container>
</template>

<script>

export default {
  name: 'Tables',

  props: {
    tables: {
      type: Array,
      required: true
    },
    isLoadingTables: {
      type: Boolean,
      required: true
    }
  },
  data () {
    return {
      // tables: [],
      // isLoading: false
    }
  },

  created () {
    this.init()
  },
  methods: {

    init () {
      // this.getTables()
    }

  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style lang="scss" >
.v-icon.danger {
  color: #ff3c00ce;
}
</style>
