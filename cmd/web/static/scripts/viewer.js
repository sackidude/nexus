const ctx = document.getElementById("mainCanvas")

const updateData = (setData)=>{
    trial3 = setData
}

const data = {
    datasets: [{
        label: 'trial-3',
        data: trial3,
        backgroundColor: 'rgb(255, 99, 132)'
    }]
};
new Chart(ctx, {
    type: 'scatter',
    data: data,
    options: {
        scales: {
            x: {
                type: "linear",
                position: "bottom"
            }
        }
    }
});