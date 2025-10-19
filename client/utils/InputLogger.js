class InputLogger{

    // lists of inputs
    constructor(video){
        this.logs = []
        this.maxLogs = 1000
        this.video = video;
    }

    log(event){
        // log information
        const logEntry = {

            timestamp: Date.now(),
            type: event.type,
            x: event.x,
            y: event.y,
            key: event.key,
            shiftKey :event.shiftKey,
            altKey : event.altKey,
            ctrlKey : event.ctrlKey,
            //metaKey: event.metaKey
            // button: event.button
            //code:event.code
        }
        this.logs.push(logEntry)

        if (this.logs.length>this.maxLogs){
            this.logs.shift()
        }

        console.log('input:' ,logEntry)
        return logEntry

    }
    // method to get complete logs
    getlogs(){
        return this.logs
    }
    exportlogs(){
        return JSON.stringify(this.logs,null,2)
        
    }
    mouseListener(e){
        const rec = this.video.getBoundingClientRect()
        const event = {
            type:e.type,
            x: e.clientX - rec.left,
            y: e.clientY - rec.top
        }
        this.log(event)
    
    }
    keyListener(e){
        const event = {
            type: e.type,
            key: e.key,
            ctrlKey: e.ctrlKey,
            altKey: e.altKey,
            shiftKey: e.shiftKey
        }
        this.log(event)
    }
}

export default InputLogger;