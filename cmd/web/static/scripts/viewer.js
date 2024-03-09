for (const trial_num in data) {
  const trial_data = data[trial_num];
  const ctx = document.getElementById("canvas-" + trial_num);

  const chartData = {
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
}
