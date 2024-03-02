var ctx = document.getElementById("mainCanvas")

var data = {
    datasets: [{
        label: 'Done',
        data: data[3].done,
        backgroundColor: '#71bd42'
    }, {
        label: 'In Progress',
        data: data[3].inProgress,
        backgroundColor: '#e05f19'
    }, {
        label: 'Unlabeled',
        data: data[3].unlabeled,
        backgroundColor: "#382911"
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