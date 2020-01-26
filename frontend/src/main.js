import 'babel-polyfill'

import Vue from 'vue'
import 'vuetify/dist/vuetify.min.css'
import '@mdi/font/css/materialdesignicons.css'

import App from './App.vue'

import Vuetify from 'vuetify'
import '@/plugins'

import * as Wails from '@wailsapp/runtime'

import { WAILSINIT } from '@/store/events'

Vue.use(Vuetify)
const opts = { theme: { dark: false } }

var vuetify = new Vuetify(opts)
Vue.config.productionTip = false
Vue.config.devtools = true

var cApp = null

function parse (key) {
	return async function () {
		const transformedArguments = []
		let i
		for (i in arguments) {
			transformedArguments.push(JSON.stringify(arguments[i]))
		}

		console.log('SEND ==>', transformedArguments)
		var response = await window.backend[key](...transformedArguments)
		try {
			if (typeof response === 'string') {
				response = JSON.parse(response)
			}
			console.log('RECEIVE <==', response)

			return response
		} catch (e) { console.log(e) }

		return false
	}
}
Wails.Init(() => {
	if (cApp == null) {
		var functions = {}
		Object.keys(window.backend).map(key => {
			functions[key] = parse(key)
		})
		Vue.prototype.$backend = functions

		cApp = new Vue({
			vuetify,
			render: h => h(App)
		}).$mount('#app')
	} else {
		cApp.$root.$emit(WAILSINIT)
	}
})
