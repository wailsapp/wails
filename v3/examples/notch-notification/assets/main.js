if (window.location.protocol === "wails:" || window.location.hostname === "wails.localhost") {
  await import("/wails/runtime.js");
}

const cpuValue = document.getElementById("cpu-value");
const memoryValue = document.getElementById("memory-value");
const diskValue = document.getElementById("disk-value");
const memoryPercentLabel = document.getElementById("memory-percent");
const memoryFreeLabel = document.getElementById("memory-free");
const diskPercentLabel = document.getElementById("disk-percent");
const diskFreeLabel = document.getElementById("disk-free");
const updatedAt = document.getElementById("updated-at");
const cpuCanvas = document.getElementById("cpu-chart");
const memoryCanvas = document.getElementById("memory-chart");
const diskCanvas = document.getElementById("disk-chart");

let cpu = 17;
let memory = 10.7;
let disk = 312;

const histories = [
  { values: Array.from({ length: 42 }, () => 18 + Math.random() * 25), canvas: cpuCanvas, color: "#bca8ff", fill: ["rgba(151,116,255,.24)", "rgba(151,116,255,0)"] },
  { values: Array.from({ length: 32 }, () => 58 + Math.random() * 10), canvas: memoryCanvas, color: "#b98dff", fill: ["rgba(185,141,255,.20)", "rgba(185,141,255,0)"] },
  { values: Array.from({ length: 32 }, () => 55 + Math.random() * 8), canvas: diskCanvas, color: "#8ba1ff", fill: ["rgba(109,140,255,.20)", "rgba(109,140,255,0)"] },
];

function nudge(value, min, max, amount) {
  return Math.max(min, Math.min(max, value + (Math.random() - .47) * amount));
}

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

function render() {
  const memoryPercent = Math.round((memory / 16) * 100);
  const diskPercent = Math.round((disk / 512) * 100);
  cpuValue.textContent = String(Math.round(cpu));
  memoryValue.textContent = memory.toFixed(1);
  diskValue.textContent = String(Math.round(disk));
  memoryPercentLabel.textContent = `${memoryPercent}% used`;
  memoryFreeLabel.textContent = `${(16 - memory).toFixed(1)} GB free`;
  diskPercentLabel.textContent = `${diskPercent}% used`;
  diskFreeLabel.textContent = `${Math.round(512 - disk)} GB free`;
  histories.forEach(drawGraph);
}

function sample() {
  cpu = nudge(cpu, 8, 68, 21);
  memory = nudge(memory, 9.1, 12.3, .55);
  disk = nudge(disk, 308, 318, 1.2);
  const nextValues = [cpu, (memory / 16) * 100, (disk / 512) * 100];
  histories.forEach((history, index) => {
    history.values.shift();
    history.values.push(nextValues[index]);
  });
  updatedAt.textContent = "Updated just now";
  render();
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

window.addEventListener("resize", render);
window.setInterval(sample, 1100);
render();
