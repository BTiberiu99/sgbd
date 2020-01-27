<template>
  <v-container fluid class="px-0">
    <v-row style="padding:25px;">
      <v-col cols="4" class="TablesList">
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

        <!-- <v-btn v-if="table && !table.HasPrimaryKey()" @click="addPrimaryKey(table)">
          Add Primary Key
        </v-btn> -->
      </v-col>
    </v-row>

    <v-dialog
      v-model="dialog"
      max-width="350"
    >
      <v-card v-if="table">
        <v-card-title class="headline">
          Add Not Null
        </v-card-title>

        <v-card-text :key="table.vueKey">
          <h3>Columns</h3>
          <template v-for="(column,index) in table.Columns">
            <v-checkbox v-if="!column.HasNotNull() && !column.HasPrimaryKey()"
                        :key="index"
                        v-model="columnsNotNull"
                        :label="`${column.Name} - ${column.Type}`"
                        :value="{index:index, column:column}"
            />
          </template>
        </v-card-text>

        <v-card-actions>
          <v-spacer />

          <v-btn
            color="green darken-1"
            text
            @click="cancelAction"
          >
            Cancel
          </v-btn>

          <v-btn
            color="red darken-1"
            text
            @click="continueAction"
          >
            Continue
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </v-container>
</template>

<script>
import { WAILSINIT } from '@/store/events'
import Column from '@/utils/Column'
import { getInstanceQueueMessage } from '@/utils/Queue.js'
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
      table: null,
      dialog: false,
      columnsNotNull: [],
      continue: false,
      wait: false
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
    cancelAction () {
      this.dialog = false
      this.continue = false
      this.wait = false
    },
    continueAction () {
      this.dialog = false
      this.continue = true
      this.wait = false
    },
    async addPrimaryKey (table) {
      this.wait = true
      this.dialog = true
    },
    async addNotNull (table) {
      this.wait = true
      this.dialog = true
      this.columnsNotNull = []
      await this.$sync(() => !this.wait)

      if (!this.continue) return

      var items = await Promise.all(this.columnsNotNull.map(async (item) => {
        var rez = item
        if (!item.column.HasNotNull()) {
          const response = await this.$backend.AddNotNull(table.Name, item.column)
          if (response.type === 'success') {
            rez.modified = true
            rez.column = response.data
          }

          getInstanceQueueMessage().addMessage(response)
        }

        return rez
      }))

      items.map(item => {
        if (item.modified) {
          table.Columns[item.index] = new Column(item.column)
        }
      })
    }
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style lang="scss">
.TablesList {
  text-align: center;
  .v-text-field__details {
    .v-messages__message.message-transition-enter-to {
      color: #ff3c00ce;
    }
  }
}
</style>
