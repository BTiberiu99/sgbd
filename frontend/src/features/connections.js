import { getInstanceQueueMessage } from '@/utils/Queue.js'
import { ENDCREATECONNECTION } from '@/store/events'
export const AppConnections = function () {
  return {
    data () {
      return {
        connections: [],
        connectionActive: '',
        connectionIsChanging: false,
        isDeletingConnection: false,
        delete: false
      }
    },
    methods: {

      // take all connections from the store
      takeConnections () {
        this.connectionIsChangingRun(async () => {
          const rez = await this.$backend.GetConnections()
          if (rez.data) {
            this.connections = rez.data.Connections
            this.connectionActive = rez.data.Index
          }
        })
      },

      // switch between connections
      switchConnection (conn) {
        this.connectionIsChangingRun(async () => {
          const rez = await this.$backend.SwitchConnection(conn)

          this.connectionActive = rez.data.Index
        })
      },

      // continue the delete of a connection
      continueDelete () {
        this.delete = true
        this.isDeletingConnection = false
      },
      // cancel the delete of a connection
      cancelDelete () {
        this.delete = false
        this.isDeletingConnection = false
      },
      // delete connection from the store
      async deleteConnection (conn) {
        this.delete = false
        this.isDeletingConnection = true

        await this.$sync(() => { return !this.isDeletingConnection })

        if (!this.delete) return

        this.connectionIsChangingRun(async () => {
          const rez = await this.$backend.RemoveConnection(conn)

          this.connections = rez.Data

          getInstanceQueueMessage().addMessage(rez)
        })
      },
      // create new connection
      createConnection (conn) {
        var rez
        this.connectionIsChangingRun(async () => {
          rez = await this.$backend.CreateConnection(conn)

          if (rez.data) {
            var add = true
            this.connections.map(conn => {
              if (conn.Index === rez.data.Index) {
                add = false
              }
            })
            if (add) {
              this.connections.push(rez.data)
            }

            this.connectionActive = rez.data.Index
          }

          getInstanceQueueMessage().addMessage(rez)
        }, () => {
          this.$root.$emit(ENDCREATECONNECTION, rez)
        })
      },

      // run changes to the connections
      async connectionIsChangingRun (call, final) {
        if (typeof call !== 'function') {
          throw new Error('First parameter must be an callback')
        }
        console.log(this.connectionIsChanging)
        if (this.connectionIsChanging) return

        this.connectionIsChanging = true
        var start = Date.now()
        try {
          await call()
        } catch (e) {
          // console.log(e)
        }
        // the event you'd like to time goes here:
        var end = Date.now()
        var elapsed = end - start // time in milliseconds
        setTimeout(() => {
          this.connectionIsChanging = false
          if (typeof final === 'function') {
            final()
          }
          this.getTables()
        }, 300 - elapsed)
      }
    }
  }
}
