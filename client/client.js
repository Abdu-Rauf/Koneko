import InputLogger from './utils/InputLogger.js'
import { setupPeerConnection } from './utils/signalling.js'

const video = document.getElementById('display-box')
const logger = new InputLogger(video);

// SIGNALING WITH WEBSOCKET

const socket = new WebSocket('ws://localhost:8080/ws')
let peerConnection;
let dataChannel

function attachInputListeners() {
    console.log('Attaching input listeners')
    
    const video = document.getElementById('display-box')
    let lastMouseMove = 0
    const MOUSE_THROTTLE = 16
    
    // Mouse movement on video element
    video.addEventListener('mousemove', (e) => {
        const now = Date.now()
        if (now - lastMouseMove < MOUSE_THROTTLE) return
        const rect = video.getBoundingClientRect()
        const x = e.clientX - rect.left
        const y = e.clientY - rect.top
        
        // Scale to container resolution 
        const scaleX = 1920 / rect.width
        const scaleY = 1080 / rect.height
        
        const containerX = Math.floor(x * scaleX)
        const containerY = Math.floor(y * scaleY)
        
        console.log('Mouse moved:', containerX, containerY)
        dataChannel.send(JSON.stringify({
            type: 'mouse_move',
            x: containerX,
            y: containerY
        }))
    })
    
    // Mouse click on video element
    video.addEventListener('mousedown', (e) => {
        e.preventDefault()
        
        const rect = video.getBoundingClientRect()
        const x = e.clientX - rect.left
        const y = e.clientY - rect.top
        
        const scaleX = 1920 / rect.width
        const scaleY = 1080 / rect.height
        
        const containerX = Math.floor(x * scaleX)
        const containerY = Math.floor(y * scaleY)
        
        console.log('Mouse clicked:', containerX, containerY, 'button:', e.button)
        dataChannel.send(JSON.stringify({
            type: 'mouse_click',
            click: e.button,
            x: containerX,
            y: containerY
        }))
    })
    
    document.addEventListener('keydown', (e) => {
        e.preventDefault()
        console.log('Key pressed:', e.key)
        dataChannel.send(JSON.stringify({
            type: 'key_press',
            key: e.key
        }))
    })
    
    // Prevent right-click menu on video
    video.addEventListener('contextmenu', (e) => {
        e.preventDefault()
    })
    
    console.log('Input listeners attached to video element')
}

socket.onopen = async () =>{
    console.log('WebSocket connected!')

    const selectBrowser = (Browser) => {
        console.log(`${Browser} selected`)
        peerConnection = setupPeerConnection(socket ,(channel)=>{
            dataChannel = channel
        })
        socket.send(JSON.stringify({
            type: 'browser_select',
            browser: Browser
        }))
        // Change the Video elements Visibility
        document.getElementById('selection-screen').classList.add('hidden')
        document.getElementById('streaming-screen').classList.remove('hidden')
        attachInputListeners()
    }
    document.getElementById('chrome-btn').addEventListener('click',()=>selectBrowser('chrome'))
    document.getElementById('firefox-btn').addEventListener('click' , ()=>selectBrowser('firefox'))

}
socket.onmessage = async (e) =>{
    // parse the message
    const msg = JSON.parse(e.data)

    if (msg.type === "offer") {
        console.log("Received offer from server")
        
        await peerConnection.setRemoteDescription(new RTCSessionDescription({
            type: "offer",
            sdp: msg.sdp
        }))
        
        // Create and send answer
        const answer = await peerConnection.createAnswer()
        await peerConnection.setLocalDescription(answer)
        
        console.log("Sending answer to server")
        socket.send(JSON.stringify({
            type: 'answer',
            sdp: answer.sdp
        }))
    }

    // add the recieved ice candidates
    if (msg.type == "ice-candidate"){
        console.log("recieved ice candidate",msg)
        var ice = new RTCIceCandidate(msg.candidate)
        peerConnection.addIceCandidate(ice)
    }
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