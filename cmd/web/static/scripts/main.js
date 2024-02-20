let pxHeight = 0;
var trial3 = []

function setData (data){
    trial3 = data
}

document.onload = (e) => {
    document.addEventListener("mousemove", (e) => {
        let rect = e.target.getBoundingClientRect();
        pxHeight = e.clientY - rect.top;
    });
}

