export function setupPeerConnection(socket , onDataChannelCreated){

    const peerConnection = new RTCPeerConnection()
    // SETTING UP DATA CHANNEL
    peerConnection.ondatachannel = (event) => {
        console.log('Data channel recieved')
        const dataChannel = event.channel
        onDataChannelCreated(dataChannel)
        
        dataChannel.onopen = ()=>{
            console.log("data channel opened inform server")

            dataChannel.send(JSON.stringify({
                type: "dc_ready"
            }))
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
    }

    // Configure recieving video bytes
    peerConnection.ontrack = (event) => {
        console.log("Video Track Received", event) 
        console.log("Streams:", event.streams)       
        const video = document.getElementById('display-box')
        video.srcObject = event.streams[0]
        video.play()
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


    peerConnection.onconnectionstatechange = () => {
        console.log('Peer connection state:', peerConnection.connectionState)

        if (peerConnection.connectionState === 'connected') {
            console.log('Peer connection established, closing WebSocket')
            socket.close()  
        }
        if (peerConnection.connectionState === 'failed' ||
            peerConnection.connectionState === 'closed' ||
            peerConnection.connectionState === 'disconnected') {
            console.log('Peer connection died')
            // Cleanup or reconnect logic here
        }
    }
    return peerConnection

}