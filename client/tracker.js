import InputLogger from './utils/InputLogger.js'

const video = document.getElementById('display-box')
const logger = new InputLogger(video);

const peerConnection = new RTCPeerConnection()
const dataChannel = peerConnection.createDataChannel('inputs')

const offer = await peerConnection.createOffer()

await peerConnection.setLocalDescription(offer)

peerConnection.onicecandidate = (e) => {
    if (e.candidate){
        console.log('New ICE candidate:', e.candidate.candidate)
    }
    // send ice to go server
}

// Log client inputs
video.addEventListener('click',logger.mouseListener.bind(logger))
video.addEventListener('mousemove',logger.mouseListener.bind(logger))
document.addEventListener('keydown', logger.keyListener.bind(logger))


// window.exportInputLogs = () => console.log(logger.exportLogs());
// window.showLogs = () => console.log(logger.getLogs()); 