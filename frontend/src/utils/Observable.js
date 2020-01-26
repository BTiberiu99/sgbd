// define a class
export default function () {
    var observers = []

    // add the ability to subscribe to a new object / DOM element
    // essentially, add something to the observers array
    this.subscribe = function (f) {
        observers.push(f)
    }

    // add the ability to unsubscribe from a particular object
    // essentially, remove something from the observers array
    this.unsubscribe = function (f) {
        observers = observers.filter(subscriber => subscriber !== f)
    }

    // update all subscribed objects / DOM elements
    // and pass some data to each of them
    this.notify = function () {
        observers.forEach(observer => observer(...arguments))
    }
}
