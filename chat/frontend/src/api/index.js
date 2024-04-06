var socket = new WebSocket("ws://localhost:8080/ws")


let connect = (cb) => { //不是很懂cb是干嘛的，或者说cb是什么东西
    console.log("Attempting to connect...")

    socket.onopen = () => {
        console.log("Connected to server successfully");
    };

    socket.onmessage = (message) => {
        console.log("Received message from server: " + message.data);
        cb(message)
    };
    
    socket.onclose = (event) => {
        console.log("Socket closed connection: ", event);
    };

    socket.onerror = (error) => {
        console.log("Socket Error: ", error);
    };
};

let sendMsg = (msg) => {
    console.log("sending msg: ", msg);
    socket.send(msg);
};

export{connect,sendMsg}



