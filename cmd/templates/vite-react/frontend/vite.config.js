import path from "path";
import reactRefresh from "@vitejs/plugin-react-refresh";

export default {
  plugins: [reactRefresh()],
  build: {
    lib: {
      entry: path.resolve(__dirname, "src/main.tsx"),
      formats: ["es"],
      name: "main",
    },
  },
};
