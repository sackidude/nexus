let pxHeight = 0;
document.addEventListener("mousemove", (e) => {
    let rect = e.target.getBoundingClientRect();
    pxHeight = e.clientY - rect.top;
});