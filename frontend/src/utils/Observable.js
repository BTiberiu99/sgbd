// define a class
export default function (vm) {
    if (!vm) return

    var observers = []

    const _notify = async function () {
        let i = 0

        for (; i < observers.length; i++) {
            observers[i](...arguments)
        }
    }
    // add the ability to subscribe to a new object / DOM element
    // essentially, add something to the observers array
    vm.subscribe = function (f) {
        observers.push(f)

        return this
    }

    // add the ability to unsubscribe from a particular object
    // essentially, remove something from the observers array
    vm.unsubscribe = function (f) {
        observers = observers.filter(subscriber => subscriber !== f)

        return this
    }

    // update all subscribed objects / DOM elements
    // and pass some data to each of them
    vm.notify = function () {
        _notify(...arguments)

        return this
    }
}
