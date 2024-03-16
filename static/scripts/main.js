let pxHeight = 0;
var data = {}

window.onload = (e) => {
    document.addEventListener("mousemove", (e) => {
        let rect = e.target.getBoundingClientRect();
        pxHeight = e.clientY - rect.top;
    });
}
