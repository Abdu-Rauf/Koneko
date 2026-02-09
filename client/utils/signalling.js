export async function signalling(msg,peerConnection,socket){

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