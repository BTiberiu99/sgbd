<template>
  <v-container fluid class="px-0">
    <v-row style="padding:25px;">
      <v-col cols="4" style="text-align:center;">
        <div style="width:200px;">
          Tables
          <v-select v-model="table"
                    :items="tables"
                    label="Tables"
                    dense
                    return-object
                    persistent-hint
                    single-line
                    item-text="Name"
                    :hint="hint(table)"
          />
        </div>
      </v-col>
      <v-col cols="8" style="text-align:center;" />
      <v-col cols="8" offset="2" style="text-align:center;">
        <Table v-if="table" :table="table" style="margin:0 0 20px 0 ;" />
        <v-btn v-if="table && !table.HasOneNotNull()" style="margin:0 40px 0 0;" @click="addNotNull(table)">
          Add Not Null
        </v-btn>

        <v-btn v-if="table && !table.HasPrimaryKey()" @click="addPrimaryKey(table)">
          Add Primary Key
        </v-btn>
      </v-col>
    </v-row>
  </v-container>
</template>

<script>
import { WAILSINIT } from '@/store/events'
export default {
  name: 'CheckIntegriy',
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
      table: null
    }
  },

  computed: {
    // tablesName () {

    // }
  },

  created () {
    this.$root.$on(WAILSINIT, this.init)
  },
  methods: {
    init () {

    },
    getMessage () {
      // var self = this;

    },
    hint (table) {
      var hint = ''
      if (table && !table.IsSafe()) {
        var hasPrimaryKey = table.HasPrimaryKey()
        var hasOneNotNull = table.HasOneNotNull()

        hint = `Tabelul ${table.Name} nu are ${!hasPrimaryKey && !hasOneNotNull ? ' cheie primara si cel putin o coloana not null' : (hasPrimaryKey ? ' cel putin o coloana not null' : ' cheie primara')}`
      }
      return hint
    },

    addPrimaryKey (table) {

    },
    addNotNull (table) {

    }
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
</style>
