var trial_data = data[trial_num];
var ctx = document.getElementById("main-canvas-" + trial_num);

var chartData = {
  datasets: [
    {
      label: "Done",
      data: trial_data.done,
      backgroundColor: "#71bd42",
    },
    {
      label: "In Progress",
      data: trial_data.inProgress,
      backgroundColor: "#e05f19",
    },
    {
      label: "Unlabeled",
      data: trial_data.unlabeled,
      backgroundColor: "#382911",
    },
  ],
};
new Chart(ctx, {
  type: "scatter",
  data: chartData,
  options: {
    scales: {
      x: {
        type: "linear",
        position: "bottom",
      },
    },
  },
});
