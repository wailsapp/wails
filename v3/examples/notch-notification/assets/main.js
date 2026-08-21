if (window.location.protocol === "wails:" || window.location.hostname === "wails.localhost") {
  await import("/wails/runtime.js");
}

const cpuValue = document.getElementById("cpu-value");
const coreCount = document.getElementById("core-count");
const machineName = document.getElementById("machine-name");
const memoryTotal = document.getElementById("memory-total");
const memoryValue = document.getElementById("memory-value");
const diskValue = document.getElementById("disk-value");
const memoryPercentLabel = document.getElementById("memory-percent");
const memoryFreeLabel = document.getElementById("memory-free");
const diskPercentLabel = document.getElementById("disk-percent");
const diskFreeLabel = document.getElementById("disk-free");
const peakValue = document.getElementById("peak-value");
const cpuCanvas = document.getElementById("cpu-chart");
const memoryProgress = document.getElementById("memory-progress");
const diskProgress = document.getElementById("disk-progress");

const history = {
  values: Array(52).fill(0),
  canvas: cpuCanvas,
  color: "#38a8ff",
  fill: ["rgba(56,168,255,.24)", "rgba(56,168,255,0)"],
};
let hasSample = false;

function drawGraph({ values, canvas, color, fill }) {
  const bounds = canvas.getBoundingClientRect();
  const scale = window.devicePixelRatio || 1;
  const width = Math.max(1, Math.round(bounds.width * scale));
  const height = Math.max(1, Math.round(bounds.height * scale));
  if (canvas.width !== width || canvas.height !== height) {
    canvas.width = width;
    canvas.height = height;
  }

  const context = canvas.getContext("2d");
  if (!context) return;
  const pointX = index => (index / (values.length - 1)) * width;
  const pointY = value => height - 3 * scale - (value / 100) * (height - 6 * scale);

  context.clearRect(0, 0, width, height);
  const area = context.createLinearGradient(0, 0, 0, height);
  area.addColorStop(0, fill[0]);
  area.addColorStop(1, fill[1]);

  context.beginPath();
  values.forEach((value, index) => {
    const x = pointX(index);
    const y = pointY(value);
    if (index === 0) context.moveTo(x, y);
    else context.lineTo(x, y);
  });
  context.lineTo(width, height);
  context.lineTo(0, height);
  context.closePath();
  context.fillStyle = area;
  context.fill();

  context.beginPath();
  values.forEach((value, index) => {
    const x = pointX(index);
    const y = pointY(value);
    if (index === 0) context.moveTo(x, y);
    else context.lineTo(x, y);
  });
  context.strokeStyle = color;
  context.lineWidth = 1.5 * scale;
  context.lineJoin = "round";
  context.lineCap = "round";
  context.stroke();
}

function render(stats) {
  const memoryPercent = Math.round(stats.memoryPercent);
  const diskPercent = Math.round(stats.diskPercent);
  cpuValue.textContent = String(Math.round(stats.cpuPercent));
  coreCount.textContent = String(stats.coreCount);
  machineName.textContent = stats.machineName || "Apple silicon";
  memoryTotal.textContent = String(Math.round(stats.memoryUsedGB + stats.memoryFreeGB));
  memoryValue.textContent = stats.memoryUsedGB.toFixed(1);
  diskValue.textContent = String(Math.round(stats.diskUsedGB));
  memoryPercentLabel.textContent = `${memoryPercent}% used`;
  memoryFreeLabel.textContent = `${stats.memoryFreeGB.toFixed(1)} GB free`;
  diskPercentLabel.textContent = `${diskPercent}% used`;
  diskFreeLabel.textContent = `${Math.round(stats.diskFreeGB)} GB free`;
  memoryProgress.value = memoryPercent;
  diskProgress.value = diskPercent;
  peakValue.textContent = `${Math.round(Math.max(...history.values))}% peak`;
  drawGraph(history);
}

function renderUnavailable() {
  history.values.fill(0);
  hasSample = false;
  cpuValue.textContent = "—";
  coreCount.textContent = "—";
  machineName.textContent = "Telemetry unavailable";
  memoryTotal.textContent = "—";
  memoryValue.textContent = "—";
  diskValue.textContent = "—";
  memoryPercentLabel.textContent = "Unavailable";
  memoryFreeLabel.textContent = "—";
  diskPercentLabel.textContent = "Unavailable";
  diskFreeLabel.textContent = "—";
  peakValue.textContent = "Unavailable";
  memoryProgress.value = 0;
  diskProgress.value = 0;
  drawGraph(history);
}

async function sample() {
  if (!globalThis.wails?.Call?.ByName) return;
  try {
    const stats = await globalThis.wails.Call.ByName("main.NotificationController.Stats");
    if (!stats?.available) {
      renderUnavailable();
      return;
    }
    if (!hasSample) {
      history.values.fill(stats.cpuPercent);
      hasSample = true;
    } else {
      history.values.shift();
      history.values.push(stats.cpuPercent);
    }
    render(stats);
  } catch (error) {
    renderUnavailable();
    console.error("reading system telemetry", error);
  }
}

async function hideWindow() {
  if (globalThis.wails?.Call?.ByName) {
    await globalThis.wails.Call.ByName("main.NotificationController.Hide");
  }
}

async function quitApplication() {
  if (globalThis.wails?.Call?.ByName) {
    await globalThis.wails.Call.ByName("main.NotificationController.Quit");
  }
}

document.getElementById("hide").addEventListener("click", hideWindow);
document.getElementById("quit").addEventListener("click", quitApplication);

window.addEventListener("resize", () => drawGraph(history));
window.setInterval(sample, 1100);
sample();
