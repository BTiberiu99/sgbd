const timeoutInterval = 2 * 1000 // 2 seconds
const queue = function () {
    // Private props
    var messages = []

    // Self instance
    var _self

    // Interval
    var interval = null

    // Stop queue and hide
    function stopQueue () {
        clearInterval(interval)
        _self.show = false
        _self.currentMessage = ''
        _self.color = ''
        interval = null
    }

    function next () {
        var msg = nextMessage()
        if (!msg) {
            return stopQueue()
        }

        _self.currentMessage = msg.text
        _self.color = msg.color

        if (!_self.show) {
            _self.show = true
        }
    }
    // start interval
    function startInterval () {
        next()

        interval = setInterval(function () {
            next()
        }, timeoutInterval)
    }

    // startQueue
    function startQueue () {
        if (interval !== null) return

        startInterval()
    }

    // get next message
    function nextMessage () {
        if (messages.length === 0) {
            return false
        }
        return messages.pop()
    }

    // create object
    _self = {
        show: false,
        currentMessage: '',
        color: '',
        addMessage: function (response) {
            if (response.message && response.type) {
                messages.push({
                    text: response.message,
                    color: response.type
                })

                startQueue()
            }
        },
        reset: function () {
            clearInterval(interval)
            startInterval()
        }
    }

    return _self
}

var instance = null

// Singleton
export const getInstanceQueueMessage = () => {
    if (instance === null) {
        instance = queue()
    }
    return instance
}
