import InputLogger from './utils/InputLogger.js'

const video = document.getElementById('display-box')
const logger = new InputLogger(video);

// SIGNALING WITH WEBSOCKET

const socket = new WebSocket('ws://localhost:8080/ws')
let peerConnection;
let dataChannel

socket.onopen = async () =>{
    console.log('WebSocket connected!')

    peerConnection = new RTCPeerConnection()
    dataChannel = peerConnection.createDataChannel('inputs')

    // SETTING UP DATA CHANNEL
    dataChannel.onopen = () => {
        console.log('Data channel opened')
        // socket.close()
    }

    dataChannel.onmessage = (event) => {
        console.log('Received via data channel:', event.data)
    }

    dataChannel.onerror = (error) => {
        console.error('Data channel error:', error)
    }

    dataChannel.onclose = () => {
        console.log('Data channel closed')
    }

    // SEND ICE CANDIDATES TO GO SERVER
    peerConnection.onicecandidate = (e) => {
        console.log('sending ice candidate', e.candidate)
        if (e.candidate){
            socket.send(JSON.stringify({
                type:'ice-candidate',
                candidate: e.candidate
            }));
        }
    }
    // CREATE OFFER AND SEND SDP TO GO SERVER
    const offer = await peerConnection.createOffer()
    await peerConnection.setLocalDescription(offer)

    console.log('sending offer', offer.sdp) 
    socket.send(JSON.stringify({
        type: 'offer',
        sdp: offer.sdp
    }))
}
socket.onmessage = (e) =>{
    // parse the message
    const msg = JSON.parse(e.data)

    // set remote desc to answer
    if (msg.type == "answer"){
        console.log("recieved answer candidate",msg)
        var rd = new RTCSessionDescription({type:"answer",sdp:msg.sdp})
        peerConnection.setRemoteDescription(rd)
    }
    // add the recieved ice candidates
    if (msg.type == "ice-candidate"){
        console.log("recieved ice candidate",msg)
        var ice = new RTCIceCandidate(msg.candidate)
        peerConnection.addIceCandidate(ice)
    }
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
video.addEventListener('click',logger.mouseListener.bind(logger))
video.addEventListener('mousemove',logger.mouseListener.bind(logger))
document.addEventListener('keydown', logger.keyListener.bind(logger))














































// // window.exportInputLogs = () => console.log(logger.exportLogs());
// // window.showLogs = () => console.log(logger.getLogs()); 