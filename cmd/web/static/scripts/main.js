let pxHeight = 0;
var trial3 = []

function setData (data){
    trial3 = data
}

window.onload = (e) => {
    document.addEventListener("mousemove", (e) => {
        let rect = e.target.getBoundingClientRect();
        pxHeight = e.clientY - rect.top;
    });
}


const updateData = (setData)=>{
    trial3 = setData
}