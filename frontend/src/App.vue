<template>
  <v-app id="inspire" dark>
    <v-navigation-drawer v-model="drawer" clipped fixed dark app>
      <v-list dense dark>
        <v-list-item dark>
          <v-list-item-action @click="component = 'check_integrity'">
            <v-icon>mdi-database-check</v-icon>
          </v-list-item-action>
          <v-list-item-content @click="component = 'check_integrity'">
            <v-list-item-title>Check Integrity</v-list-item-title>
          </v-list-item-content>
        </v-list-item>
        <v-list-item>
          <v-list-item-action @click="component = 'connection'" >
            <v-icon>mdi-settings</v-icon>
          </v-list-item-action>
          <v-list-item-content @click="component = 'connection'" >
            <v-list-item-title>Connection</v-list-item-title>
          </v-list-item-content>
        </v-list-item>
        <v-list-item>
          <v-list-item-action @click="component = 'run_SQL'" >
            <v-icon>mdi-database</v-icon>
          </v-list-item-action>
          <v-list-item-content @click="component = 'run_SQL'" >
            <v-list-item-title>Run SQL</v-list-item-title>
          </v-list-item-content>
        </v-list-item>

        <v-list-item style="margin-top:20px;">
          <v-list-item-content style="display:block">
              Connections
              <v-progress-circular
               v-if="connectionIsChanging"
              indeterminate
              color="primary"
              style="margin-left:25px;"
            ></v-progress-circular>
          </v-list-item-content>
          
        </v-list-item>
        <v-list-item v-for="(conn,index) in  connections" :key="index" @click="switchConnection(conn)">
          <v-list-item-content class="Connection" >
            <v-icon  class="Connection__delete" v-show="conn.Index !== connectionActive" @click="deleteConnection(conn)">mdi-delete</v-icon>
            <span class="Connection__name" >
              {{conn.Name}}
            </span>
            <v-icon class="Connection__active" v-show="true">mdi-check</v-icon>
         </v-list-item-content>
        </v-list-item>
      </v-list>
    </v-navigation-drawer>

     <v-app-bar
     dark
      app fixed clipped-left
    >
      <v-app-bar-nav-icon @click.stop="drawer = !drawer"></v-app-bar-nav-icon>
      <v-toolbar-title>{{component | capitalize }}</v-toolbar-title>
    </v-app-bar>
    <v-content dark>
      <v-container dark fluid class="px-0">
        <v-layout dark ustify-center align-center class="px-0">
          <component dark  :is="component" @createconnection="createConnection"/>
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
          <v-card-title class="headline">Delete</v-card-title>
  
          <v-card-text>
            Are you sure that you want to delete the connection ?
          </v-card-text>
  
          <v-card-actions>
            <v-spacer></v-spacer>
  
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
import CheckIntegriy from "./components/CheckIntegrity.vue";
import Connection from './components/Connection.vue';
import RunSQL from './components/Run.vue';
import { getInstanceQueueMessage } from './utils/Queue.js'
import { WAILSINIT, CREATECONNECTION } from '@/store/events'
import { AppConnections } from '@/features/connections'
export default {
  components: {
    check_integrity: CheckIntegriy,
    connection: Connection,
    run_SQL: RunSQL,
  },
  props: {
    source: String
  },
  filters: {
    capitalize: function (value) {
      if (!value) return ''
      value = value.toString().replace('_', ' ')
      return value.charAt(0).toUpperCase() + value.slice(1)
    }
  },
  data: () => ({

    ...AppConnections().data(),
    drawer: true,
    component: 'check_integrity',



    snack: getInstanceQueueMessage(),
  }),

  created () {
    this.init()
  },



  methods: {

    registerEvents () {
      this.$root.$on(WAILSINIT, this.resetLoadings)
      this.$root.$on(CREATECONNECTION, this.createConnection)
    },
    init () {
      this.registerEvents()
      this.takeConnections()

    },

    resetLoadings () {
      this.connectionIsChanging = false
    },

    ...AppConnections().methods
  }

};
</script>

<style lang="scss">
.logo {
  width: 16em;
}

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
    margin: 0 0 0 15px;
    &:hover {
      opacity: 0.6;
    }
  }
}
</style>