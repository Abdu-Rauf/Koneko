export function attachInputListeners(dataChannel) {
    console.log('Attaching input listeners')
    
    const video = document.getElementById('display-box')
    // Container's actual resolution
    const CONTAINER_WIDTH = 1920
    const CONTAINER_HEIGHT = 1080
    
    let lastMouseMove = 0
    const MOUSE_THROTTLE = 16
    
    video.addEventListener('mousemove', (e) => {
        const now = Date.now()
        if (now - lastMouseMove < MOUSE_THROTTLE) return
        lastMouseMove = now
        
        const rect = video.getBoundingClientRect()
        // Map displayed coordinates to container resolution
        const x = Math.round((e.clientX - rect.left) / rect.width * CONTAINER_WIDTH)
        const y = Math.round((e.clientY - rect.top) / rect.height * CONTAINER_HEIGHT)
        
        dataChannel.send(JSON.stringify({
            type: 'mouse_move',
            x: x,
            y: y
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
