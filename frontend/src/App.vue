<template>
  <v-app id="inspire" dark>
    <v-navigation-drawer v-model="drawer" clipped fixed dark app>
      <v-list dense dark>
        <!-- Components -->
        <template v-for="(component,index) in $options.components">
          <v-list-item v-if="index !== 'App'" :key="index" dark>
            <v-list-item-action @click="changeComponent(index)">
              <v-icon>{{ icons[index] }}</v-icon>
            </v-list-item-action>
            <v-list-item-content @click="changeComponent(index)">
              <v-list-item-title>{{ index | capitalize }}</v-list-item-title>
            </v-list-item-content>
          </v-list-item>
        </template>

        <!-- Connections title -->
        <v-list-item class="Connection__title">
          <v-list-item-content>
            Connections
            <v-progress-circular
              v-show="connectionIsChanging"
              indeterminate
              color="primary"
              class="Connection__loading"
              style="width:22px;height:22px;"
            />
          </v-list-item-content>
        </v-list-item>

        <!-- Connections -->
        <v-list-item v-for="(conn,index) in connections" :key="index">
          <v-list-item-content class="Connection">
            <v-icon v-show="conn.Index !== connectionActive" class="Connection__delete" @click="deleteConnection(conn)">
              mdi-delete
            </v-icon>
            <v-icon v-show="conn.Index === connectionActive" class="Connection__active">
              mdi-check
            </v-icon>
            <span class="Connection__name" @click="switchConnection(conn)">
              {{ conn.Name }}
            </span>
            <v-progress-circular
              v-show="conn.Index === connectionActive && isLoadingTables"
              indeterminate
              color="primary"
              class="Connection__loading"
              style="width:22px;height:22px;"
            />
          </v-list-item-content>
        </v-list-item>
      </v-list>
    </v-navigation-drawer>

    <v-app-bar
      dark
      app fixed clipped-left
    >
      <v-app-bar-nav-icon @click.stop="drawer = !drawer" />
      <v-toolbar-title>{{ component | capitalize }}</v-toolbar-title>
    </v-app-bar>
    <v-content dark>
      <v-container dark fluid class="px-0">
        <v-layout dark ustify-center align-center class="px-0">
          <component :is="component" :key="update" :tables="tables" :is-loading-tables="isLoadingTables" dark @createconnection="createConnection" />
        </v-layout>
      </v-container>
    </v-content>
    <v-footer app dark fixed>
      <span style="margin-left:1em">&copy; Baron Tiberiu</span>
    </v-footer>

    <v-snackbar
      v-model="snack.show"
      :top="true"
      :color="snack.color"
      @click:outside="cancelDelete"
    >
      {{ snack.currentMessage }}
      <v-btn
        text
        @click="snack.reset"
      >
        Close
      </v-btn>
    </v-snackbar>

    <v-dialog
      v-model="isDeletingConnection"
      max-width="350"
    >
      <v-card>
        <v-card-title class="headline">
          Delete
        </v-card-title>

        <v-card-text>
          Are you sure that you want to delete the connection ?
        </v-card-text>

        <v-card-actions>
          <v-spacer />

          <v-btn
            color="green darken-1"
            text
            @click="cancelDelete"
          >
            Cancel
          </v-btn>

          <v-btn
            color="red darken-1"
            text
            @click="continueDelete"
          >
            Continue
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </v-app>
</template>

<script>
import Vue from 'vue'
import CheckIntegriy from '@/views/CheckIntegrity.vue'
import Connection from '@/views/Connection.vue'
import RunSQL from '@/views/Run.vue'
import Tables from '@/views/Tables.vue'

import TableComponent from '@/components/Table'
import { migrate, createFakeData } from '@/utils/migration'
import { getInstanceQueueMessage } from '@/utils/Queue.js'
import { WAILSINIT, CREATECONNECTION, REFRESHTABLES } from '@/store/events'
import { AppConnections } from '@/features/connections'
import Table from '@/utils/Table'

Vue.component('Table', TableComponent)

const icons = {
  check_integrity: 'mdi-database-check',
  connection: 'mdi-settings',
  run_SQL: 'mdi-database',
  tables: 'mdi-database'
}
export default {
  name: 'App',
  components: {
    check_integrity: CheckIntegriy,
    connection: Connection,
    run_SQL: RunSQL,
    tables: Tables
  },
  filters: {
    capitalize: function (value) {
      if (!value) return ''
      value = value.toString().replace('_', ' ')
      return value.charAt(0).toUpperCase() + value.slice(1)
    }
  },
  props: {
    // source: String
  },
  data: () => ({

    ...AppConnections().data(),

    drawer: true,
    component: 'check_integrity',

    icons,

    tables: [],
    isLoadingTables: false,
    update: 0,

    snack: getInstanceQueueMessage()
  }),

  created () {
    this.init()

    var _vm = this
    window.run = async function (obj) {
      var migr = true; var create = true
      if (obj && typeof obj.migrate !== 'undefined') {
        migr = obj.migrate
      }

      if (obj && typeof obj.create !== 'undefined') {
        create = obj.create
      }
      try {
        if (migr) {
          await migrate(async (sql) => {
            await _vm.$backend.Run({ run: sql })
          })
        }
        if (create) {
          createFakeData(async (sql) => {
            await _vm.$backend.Run({ run: sql })
          })
        }
      } catch (e) {
        console.log(e)
      }
    }
  },

  methods: {

    registerEvents () {
      this.$root.$on(WAILSINIT, this.resetLoadings)
      this.$root.$on(CREATECONNECTION, this.createConnection)
      this.$root.$on(REFRESHTABLES, this.getTables)
    },
    async init () {
      this.registerEvents()
      this.takeConnections()
      this.getTables()

      // const test = () => {
      //   const table = {
      //     Name: 'test',
      //     Columns: [
      //       {
      //         Name: 'test1',
      //         Constraints: [
      //           {
      //             Name: 't1',
      //             Type: ''
      //           }
      //         ]
      //       }
      //     ]
      //   }
      //   console.log(module)
      //   var m = new Table(table)

      //   console.log(m)
      //   window.test = function () {
      //     m.Columns[0].Constraints.push({ Name: 's', Type: 'NOT NULL' })
      //     console.log(m)
      //   }
      // }

      // test()
    },

    async getTables () {
      if (this.isLoadingTables) return

      this.isLoadingTables = true

      var start = Date.now()
      try {
        const rez = await this.$backend.GetTables()

        if (rez.data) {
          this.tables = rez.data.map(table => {
            return new Table(table)
          })
        }

        getInstanceQueueMessage().addMessage(rez)
      } catch (e) {
        console.log(e)
      }

      var end = Date.now()
      var elapsed = end - start // time in milliseconds
      console.log(300 - elapsed)
      setTimeout(() => {
        this.isLoadingTables = false
        this.update++

        console.log(this)
      }, 300 - elapsed)
    },

    resetLoadings () {
      this.connectionIsChanging = false
    },

    changeComponent (index) {
      if (this.component !== index) {
        this.component = index
      }
    },

    ...AppConnections().methods
  }

}
</script>

<style lang="scss">
.v-list-item {
  cursor: pointer;
  &:hover {
    opacity: 0.8;
  }
}

.Connection {
  display: block;

  &__name {
    display: inline-block;
    margin: 0 0 0 15px;
    width: 50%;
    &:hover {
      opacity: 0.6;
    }
  }

  &__delete {
    display: inline;
    padding: 2px;
    &:hover {
      opacity: 0.6;
    }
  }

  &__active {
    display: inline;
    padding: 2px;
    &:hover {
      opacity: 0.6;
    }
  }

  &__title {
    margin: 20px 20px 0px 20px;
    padding: 0 10px 0 0;
    .v-list-item__content {
      display: block;
    }
  }

  &__loading {
    margin: 0 0 0 25px;
  }
}
</style>
