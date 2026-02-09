import InputLogger from './utils/InputLogger.js'
import { setupPeerConnection } from './utils/peerSetup.js'
import {signalling } from './utils/signalling.js'
import {attachInputListeners} from './utils/UserInputs.js'


const video = document.getElementById('display-box')
const logger = new InputLogger(video);

// Websocket To Handle Signalling

const socket = new WebSocket('ws://localhost:8080/ws')
let peerConnection;
let dataChannel

socket.onopen = async () =>{
    console.log('WebSocket connected!')

    const selectBrowser = (Browser) => {
        console.log(`${Browser} selected`)

        socket.send(JSON.stringify({
            type: 'browser_select',
            browser: Browser
        }))
        peerConnection = setupPeerConnection(socket ,(channel)=>{
            dataChannel = channel
            attachInputListeners(dataChannel)
        })
        // Change the Video elements Visibility
        document.getElementById('selection-screen').classList.add('hidden')
        document.getElementById('streaming-screen').classList.remove('hidden')
    }
    document.getElementById('chrome-btn').addEventListener('click',()=>selectBrowser('chrome'))
    document.getElementById('firefox-btn').addEventListener('click' , ()=>selectBrowser('firefox'))

}
socket.onmessage = async (e) =>{
    // parse the message
    const msg = JSON.parse(e.data)
    signalling(msg,peerConnection,socket)

}

socket.onclose = (event) => {
    console.log('WebSocket closed:', event.code, event.reason)
    // WebSocket closed, but peer connection will work
}

socket.onerror = (error) => {
    console.error('WebSocket error:', error)
}

function InputStream(data){
    if (dataChannel && dataChannel.readyState==='open'){
        dataChannel.send(JSON.stringify(data))
    }
    else{
        console.warn('Data channel not read, state:', dataChannel?.readyState)
    }
}
function sendTestMessage() {
    if (dataChannel && dataChannel.readyState === 'open') {
        dataChannel.send('Hello from test!')
        console.log('Test message sent!')
    } else {
        console.warn('Data channel not ready')
    }
}


window.sendTest = sendTestMessage

// Log client inputs
// video.addEventListener('click',logger.mouseListener.bind(logger))
// video.addEventListener('mousemove',logger.mouseListener.bind(logger))
// document.addEventListener('keydown', logger.keyListener.bind(logger))














































// // window.exportInputLogs = () => console.log(logger.exportLogs());
// // window.showLogs = () => console.log(logger.getLogs());